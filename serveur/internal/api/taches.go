package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"historykanban/serveur/internal/modeles"
)

func entierRequete(c *gin.Context, cle string) *int {
	valeur := c.Query(cle)
	if valeur == "" {
		return nil
	}
	nombre, erreur := strconv.Atoi(valeur)
	if erreur != nil {
		return nil
	}
	return &nombre
}

func (s *Serveur) listerTaches(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "lecteur")
	if !autorise {
		return
	}
	filtre := modeles.Filtre{
		Texte:     c.Query("texte"),
		Membre:    c.Query("membre"),
		Etiquette: c.Query("etiquette"),
		Lot:       c.Query("lot"),
		Echeance:  c.Query("echeance"),
		Urgence:   c.Query("urgence"),
		PointsMin: entierRequete(c, "pointsmin"),
		PointsMax: entierRequete(c, "pointsmax"),
	}
	taches, erreur := s.Depot.Taches(c.Request.Context(), projet.ID, filtre)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des taches impossible"})
		return
	}
	c.JSON(http.StatusOK, s.remplirURLs(taches))
}

type corpsTache struct {
	Titre        string     `json:"titre" binding:"required"`
	Description  string     `json:"description"`
	Colonne      string     `json:"colonne"`
	Lot          *string    `json:"lot"`
	Points       int        `json:"points"`
	Urgence      string     `json:"urgence"`
	Echeance     *time.Time `json:"echeance"`
	URLs         []string   `json:"urls"`
	Affectations []string   `json:"affectations"`
	Etiquettes   []string   `json:"etiquettes"`
}

func nettoyerURLs(valeurs []string) ([]string, string) {
	urls := []string{}
	vues := map[string]bool{}
	for _, valeur := range valeurs {
		lien := strings.TrimSpace(valeur)
		if lien == "" || vues[lien] {
			continue
		}
		if len([]rune(lien)) > 500 {
			return nil, "un lien ne doit pas depasser 500 caracteres"
		}
		analyse, erreur := url.Parse(lien)
		if erreur != nil || (analyse.Scheme != "http" && analyse.Scheme != "https") || analyse.Host == "" {
			return nil, "chaque lien doit etre une URL complete commencant par http:// ou https://"
		}
		vues[lien] = true
		urls = append(urls, lien)
	}
	if len(urls) > 20 {
		return nil, "une tache ne peut pas porter plus de 20 liens"
	}
	return urls, ""
}

var urgencesAutorisees = map[string]bool{
	"faible":  true,
	"normale": true,
	"elevee":  true,
	"urgente": true,
}

func urgenceValide(valeur string) (string, bool) {
	if valeur == "" {
		return "normale", true
	}
	if !urgencesAutorisees[valeur] {
		return "", false
	}
	return valeur, true
}

func (s *Serveur) creerTache(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "membre")
	if !autorise {
		return
	}
	var corps corpsTache
	if erreur := c.ShouldBindJSON(&corps); erreur != nil || corps.Colonne == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "titre et colonne requis"})
		return
	}
	urgence, valide := urgenceValide(corps.Urgence)
	if !valide {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "urgence invalide : faible, normale, elevee ou urgente"})
		return
	}
	tache, erreur := s.Depot.CreerTache(c.Request.Context(), modeles.Tache{
		Projet:       projet.ID,
		Colonne:      corps.Colonne,
		Lot:          corps.Lot,
		Titre:        corps.Titre,
		Description:  corps.Description,
		Points:       corps.Points,
		Urgence:      urgence,
		Echeance:     corps.Echeance,
		Createur:     s.revendications(c).Utilisateur,
		Affectations: corps.Affectations,
		Etiquettes:   corps.Etiquettes,
	})
	if erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "creation de la tache impossible : " + erreur.Error()})
		return
	}
	s.remplirURLsTache(tache)
	s.publier("tache.creee", projet, tache.Titre, tache, tache.Affectations, c)
	c.JSON(http.StatusCreated, tache)
}

func (s *Serveur) tacheAutorisee(c *gin.Context, minimum string) (*modeles.Tache, *modeles.Projet, bool) {
	tache, erreur := s.Depot.Tache(c.Request.Context(), c.Param("id"))
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "tache introuvable"})
		return nil, nil, false
	}
	projet, autorise := s.exigerRoleProjet(c, tache.Projet, minimum)
	if !autorise {
		return nil, nil, false
	}
	return tache, projet, true
}

