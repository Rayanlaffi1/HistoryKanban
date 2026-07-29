package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var granularites = map[string]string{
	"jour":    "day",
	"semaine": "week",
	"mois":    "month",
}

var fuseauStatistiques = chargerFuseauStatistiques()

func chargerFuseauStatistiques() *time.Location {
	lieu, erreur := time.LoadLocation("Europe/Paris")
	if erreur != nil {
		return time.Local
	}
	return lieu
}

func debutJourLocal(moment time.Time) time.Time {
	local := moment.In(fuseauStatistiques)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, fuseauStatistiques)
}

func lireDate(valeur string, defaut time.Time) time.Time {
	if valeur == "" {
		return debutJourLocal(defaut)
	}
	moment, erreur := time.ParseInLocation("2006-01-02", valeur, fuseauStatistiques)
	if erreur != nil {
		return debutJourLocal(defaut)
	}
	return moment
}

func lireEntier(valeur string, defaut int) int {
	if valeur == "" {
		return defaut
	}
	entier, erreur := strconv.Atoi(valeur)
	if erreur != nil {
		return defaut
	}
	return entier
}

func (s *Serveur) statistiquesGroupe(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "lecteur") {
		return
	}
	contexte := c.Request.Context()
	maintenant := time.Now()
	debut := lireDate(c.Query("debut"), maintenant.AddDate(0, 0, -30))
	fin := lireDate(c.Query("fin"), maintenant).AddDate(0, 0, 1)
	limite := lireEntier(c.Query("limite"), 100)
	offset := lireEntier(c.Query("offset"), 0)
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
	statistiques, erreur := s.Depot.Statistiques(contexte, projets, debut, fin, granularite, limite, offset)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "calcul des statistiques impossible"})
		return
	}
	c.JSON(http.StatusOK, statistiques)
}
