package depots

import (
	"context"

	"github.com/jackc/pgx/v5"

	"historykanban/serveur/internal/modeles"
)

func (d *Depot) Cle(ctx context.Context, projet, utilisateur string) (*modeles.Cle, error) {
	var cle modeles.Cle
	erreur := d.bd.QueryRow(ctx,
		`SELECT projet, utilisateur, cle, creation FROM cles WHERE projet = $1 AND utilisateur = $2`,
		projet, utilisateur).
		Scan(&cle.Projet, &cle.Utilisateur, &cle.Cle, &cle.Creation)
	if erreur == pgx.ErrNoRows {
		return nil, nil
	}
	if erreur != nil {
		return nil, erreur
	}
	return &cle, nil
}

func (d *Depot) EnregistrerCle(ctx context.Context, projet, utilisateur, valeur string) (*modeles.Cle, error) {
	var cle modeles.Cle
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO cles (projet, utilisateur, cle) VALUES ($1, $2, $3)
		ON CONFLICT (projet, utilisateur) DO UPDATE SET cle = $3, creation = now()
		RETURNING projet, utilisateur, cle, creation`,
		projet, utilisateur, valeur).
		Scan(&cle.Projet, &cle.Utilisateur, &cle.Cle, &cle.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &cle, nil
}

func (d *Depot) SupprimerCle(ctx context.Context, projet, utilisateur string) error {
	_, erreur := d.bd.Exec(ctx,
		`DELETE FROM cles WHERE projet = $1 AND utilisateur = $2`, projet, utilisateur)
	return erreur
}

func (d *Depot) CleParValeur(ctx context.Context, valeur string) (*modeles.Cle, *modeles.Utilisateur, error) {
	var cle modeles.Cle
	var utilisateur modeles.Utilisateur
	erreur := d.bd.QueryRow(ctx, `
		SELECT c.projet, c.utilisateur, c.cle, c.creation, u.id, u.courriel, u.nom, u.prenom
		FROM cles c JOIN utilisateurs u ON u.id = c.utilisateur
		WHERE c.cle = $1`, valeur).
		Scan(&cle.Projet, &cle.Utilisateur, &cle.Cle, &cle.Creation,
			&utilisateur.ID, &utilisateur.Courriel, &utilisateur.Nom, &utilisateur.Prenom)
	if erreur != nil {
		return nil, nil, erreur
	}
	return &cle, &utilisateur, nil
}