func (s *Serveur) modifierTache(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	var corps corpsTache
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "titre requis"})
		return
	}
	contexte := c.Request.Context()
	urls := tache.URLs
	if corps.URLs != nil {
		nettoyees, message := nettoyerURLs(corps.URLs)
		if message != "" {
			c.JSON(http.StatusBadRequest, gin.H{"erreur": message})
			return
		}
		urls = nettoyees
	}
	if urls == nil {
		urls = []string{}
	}
	urgence := tache.Urgence
	if corps.Urgence != "" {
		if !urgencesAutorisees[corps.Urgence] {
			c.JSON(http.StatusBadRequest, gin.H{"erreur": "urgence invalide : faible, normale, elevee ou urgente"})
			return
		}
		urgence = corps.Urgence
	}
	if erreur := s.Depot.ModifierTache(contexte, modeles.Tache{
		ID:          tache.ID,
		Titre:       corps.Titre,
		Description: corps.Description,
		Points:      corps.Points,
		Urgence:     urgence,
		Echeance:    corps.Echeance,
		Lot:         corps.Lot,
		URLs:        urls,
	}); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "modification de la tache impossible"})
		return
	}
	if corps.Affectations != nil {
		if erreur := s.Depot.Affecter(contexte, tache.ID, corps.Affectations); erreur != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erreur": "affectation impossible : " + erreur.Error()})
			return
		}
	}
	if corps.Etiquettes != nil {
		if erreur := s.Depot.Etiqueter(contexte, tache.ID, corps.Etiquettes); erreur != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erreur": "etiquetage impossible : " + erreur.Error()})
			return
		}
	}
	resultat, erreur := s.Depot.Tache(contexte, tache.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "relecture de la tache impossible"})
		return
	}
	s.remplirURLsTache(resultat)
	s.publier("tache.modifiee", projet, resultat.Titre, resultat, resultat.Affectations, c)
	c.JSON(http.StatusOK, resultat)
}

type corpsDeplacement struct {
	Colonne  string `json:"colonne" binding:"required"`
	Position int    `json:"position"`
}

func (s *Serveur) deplacerTache(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	var corps corpsDeplacement
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "colonne requise"})
		return
	}
	if erreur := s.Depot.DeplacerTache(c.Request.Context(), tache.ID, corps.Colonne, corps.Position); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "deplacement impossible"})
		return
	}
	resultat, erreur := s.Depot.Tache(c.Request.Context(), tache.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "relecture de la tache impossible"})
		return
	}
	s.remplirURLsTache(resultat)
	s.publier("tache.deplacee", projet, resultat.Titre, resultat, resultat.Affectations, c)
	c.JSON(http.StatusOK, resultat)
}

func (s *Serveur) supprimerTache(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	if erreur := s.Depot.SupprimerTache(c.Request.Context(), tache.ID); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression de la tache impossible"})
		return
	}
	s.publier("tache.supprimee", projet, tache.Titre, gin.H{"id": tache.ID, "colonne": tache.Colonne}, tache.Affectations, c)
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}

func (s *Serveur) listerCorbeille(c *gin.Context) {
	projet, autorise := s.exigerRoleProjet(c, c.Param("id"), "lecteur")
	if !autorise {
		return
	}
	taches, erreur := s.Depot.Corbeille(c.Request.Context(), projet.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture de la corbeille impossible"})
		return
	}
	c.JSON(http.StatusOK, taches)
}

func (s *Serveur) restaurerTache(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	if tache.Suppression == nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "cette tache n'est pas dans la corbeille"})
		return
	}
	if erreur := s.Depot.RestaurerTache(c.Request.Context(), tache.ID); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "restauration impossible"})
		return
	}
	resultat, erreur := s.Depot.Tache(c.Request.Context(), tache.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "relecture de la tache impossible"})
		return
	}
	s.remplirURLsTache(resultat)
	s.publier("tache.creee", projet, resultat.Titre, resultat, resultat.Affectations, c)
	c.JSON(http.StatusOK, resultat)
}

func (s *Serveur) purgerTache(c *gin.Context) {
	tache, _, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	if tache.Suppression == nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "cette tache n'est pas dans la corbeille"})
		return
	}
	for _, image := range tache.Images {
		s.Stockage.Supprimer(c.Request.Context(), image.Chemin)
	}
	if erreur := s.Depot.PurgerTache(c.Request.Context(), tache.ID); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression definitive impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "purge"})
}

