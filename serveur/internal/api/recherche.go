package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const limiteRecherche = 40

func (s *Serveur) rechercher(c *gin.Context) {
	texte := strings.TrimSpace(c.Query("texte"))
	if len([]rune(texte)) < 2 {
		c.JSON(http.StatusOK, []any{})
		return
	}
	resultats, erreur := s.Depot.Rechercher(c.Request.Context(),
		s.revendications(c).Utilisateur, texte, limiteRecherche)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "recherche impossible"})
		return
	}
	c.JSON(http.StatusOK, resultats)
}
