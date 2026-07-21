package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var granularites = map[string]string{
	"jour":    "day",
	"semaine": "week",
	"mois":    "month",
}

func lireDate(valeur string, defaut time.Time) time.Time {
	if valeur == "" {
		return defaut
	}
	moment, erreur := time.Parse("2006-01-02", valeur)
	if erreur != nil {
		return defaut
	}
	return moment
}

func (s *Serveur) statistiquesGroupe(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "lecteur") {
		return
	}
	contexte := c.Request.Context()
	maintenant := time.Now()
	fin := lireDate(c.Query("fin"), maintenant).Add(24 * time.Hour)
	debut := lireDate(c.Query("debut"), maintenant.AddDate(0, 0, -30))
	granularite := granularites[c.DefaultQuery("granularite", "jour")]
	if granularite == "" {
		granularite = "day"
	}
	var projets []string
	if cible := c.Query("projet"); cible != "" {
		projet, erreur := s.Depot.Projet(contexte, cible)
		if erreur != nil || projet.Groupe != identifiant {
			c.JSON(http.StatusNotFound, gin.H{"erreur": "projet introuvable dans ce groupe"})
			return
		}
		projets = []string{projet.ID}
	} else {
		var erreur error
		projets, erreur = s.Depot.IdentifiantsProjets(contexte, identifiant)
		if erreur != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des projets impossible"})
			return
		}
	}
	statistiques, erreur := s.Depot.Statistiques(contexte, projets, debut, fin, granularite)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "calcul des statistiques impossible"})
		return
	}
	c.JSON(http.StatusOK, statistiques)
}
