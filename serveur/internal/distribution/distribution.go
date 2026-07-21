package distribution

import (
	"context"
	"log"
	"time"

	"historykanban/serveur/internal/courriels"
	"historykanban/serveur/internal/depots"
	"historykanban/serveur/internal/evenements"
	"historykanban/serveur/internal/tempsreel"
)

var libelles = map[string]string{
	"tache.creee":           "Nouvelle tâche",
	"tache.modifiee":        "Tâche modifiée",
	"tache.deplacee":        "Tâche déplacée",
	"tache.supprimee":       "Tâche supprimée",
	"tache.commentee":       "Nouveau commentaire",
	"tache.image.ajoutee":   "Image ajoutée à une tâche",
	"tache.image.supprimee": "Image supprimée d'une tâche",
	"groupe.membre.ajoute":  "Ajout à un groupe",
}

type Distributeur struct {
	Depot         *depots.Depot
	Concentrateur *tempsreel.Concentrateur
	Envoyeur      *courriels.Envoyeur
}

func (d *Distributeur) Demarrer(bus *evenements.Bus) error {
	if erreur := bus.Consommer("tempsreel", "#", d.tempsreel); erreur != nil {
		return erreur
	}
	return bus.Consommer("courriels", "tache.#", d.courriels)
}

func (d *Distributeur) tempsreel(evenement evenements.Evenement) {
	contexte, annuler := context.WithTimeout(context.Background(), 10*time.Second)
	defer annuler()
	if evenement.Projet != "" {
		d.Concentrateur.DiffuserProjet(evenement.Projet, map[string]any{
			"type":      evenement.Type,
			"projet":    evenement.Projet,
			"acteur":    evenement.Acteur,
			"acteurnom": evenement.ActeurNom,
			"titre":     evenement.Titre,
			"donnees":   evenement.Donnees,
		})
	}
	for _, destinataire := range evenement.Destinataires {
		if destinataire == evenement.Acteur {
			continue
		}
		notification, erreur := d.Depot.CreerNotification(contexte, destinataire, evenement.Type, map[string]any{
			"titre":     evenement.Titre,
			"projet":    evenement.Projet,
			"acteurnom": evenement.ActeurNom,
		})
		if erreur != nil {
			log.Printf("creation de notification impossible : %v", erreur)
			continue
		}
		d.Concentrateur.EnvoyerUtilisateur(destinataire, map[string]any{
			"type":    "notification",
			"donnees": notification,
		})
	}
}

func (d *Distributeur) courriels(evenement evenements.Evenement) {
	destinataires := []string{}
	for _, destinataire := range evenement.Destinataires {
		if destinataire != evenement.Acteur {
			destinataires = append(destinataires, destinataire)
		}
	}
	if len(destinataires) == 0 {
		return
	}
	contexte, annuler := context.WithTimeout(context.Background(), 10*time.Second)
	defer annuler()
	adresses, erreur := d.Depot.CourrielsUtilisateurs(contexte, destinataires, evenement.Type)
	if erreur != nil {
		log.Printf("lecture des courriels impossible : %v", erreur)
		return
	}
	if len(adresses) == 0 {
		return
	}
	libelle := libelles[evenement.Type]
	if libelle == "" {
		libelle = "Changement sur une tâche"
	}
	corps := d.Envoyeur.GabaritChangement(evenement.Titre, libelle, evenement.ActeurNom, evenement.Projet)
	if erreur := d.Envoyeur.Envoyer(adresses, "HistoryKanban — "+libelle, corps); erreur != nil {
		log.Printf("envoi de courriel impossible : %v", erreur)
	}
}
