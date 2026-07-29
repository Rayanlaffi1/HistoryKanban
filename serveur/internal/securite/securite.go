package securite

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Revendications struct {
	Utilisateur string
	Courriel    string
	Nom         string
	Prenom      string
	Session     string
	Connexion   int64
	Roles       []string
}

type Verificateur struct {
	cles      keyfunc.Keyfunc
	emetteurs map[string]bool
}

func NouveauVerificateur(urlJWKS string, emetteurs []string) (*Verificateur, error) {
	acceptes := make(map[string]bool, len(emetteurs))
	for _, emetteur := range emetteurs {
		acceptes[emetteur] = true
	}
	var erreur error
	for tentative := 0; tentative < 60; tentative++ {
		var cles keyfunc.Keyfunc
		cles, erreur = keyfunc.NewDefault([]string{urlJWKS})
		if erreur == nil {
			contexte, annuler := context.WithTimeout(context.Background(), 5*time.Second)
			jeu, erreurLecture := cles.Storage().KeyReadAll(contexte)
			annuler()
			if erreurLecture == nil && len(jeu) > 0 {
				log.Printf("JWKS charge avec %d cle(s) depuis %s", len(jeu), urlJWKS)
				return &Verificateur{cles: cles, emetteurs: acceptes}, nil
			}
			erreur = errors.New("jeu de cles JWKS vide, Keycloak pas encore pret")
		}
		log.Printf("chargement du JWKS en attente : %v", erreur)
		time.Sleep(3 * time.Second)
	}
	return nil, erreur
}

func chaine(revendications jwt.MapClaims, cle string) string {
	if valeur, present := revendications[cle].(string); present {
		return valeur
	}
	return ""
}

func (v *Verificateur) Verifier(brut string) (*Revendications, error) {
	jeton, erreur := jwt.Parse(
		brut,
		v.cles.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
	)
	if erreur != nil {
		return nil, erreur
	}
	revendications, valide := jeton.Claims.(jwt.MapClaims)
	if !valide {
		return nil, errors.New("revendications invalides")
	}
	emetteur, erreur := revendications.GetIssuer()
	if erreur != nil || !v.emetteurs[emetteur] {
		return nil, errors.New("emetteur du jeton non autorise")
	}
	resultat := &Revendications{
		Utilisateur: chaine(revendications, "sub"),
		Courriel:    chaine(revendications, "email"),
		Nom:         chaine(revendications, "family_name"),
		Prenom:      chaine(revendications, "given_name"),
		Session:     chaine(revendications, "sid"),
	}
	if resultat.Utilisateur == "" {
		return nil, errors.New("identifiant absent du jeton")
	}
	if moment, present := revendications["auth_time"].(float64); present {
		resultat.Connexion = int64(moment)
	} else if moment, present := revendications["iat"].(float64); present {
		resultat.Connexion = int64(moment)
	}
	if acces, present := revendications["realm_access"].(map[string]any); present {
		if roles, present := acces["roles"].([]any); present {
			for _, role := range roles {
				if nom, correct := role.(string); correct {
					resultat.Roles = append(resultat.Roles, nom)
				}
			}
		}
	}
	return resultat, nil
}
