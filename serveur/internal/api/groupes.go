package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Serveur) listerGroupes(c *gin.Context) {
	groupes, erreur := s.Depot.GroupesParUtilisateur(c.Request.Context(), s.revendications(c).Utilisateur)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des groupes impossible"})
		return
	}
	c.JSON(http.StatusOK, groupes)
}

type corpsGroupe struct {
	Nom         string `json:"nom" binding:"required"`
	Description string `json:"description"`
}

func (s *Serveur) creerGroupe(c *gin.Context) {
	var corps corpsGroupe
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	groupe, erreur := s.Depot.CreerGroupe(c.Request.Context(), corps.Nom, corps.Description, s.revendications(c).Utilisateur)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "creation du groupe impossible"})
		return
	}
	c.JSON(http.StatusCreated, groupe)
}

func (s *Serveur) modifierGroupe(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "administrateur") {
		return
	}
	var corps corpsGroupe
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if erreur := s.Depot.ModifierGroupe(c.Request.Context(), identifiant, corps.Nom, corps.Description); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "modification du groupe impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "modifie"})
}

func (s *Serveur) supprimerGroupe(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "proprietaire") {
		return
	}
	if erreur := s.Depot.SupprimerGroupe(c.Request.Context(), identifiant); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression du groupe impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}

func (s *Serveur) listerMembres(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "lecteur") {
		return
	}
	membres, erreur := s.Depot.Membres(c.Request.Context(), identifiant)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des membres impossible"})
		return
	}
	c.JSON(http.StatusOK, membres)
}

type corpsMembre struct {
	Courriel string `json:"courriel" binding:"required"`
	Role     string `json:"role"`
}

func (s *Serveur) ajouterMembre(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "administrateur") {
		return
	}
	var corps corpsMembre
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "courriel requis"})
		return
	}
	if corps.Role == "" {
		corps.Role = "membre"
	}
	if corps.Role == "proprietaire" || niveaux[corps.Role] == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "role invalide"})
		return
	}
	utilisateur, erreur := s.Depot.UtilisateurParCourriel(c.Request.Context(), corps.Courriel)
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "aucun utilisateur avec ce courriel"})
		return
	}
	if erreur := s.Depot.AjouterMembre(c.Request.Context(), identifiant, utilisateur.ID, corps.Role); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "ajout du membre impossible"})
		return
	}
	groupe, _ := s.Depot.Groupe(c.Request.Context(), identifiant)
	if groupe != nil {
		s.publier("groupe.membre.ajoute", nil, groupe.Nom, groupe, []string{utilisateur.ID}, c)
	}
	c.JSON(http.StatusCreated, gin.H{"etat": "ajoute"})
}

type corpsRole struct {
	Role string `json:"role" binding:"required"`
}

func (s *Serveur) modifierRoleMembre(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "proprietaire") {
		return
	}
	var corps corpsRole
	if erreur := c.ShouldBindJSON(&corps); erreur != nil || corps.Role == "proprietaire" || niveaux[corps.Role] == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "role invalide"})
		return
	}
	cible := c.Param("utilisateur")
	groupe, erreur := s.Depot.Groupe(c.Request.Context(), identifiant)
	if erreur != nil || groupe.Proprietaire == cible {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "le role du proprietaire ne peut pas changer"})
		return
	}
	if erreur := s.Depot.ModifierRole(c.Request.Context(), identifiant, cible, corps.Role); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "modification du role impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "modifie"})
}

func (s *Serveur) retirerMembre(c *gin.Context) {
	identifiant := c.Param("id")
	cible := c.Param("utilisateur")
	if cible != s.revendications(c).Utilisateur {
		if !s.exigerRoleGroupe(c, identifiant, "administrateur") {
			return
		}
	} else if !s.exigerRoleGroupe(c, identifiant, "lecteur") {
		return
	}
	if erreur := s.Depot.RetirerMembre(c.Request.Context(), identifiant, cible); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "retrait du membre impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "retire"})
}
