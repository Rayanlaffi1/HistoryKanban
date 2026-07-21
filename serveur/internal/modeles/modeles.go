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
	ID       string    `json:"id"`
	Tache    string    `json:"tache"`
	Chemin   string    `json:"chemin"`
	Nom      string    `json:"nom"`
	Taille   int64     `json:"taille"`
	URL      string    `json:"url"`
	Creation time.Time `json:"creation"`
}

type Tache struct {
	ID           string     `json:"id"`
	Projet       string     `json:"projet"`
	Colonne      string     `json:"colonne"`
	Lot          *string    `json:"lot"`
	Titre        string     `json:"titre"`
	Description  string     `json:"description"`
	Points       int        `json:"points"`
	Echeance     *time.Time `json:"echeance"`
	Position     int        `json:"position"`
	Createur     string     `json:"createur"`
	Creation     time.Time  `json:"creation"`
	Modification time.Time  `json:"modification"`
	Affectations []string   `json:"affectations"`
	Etiquettes   []string   `json:"etiquettes"`
	Images       []Image    `json:"images"`
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

type Filtre struct {
	Texte     string
	Membre    string
	Etiquette string
	Lot       string
	Echeance  string
	PointsMin *int
	PointsMax *int
}
