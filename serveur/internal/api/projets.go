package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"historykanban/serveur/internal/modeles"
)

func (s *Serveur) listerProjets(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "lecteur") {
		return
	}
	projets, erreur := s.Depot.ProjetsParGroupe(c.Request.Context(), identifiant)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des projets impossible"})
		return
	}
	c.JSON(http.StatusOK, projets)
}

type corpsProjet struct {
	Nom         string `json:"nom" binding:"required"`
	Description string `json:"description"`
	Couleur     string `json:"couleur"`
	Archive     bool   `json:"archive"`
}

func (s *Serveur) creerProjet(c *gin.Context) {
	identifiant := c.Param("id")
	if !s.exigerRoleGroupe(c, identifiant, "membre") {
		return
	}
	var corps corpsProjet
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if corps.Couleur == "" {
		corps.Couleur = "#737373"
	}
	projet, erreur := s.Depot.CreerProjet(c.Request.Context(), modeles.Projet{
		Groupe:      identifiant,
		Nom:         corps.Nom,
		Description: corps.Description,
		Couleur:     corps.Couleur,
		Createur:    s.revendications(c).Utilisateur,
	})
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "creation du projet impossible"})
		return
	}
	c.JSON(http.StatusCreated, projet)
}

func (s *Serveur) obtenirProjet(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "lecteur")
	if !autorise {
		return
	}
	contexte := c.Request.Context()
	colonnes, erreur := s.Depot.Colonnes(contexte, projet.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des colonnes impossible"})
		return
	}
	etiquettes, erreur := s.Depot.Etiquettes(contexte, projet.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des etiquettes impossible"})
		return
	}
	lots, erreur := s.Depot.Lots(contexte, projet.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des lots impossible"})
		return
	}
	membres, erreur := s.Depot.Membres(contexte, projet.Groupe)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des membres impossible"})
		return
	}
	role, _ := s.Depot.RoleGroupe(contexte, projet.Groupe, s.revendications(c).Utilisateur)
	c.JSON(http.StatusOK, gin.H{
		"projet":     projet,
		"colonnes":   colonnes,
		"etiquettes": etiquettes,
		"lots":       lots,
		"membres":    membres,
		"role":       role,
	})
}

func (s *Serveur) modifierProjet(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "administrateur")
	if !autorise {
		return
	}
	var corps corpsProjet
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if corps.Couleur == "" {
		corps.Couleur = projet.Couleur
	}
	if erreur := s.Depot.ModifierProjet(c.Request.Context(), projet.ID, corps.Nom, corps.Description, corps.Couleur, corps.Archive); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "modification du projet impossible"})
		return
	}
	s.publier("projet.modifie", projet, corps.Nom, corps, nil, c)
	c.JSON(http.StatusOK, gin.H{"etat": "modifie"})
}

func (s *Serveur) supprimerProjet(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "administrateur")
	if !autorise {
		return
	}
	if erreur := s.Depot.SupprimerProjet(c.Request.Context(), projet.ID); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression du projet impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}

type corpsColonne struct {
	Nom     string `json:"nom" binding:"required"`
	Couleur string `json:"couleur"`
	Limite  *int   `json:"limite"`
}

func (s *Serveur) creerColonne(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "administrateur")
	if !autorise {
		return
	}
	var corps corpsColonne
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if corps.Couleur == "" {
		corps.Couleur = "#a3a3a3"
	}
	colonne, erreur := s.Depot.CreerColonne(c.Request.Context(), modeles.Colonne{
		Projet:  projet.ID,
		Nom:     corps.Nom,
		Couleur: corps.Couleur,
		Limite:  corps.Limite,
	})
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "creation de la colonne impossible"})
		return
	}
	s.publier("colonne.creee", projet, colonne.Nom, colonne, nil, c)
	c.JSON(http.StatusCreated, colonne)
}

func (s *Serveur) modifierColonne(c *gin.Context) {
	identifiant := c.Param("id")
	identifiantProjet, erreur := s.Depot.ProjetColonne(c.Request.Context(), identifiant)
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "colonne introuvable"})
		return
	}
	projet, autorise := s.exigerRoleProjet(c, identifiantProjet, "administrateur")
	if !autorise {
		return
	}
	var corps corpsColonne
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if corps.Couleur == "" {
		corps.Couleur = "#a3a3a3"
	}
	if erreur := s.Depot.ModifierColonne(c.Request.Context(), identifiant, corps.Nom, corps.Couleur, corps.Limite); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "modification de la colonne impossible"})
		return
	}
	s.publier("colonne.modifiee", projet, corps.Nom, corps, nil, c)
	c.JSON(http.StatusOK, gin.H{"etat": "modifie"})
}

func (s *Serveur) supprimerColonne(c *gin.Context) {
	identifiant := c.Param("id")
	identifiantProjet, erreur := s.Depot.ProjetColonne(c.Request.Context(), identifiant)
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "colonne introuvable"})
		return
	}
	projet, autorise := s.exigerRoleProjet(c, identifiantProjet, "administrateur")
	if !autorise {
		return
	}
	if erreur := s.Depot.SupprimerColonne(c.Request.Context(), identifiant); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression de la colonne impossible"})
		return
	}
	s.publier("colonne.supprimee", projet, "", gin.H{"id": identifiant}, nil, c)
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}

