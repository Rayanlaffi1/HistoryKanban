package securite

import (
	"errors"
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
	cles     keyfunc.Keyfunc
	emetteur string
}

func NouveauVerificateur(urlJWKS, emetteur string) (*Verificateur, error) {
	var cles keyfunc.Keyfunc
	var erreur error
	for tentative := 0; tentative < 30; tentative++ {
		cles, erreur = keyfunc.NewDefault([]string{urlJWKS})
		if erreur == nil {
			return &Verificateur{cles: cles, emetteur: emetteur}, nil
		}
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
		jwt.WithIssuer(v.emetteur),
		jwt.WithExpirationRequired(),
	)
	if erreur != nil {
		return nil, erreur
	}
	revendications, valide := jeton.Claims.(jwt.MapClaims)
	if !valide {
		return nil, errors.New("revendications invalides")
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
