package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"historykanban/serveur/internal/modeles"
)

const longueurMaximaleSousTache = 200

type corpsSousTache struct {
	Libelle string `json:"libelle"`
	Faite   *bool  `json:"faite"`
}

func (s *Serveur) listerSousTaches(c *gin.Context) {
	tache, _, autorise := s.tacheAutorisee(c, "lecteur")
	if !autorise {
		return
	}
	sousTaches, erreur := s.Depot.SousTaches(c.Request.Context(), tache.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des sous-taches impossible"})
		return
	}
	c.JSON(http.StatusOK, sousTaches)
}

func (s *Serveur) creerSousTache(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	var corps corpsSousTache
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "libelle requis"})
		return
	}
	libelle := strings.TrimSpace(corps.Libelle)
	if libelle == "" || len([]rune(libelle)) > longueurMaximaleSousTache {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "le libelle doit faire entre 1 et 200 caracteres"})
		return
	}
	sousTache, erreur := s.Depot.CreerSousTache(c.Request.Context(), tache.ID, libelle)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "creation de la sous-tache impossible"})
		return
	}
	s.rafraichirTache(c, projet, tache.ID)
	c.JSON(http.StatusCreated, sousTache)
}

func (s *Serveur) sousTacheAutorisee(c *gin.Context) (*modeles.SousTache, *modeles.Projet, bool) {
	sousTache, erreur := s.Depot.SousTache(c.Request.Context(), c.Param("id"))
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "sous-tache introuvable"})
		return nil, nil, false
	}
	tache, erreur := s.Depot.Tache(c.Request.Context(), sousTache.Tache)
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "tache introuvable"})
		return nil, nil, false
	}
	projet, autorise := s.exigerRoleProjet(c, tache.Projet, "membre")
	if !autorise {
		return nil, nil, false
	}
	return sousTache, projet, true
}

func (s *Serveur) modifierSousTache(c *gin.Context) {
	sousTache, projet, autorise := s.sousTacheAutorisee(c)
	if !autorise {
		return
	}
	var corps corpsSousTache
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "corps invalide"})
		return
	}
	libelle := sousTache.Libelle
	if strings.TrimSpace(corps.Libelle) != "" {
		libelle = strings.TrimSpace(corps.Libelle)
	}
	if len([]rune(libelle)) > longueurMaximaleSousTache {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "le libelle ne doit pas depasser 200 caracteres"})
		return
	}
	faite := sousTache.Faite
	if corps.Faite != nil {
		faite = *corps.Faite
	}
	if erreur := s.Depot.ModifierSousTache(c.Request.Context(), sousTache.ID, libelle, faite); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "modification de la sous-tache impossible"})
		return
	}
	s.rafraichirTache(c, projet, sousTache.Tache)
	sousTache.Libelle = libelle
	sousTache.Faite = faite
	c.JSON(http.StatusOK, sousTache)
}

func (s *Serveur) supprimerSousTache(c *gin.Context) {
	sousTache, projet, autorise := s.sousTacheAutorisee(c)
	if !autorise {
		return
	}
	if erreur := s.Depot.SupprimerSousTache(c.Request.Context(), sousTache.ID); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression de la sous-tache impossible"})
		return
	}
	s.rafraichirTache(c, projet, sousTache.Tache)
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}

func (s *Serveur) rafraichirTache(c *gin.Context, projet *modeles.Projet, tache string) {
	resultat, erreur := s.Depot.Tache(c.Request.Context(), tache)
	if erreur != nil {
		return
	}
	s.remplirURLsTache(resultat)
	s.publier("tache.modifiee", projet, resultat.Titre, resultat, nil, c)
}