type corpsOrdre struct {
	Ordre []string `json:"ordre" binding:"required"`
}

func (s *Serveur) reordonnerColonnes(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "administrateur")
	if !autorise {
		return
	}
	var corps corpsOrdre
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "ordre requis"})
		return
	}
	if erreur := s.Depot.ReordonnerColonnes(c.Request.Context(), projet.ID, corps.Ordre); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "reordonnancement impossible"})
		return
	}
	s.publier("colonne.reordonnee", projet, "", corps.Ordre, nil, c)
	c.JSON(http.StatusOK, gin.H{"etat": "modifie"})
}

func (s *Serveur) televerserFichier(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "membre")
	if !autorise {
		return
	}
	fichier, erreur := c.FormFile("fichier")
	if erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "fichier requis"})
		return
	}
	if fichier.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "fichier trop volumineux"})
		return
	}
	typeContenu := fichier.Header.Get("Content-Type")
	if !extensionsAutorisees[typeContenu] {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "format d'image non pris en charge"})
		return
	}
	contenu, erreur := fichier.Open()
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture du fichier impossible"})
		return
	}
	defer contenu.Close()
	chemin := "libre/" + projet.ID + "/" + uuid.NewString()
	if erreur := s.Stockage.Televerser(c.Request.Context(), chemin, contenu, fichier.Size, typeContenu); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "televersement impossible"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"url": s.Stockage.URL(chemin)})
}

type corpsEtiquette struct {
	Nom     string `json:"nom" binding:"required"`
	Couleur string `json:"couleur"`
}

func (s *Serveur) creerEtiquette(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "membre")
	if !autorise {
		return
	}
	var corps corpsEtiquette
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if corps.Couleur == "" {
		corps.Couleur = "#737373"
	}
	etiquette, erreur := s.Depot.CreerEtiquette(c.Request.Context(), modeles.Etiquette{
		Projet:  projet.ID,
		Nom:     corps.Nom,
		Couleur: corps.Couleur,
	})
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "creation de l'etiquette impossible"})
		return
	}
	s.publier("etiquette.creee", projet, etiquette.Nom, etiquette, nil, c)
	c.JSON(http.StatusCreated, etiquette)
}

func (s *Serveur) modifierEtiquette(c *gin.Context) {
	identifiantProjet, erreur := s.Depot.ProjetEtiquette(c.Request.Context(), c.Param("id"))
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "etiquette introuvable"})
		return
	}
	if _, autorise := s.exigerRoleProjet(c, identifiantProjet, "membre"); !autorise {
		return
	}
	var corps corpsEtiquette
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if erreur := s.Depot.ModifierEtiquette(c.Request.Context(), c.Param("id"), corps.Nom, corps.Couleur); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "modification de l'etiquette impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "modifie"})
}

func (s *Serveur) supprimerEtiquette(c *gin.Context) {
	identifiantProjet, erreur := s.Depot.ProjetEtiquette(c.Request.Context(), c.Param("id"))
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "etiquette introuvable"})
		return
	}
	if _, autorise := s.exigerRoleProjet(c, identifiantProjet, "membre"); !autorise {
		return
	}
	if erreur := s.Depot.SupprimerEtiquette(c.Request.Context(), c.Param("id")); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression de l'etiquette impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}

type corpsLot struct {
	Nom      string     `json:"nom" binding:"required"`
	Couleur  string     `json:"couleur"`
	Echeance *time.Time `json:"echeance"`
}

func (s *Serveur) creerLot(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "membre")
	if !autorise {
		return
	}
	var corps corpsLot
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if corps.Couleur == "" {
		corps.Couleur = "#525252"
	}
	lot, erreur := s.Depot.CreerLot(c.Request.Context(), modeles.Lot{
		Projet:   projet.ID,
		Nom:      corps.Nom,
		Couleur:  corps.Couleur,
		Echeance: corps.Echeance,
	})
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "creation du lot impossible"})
		return
	}
	s.publier("lot.cree", projet, lot.Nom, lot, nil, c)
	c.JSON(http.StatusCreated, lot)
}

func (s *Serveur) modifierLot(c *gin.Context) {
	identifiantProjet, erreur := s.Depot.ProjetLot(c.Request.Context(), c.Param("id"))
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "lot introuvable"})
		return
	}
	if _, autorise := s.exigerRoleProjet(c, identifiantProjet, "membre"); !autorise {
		return
	}
	var corps corpsLot
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "nom requis"})
		return
	}
	if corps.Couleur == "" {
		corps.Couleur = "#525252"
	}
	if erreur := s.Depot.ModifierLot(c.Request.Context(), modeles.Lot{
		ID:       c.Param("id"),
		Nom:      corps.Nom,
		Couleur:  corps.Couleur,
		Echeance: corps.Echeance,
	}); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "modification du lot impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "modifie"})
}

func (s *Serveur) supprimerLot(c *gin.Context) {
	identifiantProjet, erreur := s.Depot.ProjetLot(c.Request.Context(), c.Param("id"))
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "lot introuvable"})
		return
	}
	if _, autorise := s.exigerRoleProjet(c, identifiantProjet, "membre"); !autorise {
		return
	}
	if erreur := s.Depot.SupprimerLot(c.Request.Context(), c.Param("id")); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression du lot impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}
