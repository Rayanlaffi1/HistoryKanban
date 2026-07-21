package api

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"historykanban/serveur/internal/modeles"
)

func (s *Serveur) moi(c *gin.Context) {
	revendications := s.revendications(c)
	c.JSON(http.StatusOK, gin.H{
		"id":       revendications.Utilisateur,
		"courriel": revendications.Courriel,
		"nom":      revendications.Nom,
		"prenom":   revendications.Prenom,
		"roles":    revendications.Roles,
	})
}

func (s *Serveur) obtenirPreferences(c *gin.Context) {
	preferences, erreur := s.Depot.Preferences(c.Request.Context(), s.revendications(c).Utilisateur)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des preferences impossible"})
		return
	}
	c.JSON(http.StatusOK, preferences)
}

type corpsPreferences struct {
	Courriels bool            `json:"courriels"`
	Types     map[string]bool `json:"types"`
}

func (s *Serveur) enregistrerPreferences(c *gin.Context) {
	var corps corpsPreferences
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "corps invalide"})
		return
	}
	if corps.Types == nil {
		corps.Types = map[string]bool{}
	}
	preferences := modeles.Preferences{
		Utilisateur: s.revendications(c).Utilisateur,
		Courriels:   corps.Courriels,
		Types:       corps.Types,
	}
	if erreur := s.Depot.EnregistrerPreferences(c.Request.Context(), preferences); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "enregistrement des preferences impossible"})
		return
	}
	c.JSON(http.StatusOK, preferences)
}

func (s *Serveur) listerNotifications(c *gin.Context) {
	notifications, erreur := s.Depot.Notifications(c.Request.Context(), s.revendications(c).Utilisateur)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des notifications impossible"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func (s *Serveur) marquerLue(c *gin.Context) {
	if erreur := s.Depot.MarquerLue(c.Request.Context(), c.Param("id"), s.revendications(c).Utilisateur); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "mise a jour impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "lue"})
}

func (s *Serveur) toutMarquerLu(c *gin.Context) {
	if erreur := s.Depot.MarquerToutesLues(c.Request.Context(), s.revendications(c).Utilisateur); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "mise a jour impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "lues"})
}

func (s *Serveur) webhookKeycloak(c *gin.Context) {
	corps, erreur := io.ReadAll(io.LimitReader(c.Request.Body, 65536))
	if erreur == nil && len(corps) > 0 {
		log.Printf("evenement keycloak recu : %s", string(corps))
	}
	c.JSON(http.StatusOK, gin.H{"etat": "recu"})
}
