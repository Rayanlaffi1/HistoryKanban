package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"historykanban/serveur/internal/evenements"
	"historykanban/serveur/internal/modeles"
	"historykanban/serveur/internal/securite"
)

var niveaux = map[string]int{
	"lecteur":        1,
	"membre":         2,
	"administrateur": 3,
	"proprietaire":   4,
}

func (s *Serveur) authentifier() gin.HandlerFunc {
	return func(c *gin.Context) {
		entete := c.GetHeader("Authorization")
		if !strings.HasPrefix(entete, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erreur": "jeton absent"})
			return
		}
		revendications, erreur := s.Verificateur.Verifier(strings.TrimPrefix(entete, "Bearer "))
		if erreur != nil {
			log.Printf("verification du jeton echouee : %v", erreur)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erreur": "jeton invalide"})
			return
		}
		if erreur := s.Depot.AssurerUtilisateur(c.Request.Context(), modeles.Utilisateur{
			ID:       revendications.Utilisateur,
			Courriel: revendications.Courriel,
			Nom:      revendications.Nom,
			Prenom:   revendications.Prenom,
		}); erreur != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"erreur": "enregistrement du profil impossible"})
			return
		}
		if revendications.Session != "" {
			valide, remplacee, erreur := s.Depot.ReclamerSession(
				c.Request.Context(), revendications.Utilisateur, revendications.Session, revendications.Connexion)
			if erreur != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"erreur": "verification de session impossible"})
				return
			}
			if !valide {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erreur": "session remplacee", "code": "session.remplacee"})
				return
			}
			if remplacee {
				s.Concentrateur.RemplacerSession(revendications.Utilisateur, revendications.Session)
			}
		}
		c.Set("revendications", revendications)
		c.Next()
	}
}

func (s *Serveur) revendications(c *gin.Context) *securite.Revendications {
	valeur, _ := c.Get("revendications")
	return valeur.(*securite.Revendications)
}

func (s *Serveur) exigerRoleGroupe(c *gin.Context, groupe, minimum string) bool {
	role, erreur := s.Depot.RoleGroupe(c.Request.Context(), groupe, s.revendications(c).Utilisateur)
	if erreur != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"erreur": "verification des droits impossible"})
		return false
	}
	if role == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"erreur": "groupe introuvable"})
		return false
	}
	if niveaux[role] < niveaux[minimum] {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erreur": "droits insuffisants"})
		return false
	}
	return true
}

func (s *Serveur) exigerRoleProjet(c *gin.Context, identifiant, minimum string) (*modeles.Projet, bool) {
	projet, erreur := s.Depot.Projet(c.Request.Context(), identifiant)
	if erreur != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"erreur": "projet introuvable"})
		return nil, false
	}
	if !s.exigerRoleGroupe(c, projet.Groupe, minimum) {
		return nil, false
	}
	return projet, true
}

func (s *Serveur) publier(categorie string, projet *modeles.Projet, titre string, donnees any, destinataires []string, c *gin.Context) {
	revendications := s.revendications(c)
	octets, erreur := json.Marshal(donnees)
	if erreur != nil {
		return
	}
	evenement := evenements.Evenement{
		Type:          categorie,
		Acteur:        revendications.Utilisateur,
		ActeurNom:     strings.TrimSpace(revendications.Prenom + " " + revendications.Nom),
		Titre:         titre,
		Donnees:       octets,
		Destinataires: destinataires,
		Agent:         viaAgent(c),
	}
	if projet != nil {
		evenement.Projet = projet.ID
		evenement.Groupe = projet.Groupe
	}
	s.Bus.Publier(evenement)
}
