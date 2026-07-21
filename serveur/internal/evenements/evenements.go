package evenements

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const echange = "historykanban"

type Evenement struct {
	Type          string          `json:"type"`
	Projet        string          `json:"projet"`
	Groupe        string          `json:"groupe"`
	Acteur        string          `json:"acteur"`
	ActeurNom     string          `json:"acteurnom"`
	Titre         string          `json:"titre"`
	Donnees       json.RawMessage `json:"donnees"`
	Destinataires []string        `json:"destinataires"`
	Agent         bool            `json:"agent"`
}

type Bus struct {
	connexion *amqp.Connection
	canal     *amqp.Channel
	verrou    sync.Mutex
}

func Connecter(url string) (*Bus, error) {
	var connexion *amqp.Connection
	var erreur error
	for tentative := 0; tentative < 20; tentative++ {
		connexion, erreur = amqp.Dial(url)
		if erreur == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if erreur != nil {
		return nil, erreur
	}
	canal, erreur := connexion.Channel()
	if erreur != nil {
		return nil, erreur
	}
	if erreur = canal.ExchangeDeclare(echange, "topic", true, false, false, false, nil); erreur != nil {
		return nil, erreur
	}
	return &Bus{connexion: connexion, canal: canal}, nil
}

func (b *Bus) Publier(evenement Evenement) {
	octets, erreur := json.Marshal(evenement)
	if erreur != nil {
		log.Printf("serialisation evenement impossible : %v", erreur)
		return
	}
	b.verrou.Lock()
	defer b.verrou.Unlock()
	contexte, annuler := context.WithTimeout(context.Background(), 5*time.Second)
	defer annuler()
	erreur = b.canal.PublishWithContext(contexte, echange, evenement.Type, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        octets,
	})
	if erreur != nil {
		log.Printf("publication evenement impossible : %v", erreur)
	}
}

func (b *Bus) Consommer(file, cle string, gestionnaire func(Evenement)) error {
	canal, erreur := b.connexion.Channel()
	if erreur != nil {
		return erreur
	}
	declaree, erreur := canal.QueueDeclare(file, true, false, false, false, nil)
	if erreur != nil {
		return erreur
	}
	if erreur = canal.QueueBind(declaree.Name, cle, echange, false, nil); erreur != nil {
		return erreur
	}
	livraisons, erreur := canal.Consume(declaree.Name, "", true, false, false, false, nil)
	if erreur != nil {
		return erreur
	}
	go func() {
		for livraison := range livraisons {
			var evenement Evenement
			if erreur := json.Unmarshal(livraison.Body, &evenement); erreur != nil {
				log.Printf("evenement illisible : %v", erreur)
				continue
			}
			gestionnaire(evenement)
		}
	}()
	return nil
}

func (b *Bus) Fermer() {
	b.canal.Close()
	b.connexion.Close()
}
