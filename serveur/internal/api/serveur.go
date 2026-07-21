package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"historykanban/serveur/internal/configuration"
	"historykanban/serveur/internal/depots"
	"historykanban/serveur/internal/evenements"
	"historykanban/serveur/internal/modeles"
	"historykanban/serveur/internal/securite"
	"historykanban/serveur/internal/stockage"
	"historykanban/serveur/internal/tempsreel"
)

type Serveur struct {
	Config        configuration.Configuration
	Depot         *depots.Depot
	Verificateur  *securite.Verificateur
	Concentrateur *tempsreel.Concentrateur
	Bus           *evenements.Bus
	Stockage      *stockage.Stockage
}

func (s *Serveur) Routeur() *gin.Engine {
	moteur := gin.New()
	moteur.Use(gin.Logger(), gin.Recovery())
	moteur.Use(cors.New(cors.Config{
		AllowOrigins:     s.Config.Origines,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	moteur.GET("/sante", func(c *gin.Context) {
		c.JSON(200, gin.H{"etat": "operationnel"})
	})
	moteur.POST("/api/webhooks/keycloak", s.webhookKeycloak)
	moteur.GET("/ws", s.websocket)

	api := moteur.Group("/api", s.authentifier())

	api.GET("/moi", s.moi)
	api.GET("/presence", s.presence)
	api.GET("/moi/preferences", s.obtenirPreferences)
	api.PUT("/moi/preferences", s.enregistrerPreferences)

	api.GET("/notifications", s.listerNotifications)
	api.PUT("/notifications/tout", s.toutMarquerLu)
	api.PUT("/notifications/:id/lue", s.marquerLue)

	api.GET("/groupes", s.listerGroupes)
	api.POST("/groupes", s.creerGroupe)
	api.PUT("/groupes/:id", s.modifierGroupe)
	api.DELETE("/groupes/:id", s.supprimerGroupe)
	api.GET("/groupes/:id/membres", s.listerMembres)
	api.POST("/groupes/:id/membres", s.ajouterMembre)
	api.PUT("/groupes/:id/membres/:utilisateur", s.modifierRoleMembre)
	api.PUT("/groupes/:id/membres/:utilisateur/fonction", s.modifierFonctionMembre)
	api.DELETE("/groupes/:id/membres/:utilisateur", s.retirerMembre)
	api.GET("/groupes/:id/statistiques", s.statistiquesGroupe)
	api.GET("/groupes/:id/projets", s.listerProjets)
	api.POST("/groupes/:id/projets", s.creerProjet)

	api.GET("/projets/:id", s.obtenirProjet)
	api.PUT("/projets/:id", s.modifierProjet)
	api.DELETE("/projets/:id", s.supprimerProjet)
	api.POST("/projets/:id/colonnes", s.creerColonne)
	api.PUT("/projets/:id/colonnes/ordre", s.reordonnerColonnes)
	api.POST("/projets/:id/etiquettes", s.creerEtiquette)
	api.POST("/projets/:id/lots", s.creerLot)
	api.GET("/projets/:id/taches", s.listerTaches)
	api.POST("/projets/:id/taches", s.creerTache)
	api.POST("/projets/:id/fichiers", s.televerserFichier)

	api.PUT("/colonnes/:id", s.modifierColonne)
	api.DELETE("/colonnes/:id", s.supprimerColonne)
	api.PUT("/etiquettes/:id", s.modifierEtiquette)
	api.DELETE("/etiquettes/:id", s.supprimerEtiquette)
	api.PUT("/lots/:id", s.modifierLot)
	api.DELETE("/lots/:id", s.supprimerLot)

	api.PUT("/taches/:id", s.modifierTache)
	api.DELETE("/taches/:id", s.supprimerTache)
	api.PUT("/taches/:id/deplacer", s.deplacerTache)
	api.GET("/taches/:id/activites", s.listerActivites)
	api.GET("/taches/:id/commentaires", s.listerCommentaires)
	api.POST("/taches/:id/commentaires", s.creerCommentaire)
	api.DELETE("/commentaires/:id", s.supprimerCommentaire)
	api.POST("/taches/:id/images", s.televerserImage)
	api.DELETE("/images/:id", s.supprimerImage)

	return moteur
}

func (s *Serveur) remplirURLs(taches []modeles.Tache) []modeles.Tache {
	for indiceTache := range taches {
		for indiceImage := range taches[indiceTache].Images {
			taches[indiceTache].Images[indiceImage].URL = s.Stockage.URL(taches[indiceTache].Images[indiceImage].Chemin)
		}
	}
	return taches
}

func (s *Serveur) remplirURLsTache(tache *modeles.Tache) *modeles.Tache {
	for indice := range tache.Images {
		tache.Images[indice].URL = s.Stockage.URL(tache.Images[indice].Chemin)
	}
	return tache
}
