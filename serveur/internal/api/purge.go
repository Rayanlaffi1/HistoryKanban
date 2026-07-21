package api

import (
	"context"
	"log"
	"time"
)

const joursRetentionCorbeille = 30

func (s *Serveur) DemarrerPurgeCorbeille() {
	go func() {
		for {
			s.purgerCorbeille()
			time.Sleep(6 * time.Hour)
		}
	}()
}

func (s *Serveur) purgerCorbeille() {
	contexte, annuler := context.WithTimeout(context.Background(), time.Minute)
	defer annuler()
	identifiants, chemins, erreur := s.Depot.TachesExpirees(contexte, joursRetentionCorbeille)
	if erreur != nil {
		log.Printf("purge de la corbeille impossible : %v", erreur)
		return
	}
	if len(identifiants) == 0 {
		return
	}
	for _, chemin := range chemins {
		s.Stockage.Supprimer(contexte, chemin)
	}
	for _, identifiant := range identifiants {
		if erreur := s.Depot.PurgerTache(contexte, identifiant); erreur != nil {
			log.Printf("purge de la tache %s impossible : %v", identifiant, erreur)
		}
	}
	log.Printf("corbeille purgee : %d tache(s) supprimee(s) depuis plus de %d jours",
		len(identifiants), joursRetentionCorbeille)
}