func (s *Serveur) listerActivites(c *gin.Context) {
	tache, _, autorise := s.tacheAutorisee(c, "lecteur")
	if !autorise {
		return
	}
	activites, erreur := s.Depot.Activites(c.Request.Context(), tache.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture de l'activite impossible"})
		return
	}
	c.JSON(http.StatusOK, activites)
}

func (s *Serveur) listerCommentaires(c *gin.Context) {
	tache, _, autorise := s.tacheAutorisee(c, "lecteur")
	if !autorise {
		return
	}
	commentaires, erreur := s.Depot.Commentaires(c.Request.Context(), tache.ID)
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture des commentaires impossible"})
		return
	}
	c.JSON(http.StatusOK, commentaires)
}

type corpsCommentaire struct {
	Contenu string `json:"contenu" binding:"required"`
}

func (s *Serveur) creerCommentaire(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
	if !autorise {
		return
	}
	var corps corpsCommentaire
	if erreur := c.ShouldBindJSON(&corps); erreur != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "contenu requis"})
		return
	}
	commentaire, erreur := s.Depot.CreerCommentaire(c.Request.Context(), modeles.Commentaire{
		Tache:   tache.ID,
		Auteur:  s.revendications(c).Utilisateur,
		Contenu: corps.Contenu,
		Agent:   viaAgent(c),
	})
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "creation du commentaire impossible"})
		return
	}
	s.publier("tache.commentee", projet, tache.Titre, commentaire, tache.Affectations, c)
	c.JSON(http.StatusCreated, commentaire)
}

func (s *Serveur) supprimerCommentaire(c *gin.Context) {
	if erreur := s.Depot.SupprimerCommentaire(c.Request.Context(), c.Param("id"), s.revendications(c).Utilisateur); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression du commentaire impossible"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}

var extensionsAutorisees = map[string]bool{
	"image/png":                   true,
	"image/jpeg":                  true,
	"image/gif":                   true,
	"image/webp":                  true,
	"application/pdf":             true,
	"text/plain":                  true,
	"text/csv":                    true,
	"text/markdown":               true,
	"application/json":            true,
	"application/zip":             true,
	"application/x-7z-compressed": true,
	"application/vnd.rar":         true,
	"application/gzip":            true,
	"application/msword":          true,
	"application/vnd.ms-excel":    true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	"application/vnd.oasis.opendocument.text":                                   true,
	"application/vnd.oasis.opendocument.spreadsheet":                            true,
}

func (s *Serveur) televerserImage(c *gin.Context) {
	tache, projet, autorise := s.tacheAutorisee(c, "membre")
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
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "format de fichier non pris en charge"})
		return
	}
	contenu, erreur := fichier.Open()
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "lecture du fichier impossible"})
		return
	}
	defer contenu.Close()
	chemin := tache.ID + "/" + uuid.NewString()
	if erreur := s.Stockage.Televerser(c.Request.Context(), chemin, contenu, fichier.Size, typeContenu); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "televersement impossible"})
		return
	}
	image, erreur := s.Depot.AjouterImage(c.Request.Context(), modeles.Image{
		Tache:       tache.ID,
		Chemin:      chemin,
		Nom:         fichier.Filename,
		Taille:      fichier.Size,
		TypeContenu: typeContenu,
	})
	if erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "enregistrement de l'image impossible"})
		return
	}
	image.URL = s.Stockage.URL(image.Chemin)
	s.publier("tache.image.ajoutee", projet, tache.Titre, image, tache.Affectations, c)
	c.JSON(http.StatusCreated, image)
}

func (s *Serveur) supprimerImage(c *gin.Context) {
	image, erreur := s.Depot.Image(c.Request.Context(), c.Param("id"))
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "image introuvable"})
		return
	}
	identifiantProjet, erreur := s.Depot.ProjetTache(c.Request.Context(), image.Tache)
	if erreur != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": "tache introuvable"})
		return
	}
	projet, autorise := s.exigerRoleProjet(c, identifiantProjet, "membre")
	if !autorise {
		return
	}
	s.Stockage.Supprimer(c.Request.Context(), image.Chemin)
	if erreur := s.Depot.SupprimerImage(c.Request.Context(), image.ID); erreur != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "suppression de l'image impossible"})
		return
	}
	s.publier("tache.image.supprimee", projet, "", gin.H{"id": image.ID, "tache": image.Tache}, nil, c)
	c.JSON(http.StatusOK, gin.H{"etat": "supprime"})
}
