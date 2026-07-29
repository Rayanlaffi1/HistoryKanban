package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var clientMiseAJour = &http.Client{Timeout: 10 * time.Second}

func (s *Serveur) exigerAdministrateurPlateforme() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, role := range s.revendications(c).Roles {
			if role == "administrateur" {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erreur": "reserve aux administrateurs de la plateforme"})
	}
}

func (s *Serveur) jenkinsConfigure() bool {
	return s.Config.JenkinsURL != "" && s.Config.JenkinsJob != "" &&
		s.Config.JenkinsUtilisateur != "" && s.Config.JenkinsJeton != ""
}

func normaliserVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
}

func (s *Serveur) derniereVersionGitHub(contexte context.Context) (string, error) {
	adresse := "https://api.github.com/repos/" + s.Config.GithubDepot + "/releases/latest"
	requete, erreur := http.NewRequestWithContext(contexte, http.MethodGet, adresse, nil)
	if erreur != nil {
		return "", erreur
	}
	requete.Header.Set("Accept", "application/vnd.github+json")
	reponse, erreur := clientMiseAJour.Do(requete)
	if erreur != nil {
		return "", erreur
	}
	defer reponse.Body.Close()
	if reponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("reponse GitHub %d", reponse.StatusCode)
	}
	var corps struct {
		Tag string `json:"tag_name"`
	}
	if erreur := json.NewDecoder(io.LimitReader(reponse.Body, 1<<20)).Decode(&corps); erreur != nil {
		return "", erreur
	}
	return corps.Tag, nil
}

func (s *Serveur) etatMiseAJour(c *gin.Context) {
	etat := gin.H{
		"versionActuelle":  s.Config.VersionApplication,
		"derniereVersion":  "",
		"majDisponible":    false,
		"jenkinsConfigure": s.jenkinsConfigure(),
	}
	if s.Config.GithubDepot != "" {
		derniere, erreur := s.derniereVersionGitHub(c.Request.Context())
		if erreur != nil {
			log.Printf("lecture de la derniere release GitHub impossible : %v", erreur)
		} else {
			etat["derniereVersion"] = derniere
			etat["majDisponible"] = normaliserVersion(derniere) != "" &&
				normaliserVersion(derniere) != normaliserVersion(s.Config.VersionApplication)
		}
	}
	c.JSON(http.StatusOK, etat)
}

func cheminJob(job string) string {
	segments := strings.Split(job, "/")
	for indice := range segments {
		segments[indice] = url.PathEscape(segments[indice])
	}
	return "/job/" + strings.Join(segments, "/job/")
}

func (s *Serveur) declencherMiseAJour(c *gin.Context) {
	if !s.jenkinsConfigure() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"erreur": "declenchement Jenkins non configure sur le serveur"})
		return
	}
	adresse := strings.TrimSuffix(s.Config.JenkinsURL, "/") + cheminJob(s.Config.JenkinsJob) + "/build?delay=0sec"
	requete, erreur := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, adresse, nil)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "preparation de l'appel Jenkins impossible"})
		return
	}
	requete.SetBasicAuth(s.Config.JenkinsUtilisateur, s.Config.JenkinsJeton)
	reponse, erreur := clientMiseAJour.Do(requete)
	if erreur != nil {
		log.Printf("appel Jenkins impossible : %v", erreur)
		c.JSON(http.StatusBadGateway, gin.H{"erreur": "Jenkins injoignable"})
		return
	}
	defer reponse.Body.Close()
	if reponse.StatusCode != http.StatusCreated && reponse.StatusCode != http.StatusOK {
		log.Printf("declenchement Jenkins refuse (code %d)", reponse.StatusCode)
		c.JSON(http.StatusBadGateway, gin.H{"erreur": "declenchement Jenkins refuse"})
		return
	}
	log.Printf("mise a jour declenchee via Jenkins par %s", s.revendications(c).Utilisateur)
	c.JSON(http.StatusAccepted, gin.H{"etat": "declenchee"})
}
