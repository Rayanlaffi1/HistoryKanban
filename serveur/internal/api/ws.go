package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type messageClient struct {
	Action string `json:"action"`
	Projet string `json:"projet"`
}

func (s *Serveur) origineWebSocketAutorisee(r *http.Request) bool {
	origine := r.Header.Get("Origin")
	if origine == "" {
		return true
	}
	for _, autorisee := range s.Config.Origines {
		if strings.TrimRight(strings.TrimSpace(autorisee), "/") == strings.TrimRight(origine, "/") {
			return true
		}
	}
	return false
}

func jetonWebSocket(c *gin.Context) string {
	protocoles := websocket.Subprotocols(c.Request)
	for _, protocole := range protocoles {
		if strings.HasPrefix(protocole, "bearer.") {
			return strings.TrimPrefix(protocole, "bearer.")
		}
	}
	if entete := c.GetHeader("Authorization"); strings.HasPrefix(entete, "Bearer ") {
		return strings.TrimPrefix(entete, "Bearer ")
	}
	return c.Query("jeton")
}

func (s *Serveur) websocket(c *gin.Context) {
	jeton := jetonWebSocket(c)
	revendications, erreur := s.Verificateur.Verifier(jeton)
	if erreur != nil {
		log.Printf("verification du jeton websocket echouee : %v", erreur)
		c.JSON(http.StatusUnauthorized, gin.H{"erreur": "jeton invalide"})
		return
	}
	promoteur := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     s.origineWebSocketAutorisee,
		Subprotocols:    []string{"historykanban"},
	}
	prise, erreur := promoteur.Upgrade(c.Writer, c.Request, nil)
	if erreur != nil {
		return
	}
	if revendications.Session != "" {
		valide, remplacee, erreur := s.Depot.ReclamerSession(
			c.Request.Context(), revendications.Utilisateur, revendications.Session, revendications.Connexion)
		if erreur != nil || !valide {
			prise.WriteMessage(websocket.TextMessage, []byte(`{"type":"session.remplacee"}`))
			prise.Close()
			return
		}
		if remplacee {
			s.Concentrateur.RemplacerSession(revendications.Utilisateur, revendications.Session)
		}
	}
	connexion := s.Concentrateur.Ajouter(prise, revendications.Utilisateur, revendications.Session)
	defer s.Concentrateur.Retirer(connexion)
	prise.SetReadDeadline(time.Now().Add(90 * time.Second))
	prise.SetPongHandler(func(string) error {
		prise.SetReadDeadline(time.Now().Add(90 * time.Second))
		return nil
	})
	verification := time.NewTicker(time.Minute)
	defer verification.Stop()
	go func() {
		for range verification.C {
			revendicationsActuelles, erreur := s.Verificateur.Verifier(jeton)
			if erreur != nil || revendicationsActuelles.Utilisateur != revendications.Utilisateur {
				prise.WriteMessage(websocket.TextMessage, []byte(`{"type":"session.remplacee"}`))
				prise.Close()
				return
			}
			if revendications.Session != "" {
				valide, _, erreur := s.Depot.ReclamerSession(c.Request.Context(), revendications.Utilisateur, revendications.Session, revendications.Connexion)
				if erreur != nil || !valide {
					prise.WriteMessage(websocket.TextMessage, []byte(`{"type":"session.remplacee"}`))
					prise.Close()
					return
				}
			}
			prise.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second))
		}
	}()
	for {
		_, brut, erreur := prise.ReadMessage()
		if erreur != nil {
			return
		}
		var message messageClient
		if erreur := json.Unmarshal(brut, &message); erreur != nil {
			continue
		}
		switch message.Action {
		case "abonner":
			projet, erreur := s.Depot.Projet(c.Request.Context(), message.Projet)
			if erreur != nil {
				continue
			}
			role, erreur := s.Depot.RoleGroupe(c.Request.Context(), projet.Groupe, revendications.Utilisateur)
			if erreur != nil || role == "" {
				continue
			}
			connexion.Abonner(message.Projet)
		case "desabonner":
			connexion.Desabonner(message.Projet)
		}
	}
}
