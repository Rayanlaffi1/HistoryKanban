package depots

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"historykanban/serveur/internal/modeles"
)

type Depot struct {
	bd *pgxpool.Pool
}

func Nouveau(bd *pgxpool.Pool) *Depot {
	return &Depot{bd: bd}
}

func (d *Depot) AssurerUtilisateur(ctx context.Context, utilisateur modeles.Utilisateur) error {
	_, erreur := d.bd.Exec(ctx, `
		INSERT INTO utilisateurs (id, courriel, nom, prenom)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET courriel = EXCLUDED.courriel, nom = EXCLUDED.nom, prenom = EXCLUDED.prenom`,
		utilisateur.ID, utilisateur.Courriel, utilisateur.Nom, utilisateur.Prenom)
	return erreur
}

func (d *Depot) UtilisateurParCourriel(ctx context.Context, courriel string) (*modeles.Utilisateur, error) {
	var utilisateur modeles.Utilisateur
	erreur := d.bd.QueryRow(ctx,
		`SELECT id, courriel, nom, prenom, creation FROM utilisateurs WHERE lower(courriel) = lower($1)`,
		courriel).Scan(&utilisateur.ID, &utilisateur.Courriel, &utilisateur.Nom, &utilisateur.Prenom, &utilisateur.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &utilisateur, nil
}

func (d *Depot) Utilisateurs(ctx context.Context, identifiants []string) ([]modeles.Utilisateur, error) {
	lignes, erreur := d.bd.Query(ctx,
		`SELECT id, courriel, nom, prenom, creation FROM utilisateurs WHERE id = ANY($1)`, identifiants)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	utilisateurs := []modeles.Utilisateur{}
	for lignes.Next() {
		var utilisateur modeles.Utilisateur
		if erreur := lignes.Scan(&utilisateur.ID, &utilisateur.Courriel, &utilisateur.Nom, &utilisateur.Prenom, &utilisateur.Creation); erreur != nil {
			return nil, erreur
		}
		utilisateurs = append(utilisateurs, utilisateur)
	}
	return utilisateurs, lignes.Err()
}

func (d *Depot) ReclamerSession(ctx context.Context, utilisateur, session string, connexion int64) (bool, bool, error) {
	var sessionActuelle string
	var connexionActuelle int64
	erreur := d.bd.QueryRow(ctx,
		`SELECT session, connexion FROM sessions WHERE utilisateur = $1`, utilisateur).
		Scan(&sessionActuelle, &connexionActuelle)
	if erreur != nil {
		_, erreur = d.bd.Exec(ctx, `
			INSERT INTO sessions (utilisateur, session, connexion)
			VALUES ($1, $2, $3)
			ON CONFLICT (utilisateur) DO UPDATE SET session = EXCLUDED.session, connexion = EXCLUDED.connexion, modification = now()`,
			utilisateur, session, connexion)
		return erreur == nil, false, erreur
	}
	if sessionActuelle == session {
		return true, false, nil
	}
	if connexion >= connexionActuelle {
		_, erreur = d.bd.Exec(ctx,
			`UPDATE sessions SET session = $2, connexion = $3, modification = now() WHERE utilisateur = $1`,
			utilisateur, session, connexion)
		return erreur == nil, true, erreur
	}
	return false, false, nil
}

func (d *Depot) SupprimerSession(ctx context.Context, utilisateur, session string) error {
	_, erreur := d.bd.Exec(ctx,
		`DELETE FROM sessions WHERE utilisateur = $1 AND session = $2`, utilisateur, session)
	return erreur
}
