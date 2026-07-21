package main

import (
	"context"
	"log"

	"historykanban/serveur/internal/api"
	"historykanban/serveur/internal/basededonnees"
	"historykanban/serveur/internal/configuration"
	"historykanban/serveur/internal/courriels"
	"historykanban/serveur/internal/depots"
	"historykanban/serveur/internal/distribution"
	"historykanban/serveur/internal/evenements"
	"historykanban/serveur/internal/securite"
	"historykanban/serveur/internal/stockage"
	"historykanban/serveur/internal/tempsreel"
)

func main() {
	config := configuration.Charger()
	contexte := context.Background()

	reserve, erreur := basededonnees.Connecter(contexte, config.BDURL)
	if erreur != nil {
		log.Fatalf("base de donnees : %v", erreur)
	}
	defer reserve.Close()

	verificateur, erreur := securite.NouveauVerificateur(config.URLJWKS(), config.Emetteur())
	if erreur != nil {
		log.Fatalf("keycloak : %v", erreur)
	}

	entrepot, erreur := stockage.Nouveau(config.MinioHote, config.MinioCle, config.MinioSecret, config.MinioSeau, config.MinioURLPublique)
	if erreur != nil {
		log.Fatalf("minio : %v", erreur)
	}

	bus, erreur := evenements.Connecter(config.RabbitURL)
	if erreur != nil {
		log.Fatalf("rabbitmq : %v", erreur)
	}
	defer bus.Fermer()

	depot := depots.Nouveau(reserve)
	concentrateur := tempsreel.Nouveau()
	envoyeur := courriels.Nouveau(config.SMTPHote, config.SMTPPort, config.SMTPExpediteur)

	distributeur := &distribution.Distributeur{
		Depot:         depot,
		Concentrateur: concentrateur,
		Envoyeur:      envoyeur,
	}
	if erreur := distributeur.Demarrer(bus); erreur != nil {
		log.Fatalf("distribution : %v", erreur)
	}

	serveur := &api.Serveur{
		Config:        config,
		Depot:         depot,
		Verificateur:  verificateur,
		Concentrateur: concentrateur,
		Bus:           bus,
		Stockage:      entrepot,
	}

	serveur.DemarrerPurgeCorbeille()

	log.Printf("serveur HistoryKanban demarre sur le port %s", config.Port)
	if erreur := serveur.Routeur().Run(":" + config.Port); erreur != nil {
		log.Fatalf("serveur : %v", erreur)
	}
}
