package tempsreel

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Connexion struct {
	prise       *websocket.Conn
	Utilisateur string
	Session     string
	projets     map[string]bool
	envoi       chan []byte
	verrou      sync.Mutex
}

type Concentrateur struct {
	verrou     sync.RWMutex
	connexions map[*Connexion]bool
}

func Nouveau() *Concentrateur {
	return &Concentrateur{connexions: map[*Connexion]bool{}}
}

func (c *Concentrateur) Ajouter(prise *websocket.Conn, utilisateur, session string) *Connexion {
	connexion := &Connexion{
		prise:       prise,
		Utilisateur: utilisateur,
		Session:     session,
		projets:     map[string]bool{},
		envoi:       make(chan []byte, 64),
	}
	c.verrou.Lock()
	c.connexions[connexion] = true
	c.verrou.Unlock()
	go connexion.ecrire()
	return connexion
}

func (c *Concentrateur) Retirer(connexion *Connexion) {
	c.verrou.Lock()
	if c.connexions[connexion] {
		delete(c.connexions, connexion)
		close(connexion.envoi)
	}
	c.verrou.Unlock()
}

func (connexion *Connexion) ecrire() {
	for message := range connexion.envoi {
		if erreur := connexion.prise.WriteMessage(websocket.TextMessage, message); erreur != nil {
			return
		}
	}
	connexion.prise.Close()
}

func (connexion *Connexion) transmettre(message []byte) {
	select {
	case connexion.envoi <- message:
	default:
	}
}

func (connexion *Connexion) Abonner(projet string) {
	connexion.verrou.Lock()
	connexion.projets[projet] = true
	connexion.verrou.Unlock()
}

func (connexion *Connexion) Desabonner(projet string) {
	connexion.verrou.Lock()
	delete(connexion.projets, projet)
	connexion.verrou.Unlock()
}

func (connexion *Connexion) abonne(projet string) bool {
	connexion.verrou.Lock()
	defer connexion.verrou.Unlock()
	return connexion.projets[projet]
}

func serialiser(message any) []byte {
	octets, erreur := json.Marshal(message)
	if erreur != nil {
		return nil
	}
	return octets
}

func (c *Concentrateur) DiffuserProjet(projet string, message any) {
	octets := serialiser(message)
	if octets == nil {
		return
	}
	c.verrou.RLock()
	defer c.verrou.RUnlock()
	for connexion := range c.connexions {
		if connexion.abonne(projet) {
			connexion.transmettre(octets)
		}
	}
}

func (c *Concentrateur) EnvoyerUtilisateur(utilisateur string, message any) {
	octets := serialiser(message)
	if octets == nil {
		return
	}
	c.verrou.RLock()
	defer c.verrou.RUnlock()
	for connexion := range c.connexions {
		if connexion.Utilisateur == utilisateur {
			connexion.transmettre(octets)
		}
	}
}

func (c *Concentrateur) RemplacerSession(utilisateur, sessionActuelle string) {
	octets := serialiser(map[string]string{"type": "session.remplacee"})
	c.verrou.RLock()
	anciennes := []*Connexion{}
	for connexion := range c.connexions {
		if connexion.Utilisateur == utilisateur && connexion.Session != sessionActuelle {
			connexion.transmettre(octets)
			anciennes = append(anciennes, connexion)
		}
	}
	c.verrou.RUnlock()
	for _, connexion := range anciennes {
		c.Retirer(connexion)
	}
}
