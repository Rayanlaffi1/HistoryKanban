package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var promoteur = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type messageClient struct {
	Action string `json:"action"`
	Projet string `json:"projet"`
}

func (s *Serveur) websocket(c *gin.Context) {
	jeton := c.Query("jeton")
	revendications, erreur := s.Verificateur.Verifier(jeton)
	if erreur != nil {
		log.Printf("verification du jeton websocket echouee : %v", erreur)
		c.JSON(http.StatusUnauthorized, gin.H{"erreur": "jeton invalide"})
		return
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
