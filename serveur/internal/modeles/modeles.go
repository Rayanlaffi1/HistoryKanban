package modeles

import (
	"encoding/json"
	"time"
)

type Utilisateur struct {
	ID       string    `json:"id"`
	Courriel string    `json:"courriel"`
	Nom      string    `json:"nom"`
	Prenom   string    `json:"prenom"`
	Creation time.Time `json:"creation"`
}

type Groupe struct {
	ID           string    `json:"id"`
	Nom          string    `json:"nom"`
	Description  string    `json:"description"`
	Proprietaire string    `json:"proprietaire"`
	Creation     time.Time `json:"creation"`
	Role         string    `json:"role"`
	NbMembres    int       `json:"nbmembres"`
	NbProjets    int       `json:"nbprojets"`
}

type Membre struct {
	Groupe      string    `json:"groupe"`
	Utilisateur string    `json:"utilisateur"`
	Role        string    `json:"role"`
	Fonction    string    `json:"fonction"`
	Ajout       time.Time `json:"ajout"`
	Courriel    string    `json:"courriel"`
	Nom         string    `json:"nom"`
	Prenom      string    `json:"prenom"`
}

type Projet struct {
	ID          string    `json:"id"`
	Groupe      string    `json:"groupe"`
	Nom         string    `json:"nom"`
	Description string    `json:"description"`
	Couleur     string    `json:"couleur"`
	Archive     bool      `json:"archive"`
	Depot       string    `json:"depot"`
	Createur    string    `json:"createur"`
	Creation    time.Time `json:"creation"`
	NbTaches    int       `json:"nbtaches"`
}

type Colonne struct {
	ID       string `json:"id"`
	Projet   string `json:"projet"`
	Nom      string `json:"nom"`
	Couleur  string `json:"couleur"`
	Position int    `json:"position"`
	Limite   *int   `json:"limite"`
}

type Lot struct {
	ID       string     `json:"id"`
	Projet   string     `json:"projet"`
	Nom      string     `json:"nom"`
	Couleur  string     `json:"couleur"`
	Echeance *time.Time `json:"echeance"`
	Creation time.Time  `json:"creation"`
}

type Etiquette struct {
	ID      string `json:"id"`
	Projet  string `json:"projet"`
	Nom     string `json:"nom"`
	Couleur string `json:"couleur"`
}

type Image struct {
	ID          string    `json:"id"`
	Tache       string    `json:"tache"`
	Chemin      string    `json:"chemin"`
	Nom         string    `json:"nom"`
	Taille      int64     `json:"taille"`
	TypeContenu string    `json:"typecontenu"`
	URL         string    `json:"url"`
	Creation    time.Time `json:"creation"`
}

type Tache struct {
	ID           string     `json:"id"`
	Projet       string     `json:"projet"`
	Colonne      string     `json:"colonne"`
	Lot          *string    `json:"lot"`
	Titre        string     `json:"titre"`
	Description  string     `json:"description"`
	Points       int        `json:"points"`
	Urgence      string     `json:"urgence"`
	Echeance     *time.Time `json:"echeance"`
	Commit       string     `json:"commit"`
	Position     int        `json:"position"`
	Suppression  *time.Time `json:"suppression"`
	Createur     string     `json:"createur"`
	Creation     time.Time  `json:"creation"`
	Modification time.Time  `json:"modification"`
	Affectations []string   `json:"affectations"`
	Etiquettes   []string   `json:"etiquettes"`
	Images       []Image    `json:"images"`
}

type Cle struct {
	Projet      string    `json:"projet"`
	Utilisateur string    `json:"utilisateur"`
	Cle         string    `json:"cle"`
	Creation    time.Time `json:"creation"`
}

type Commentaire struct {
	ID       string    `json:"id"`
	Tache    string    `json:"tache"`
	Auteur   string    `json:"auteur"`
	Contenu  string    `json:"contenu"`
	Creation time.Time `json:"creation"`
	Nom      string    `json:"nom"`
	Prenom   string    `json:"prenom"`
}

type Activite struct {
	ID          string    `json:"id"`
	Tache       string    `json:"tache"`
	Utilisateur string    `json:"utilisateur"`
	Type        string    `json:"type"`
	Detail      string    `json:"detail"`
	Creation    time.Time `json:"creation"`
	Nom         string    `json:"nom"`
	Prenom      string    `json:"prenom"`
}

type Notification struct {
	ID          string          `json:"id"`
	Utilisateur string          `json:"utilisateur"`
	Type        string          `json:"type"`
	Contenu     json.RawMessage `json:"contenu"`
	Lue         bool            `json:"lue"`
	Creation    time.Time       `json:"creation"`
}

type Preferences struct {
	Utilisateur string          `json:"utilisateur"`
	Courriels   bool            `json:"courriels"`
	Types       map[string]bool `json:"types"`
}

type LigneClassement struct {
	Utilisateur  string `json:"utilisateur"`
	Nom          string `json:"nom"`
	Prenom       string `json:"prenom"`
	Points       int    `json:"points"`
	Taches       int    `json:"taches"`
	Creees       int    `json:"creees"`
	Commentaires int    `json:"commentaires"`
}

type PointSerie struct {
	Periode time.Time `json:"periode"`
	Points  int       `json:"points"`
	Taches  int       `json:"taches"`
}

type TotauxStatistiques struct {
	Points      int `json:"points"`
	Terminees   int `json:"terminees"`
	EnRetard    int `json:"enretard"`
	Total       int `json:"total"`
	TotalPoints int `json:"totalpoints"`
}

type Statistiques struct {
	Classement []LigneClassement  `json:"classement"`
	Serie      []PointSerie       `json:"serie"`
	Totaux     TotauxStatistiques `json:"totaux"`
}

type Resultat struct {
	Tache         string     `json:"tache"`
	Titre         string     `json:"titre"`
	Urgence       string     `json:"urgence"`
	Echeance      *time.Time `json:"echeance"`
	Projet        string     `json:"projet"`
	ProjetNom     string     `json:"projetnom"`
	ProjetCouleur string     `json:"projetcouleur"`
	Groupe        string     `json:"groupe"`
	GroupeNom     string     `json:"groupenom"`
	Colonne       string     `json:"colonne"`
	Origine       string     `json:"origine"`
	Extrait       string     `json:"extrait"`
	Modification  time.Time  `json:"modification"`
}

type Filtre struct {
	Texte     string
	Membre    string
	Etiquette string
	Lot       string
	Echeance  string
	Urgence   string
	PointsMin *int
	PointsMax *int
}
