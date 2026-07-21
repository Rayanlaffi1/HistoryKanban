package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"historykanban/serveur/internal/modeles"
	"historykanban/serveur/internal/securite"
)

func (s *Serveur) authentifierAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		valeur := c.GetHeader("X-Cle-API")
		if valeur == "" {
			entete := c.GetHeader("Authorization")
			if strings.HasPrefix(entete, "Bearer hk.") {
				valeur = strings.TrimPrefix(entete, "Bearer ")
			}
		}
		if valeur == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erreur": "cle d'API absente"})
			return
		}
		cle, utilisateur, erreur := s.Depot.CleParValeur(c.Request.Context(), valeur)
		if erreur != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erreur": "cle d'API invalide"})
			return
		}
		c.Set("revendications", &securite.Revendications{
			Utilisateur: utilisateur.ID,
			Courriel:    utilisateur.Courriel,
			Nom:         utilisateur.Nom,
			Prenom:      utilisateur.Prenom,
		})
		c.Set("projetagent", cle.Projet)
		c.Next()
	}
}

func (s *Serveur) forcerProjetAgent(suivant gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Params = append(c.Params, gin.Param{Key: "id", Value: c.GetString("projetagent")})
		suivant(c)
	}
}

func (s *Serveur) verifierTacheAgent(suivant gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		projet, erreur := s.Depot.ProjetTache(c.Request.Context(), c.Param("id"))
		if erreur != nil || projet != c.GetString("projetagent") {
			c.JSON(http.StatusNotFound, gin.H{"erreur": "tache introuvable dans ce projet"})
			return
		}
		suivant(c)
	}
}

func (s *Serveur) verifierEtiquetteAgent(suivant gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		projet, erreur := s.Depot.ProjetEtiquette(c.Request.Context(), c.Param("id"))
		if erreur != nil || projet != c.GetString("projetagent") {
			c.JSON(http.StatusNotFound, gin.H{"erreur": "etiquette introuvable dans ce projet"})
			return
		}
		suivant(c)
	}
}

func (s *Serveur) verifierLotAgent(suivant gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		projet, erreur := s.Depot.ProjetLot(c.Request.Context(), c.Param("id"))
		if erreur != nil || projet != c.GetString("projetagent") {
			c.JSON(http.StatusNotFound, gin.H{"erreur": "lot introuvable dans ce projet"})
			return
		}
		suivant(c)
	}
}

func (s *Serveur) agentListerMembres(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.GetString("projetagent"), "lecteur")
	if !autorise {
		return
	}
	membres, erreur := s.Depot.Membres(c.Request.Context(), projet.Groupe)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des membres impossible"})
		return
	}
	c.JSON(http.StatusOK, membres)
}

type corpsAffectations struct {
	Affectations []string `json:"affectations"`
}

func (s *Serveur) agentAffecter(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	var corps corpsAffectations
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "corps invalide"})
		return
	}
	if corps.Affectations == nil {
		corps.Affectations = []string{}
	}
	if erreur := s.Depot.Affecter(c.Request.Context(), tache.ID, corps.Affectations); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "affectation impossible, verifier les identifiants des membres"})
		return
	}
	s.repondreTache(c, tache.ID, projet)
}

type corpsEtiquetages struct {
	Etiquettes []string `json:"etiquettes"`
}

func (s *Serveur) agentEtiqueter(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	var corps corpsEtiquetages
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "corps invalide"})
		return
	}
	if corps.Etiquettes == nil {
		corps.Etiquettes = []string{}
	}
	if erreur := s.Depot.Etiqueter(c.Request.Context(), tache.ID, corps.Etiquettes); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "etiquetage impossible, verifier les identifiants des etiquettes"})
		return
	}
	s.repondreTache(c, tache.ID, projet)
}

func (s *Serveur) repondreTache(c *gin.Context, identifiant string, projet *modeles.Projet) {
	resultat, erreur := s.Depot.Tache(c.Request.Context(), identifiant)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "relecture de la tache impossible"})
		return
	}
	s.remplirURLsTache(resultat)
	s.publier("tache.modifiee", projet, resultat.Titre, resultat, resultat.Affectations, c)
	c.JSON(http.StatusOK, resultat)
}

func (s *Serveur) agentObtenirTache(c *gin.Context) {
	tache, erreur := s.Depot.Tache(c.Request.Context(), c.Param("id"))
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "tache introuvable"})
		return
	}
	c.JSON(http.StatusOK, s.remplirURLsTache(tache))
}

type corpsCommit struct {
	Commit string `json:"commit"`
}

func (s *Serveur) agentRenseignerCommit(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	var corps corpsCommit
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "corps invalide"})
		return
	}
	commit := strings.TrimSpace(corps.Commit)
	if len([]rune(commit)) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "le commit ne doit pas depasser 100 caracteres"})
		return
	}
	if erreur := s.Depot.ModifierCommit(c.Request.Context(), tache.ID, commit); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "enregistrement du commit impossible"})
		return
	}
	resultat, erreur := s.Depot.Tache(c.Request.Context(), tache.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "relecture de la tache impossible"})
		return
	}
	s.remplirURLsTache(resultat)
	s.publier("tache.modifiee", projet, resultat.Titre, resultat, resultat.Affectations, c)
	c.JSON(http.StatusOK, resultat)
}
