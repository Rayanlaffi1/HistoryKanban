package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Serveur) obtenirCle(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "membre")
	if !autorise {
		return
	}
	cle, erreur := s.Depot.Cle(c.Request.Context(), projet.ID, s.revendications(c).Utilisateur)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture de la cle impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cle": cle})
}

func (s *Serveur) genererCle(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "membre")
	if !autorise {
		return
	}
	octets := make([]byte, 32)
	if _, erreur := rand.Read(octets); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "generation de la cle impossible"})
		return
	}
	valeur := "hk." + hex.EncodeToString(octets)
	cle, erreur := s.Depot.EnregistrerCle(c.Request.Context(), projet.ID, s.revendications(c).Utilisateur, valeur)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "enregistrement de la cle impossible"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"cle": cle})
}

func (s *Serveur) supprimerCle(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "membre")
	if !autorise {
		return
	}
	if erreur := s.Depot.SupprimerCle(c.Request.Context(), projet.ID, s.revendications(c).Utilisateur); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression de la cle impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}
