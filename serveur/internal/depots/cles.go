package depots

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"regexp"

	"github.com/jackc/pgx/v5"

	"historykanban/serveur/internal/modeles"
)

var expressionEmpreinteCle = regexp.MustCompile(`^[a-f0-9]{64}$`)

func empreinteCle(valeur string) string {
	somme := sha256.Sum256([]byte(valeur))
	return hex.EncodeToString(somme[:])
}

func estEmpreinteCle(valeur string) bool {
	return expressionEmpreinteCle.MatchString(valeur)
}

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
	cle.Cle = ""
	return &cle, nil
}

func (d *Depot) EnregistrerCle(ctx context.Context, projet, utilisateur, valeur string) (*modeles.Cle, error) {
	var cle modeles.Cle
	emp := empreinteCle(valeur)
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO cles (projet, utilisateur, cle) VALUES ($1, $2, $3)
		ON CONFLICT (projet, utilisateur) DO UPDATE SET cle = $3, creation = now()
		RETURNING projet, utilisateur, cle, creation`,
		projet, utilisateur, emp).
		Scan(&cle.Projet, &cle.Utilisateur, &cle.Cle, &cle.Creation)
	if erreur != nil {
		return nil, erreur
	}
	cle.Cle = valeur
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
	emp := empreinteCle(valeur)
	erreur := d.bd.QueryRow(ctx, `
		SELECT c.projet, c.utilisateur, c.cle, c.creation, u.id, u.courriel, u.nom, u.prenom
		FROM cles c JOIN utilisateurs u ON u.id = c.utilisateur
		WHERE c.cle = $1`, emp).
		Scan(&cle.Projet, &cle.Utilisateur, &cle.Cle, &cle.Creation,
			&utilisateur.ID, &utilisateur.Courriel, &utilisateur.Nom, &utilisateur.Prenom)
	if erreur == pgx.ErrNoRows {
		// Compatibilite transitoire pour les installations non migrees : accepter une
		// ancienne cle en clair une seule fois, puis la remplacer par son empreinte.
		erreur = d.bd.QueryRow(ctx, `
			UPDATE cles c SET cle = $2, creation = now()
			FROM utilisateurs u
			WHERE u.id = c.utilisateur AND c.cle = $1
			RETURNING c.projet, c.utilisateur, c.cle, c.creation, u.id, u.courriel, u.nom, u.prenom`, valeur, emp).
			Scan(&cle.Projet, &cle.Utilisateur, &cle.Cle, &cle.Creation,
				&utilisateur.ID, &utilisateur.Courriel, &utilisateur.Nom, &utilisateur.Prenom)
	}
	if erreur != nil {
		return nil, nil, erreur
	}
	cle.Cle = ""
	return &cle, &utilisateur, nil
}
