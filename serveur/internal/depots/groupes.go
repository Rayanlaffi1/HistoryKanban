package depots

import (
	"context"

	"github.com/jackc/pgx/v5"

	"historykanban/serveur/internal/modeles"
)

func (d *Depot) GroupesParUtilisateur(ctx context.Context, utilisateur string) ([]modeles.Groupe, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT g.id, g.nom, g.description, g.proprietaire, g.creation, m.role,
			(SELECT count(*) FROM membres WHERE groupe = g.id),
			(SELECT count(*) FROM projets WHERE groupe = g.id)
		FROM groupes g
		JOIN membres m ON m.groupe = g.id AND m.utilisateur = $1
		ORDER BY g.creation`, utilisateur)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	groupes := []modeles.Groupe{}
	for lignes.Next() {
		var groupe modeles.Groupe
		if erreur := lignes.Scan(&groupe.ID, &groupe.Nom, &groupe.Description, &groupe.Proprietaire,
			&groupe.Creation, &groupe.Role, &groupe.NbMembres, &groupe.NbProjets); erreur != nil {
			return nil, erreur
		}
		groupes = append(groupes, groupe)
	}
	return groupes, lignes.Err()
}

func (d *Depot) Groupe(ctx context.Context, id string) (*modeles.Groupe, error) {
	var groupe modeles.Groupe
	erreur := d.bd.QueryRow(ctx,
		`SELECT id, nom, description, proprietaire, creation FROM groupes WHERE id = $1`, id).
		Scan(&groupe.ID, &groupe.Nom, &groupe.Description, &groupe.Proprietaire, &groupe.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &groupe, nil
}

func (d *Depot) CreerGroupe(ctx context.Context, nom, description, proprietaire string) (*modeles.Groupe, error) {
	transaction, erreur := d.bd.Begin(ctx)
	if erreur != nil {
		return nil, erreur
	}
	defer transaction.Rollback(ctx)
	var groupe modeles.Groupe
	erreur = transaction.QueryRow(ctx, `
		INSERT INTO groupes (nom, description, proprietaire)
		VALUES ($1, $2, $3)
		RETURNING id, nom, description, proprietaire, creation`,
		nom, description, proprietaire).
		Scan(&groupe.ID, &groupe.Nom, &groupe.Description, &groupe.Proprietaire, &groupe.Creation)
	if erreur != nil {
		return nil, erreur
	}
	_, erreur = transaction.Exec(ctx,
		`INSERT INTO membres (groupe, utilisateur, role) VALUES ($1, $2, 'proprietaire')`,
		groupe.ID, proprietaire)
	if erreur != nil {
		return nil, erreur
	}
	if erreur = transaction.Commit(ctx); erreur != nil {
		return nil, erreur
	}
	groupe.Role = "proprietaire"
	groupe.NbMembres = 1
	return &groupe, nil
}

func (d *Depot) ModifierGroupe(ctx context.Context, id, nom, description string) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE groupes SET nom = $2, description = $3 WHERE id = $1`, id, nom, description)
	return erreur
}

func (d *Depot) SupprimerGroupe(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM groupes WHERE id = $1`, id)
	return erreur
}

func (d *Depot) RoleGroupe(ctx context.Context, groupe, utilisateur string) (string, error) {
	var role string
	erreur := d.bd.QueryRow(ctx,
		`SELECT role FROM membres WHERE groupe = $1 AND utilisateur = $2`, groupe, utilisateur).Scan(&role)
	if erreur == pgx.ErrNoRows {
		return "", nil
	}
	return role, erreur
}

func (d *Depot) Membres(ctx context.Context, groupe string) ([]modeles.Membre, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT m.groupe, m.utilisateur, m.role, m.fonction, m.ajout, u.courriel, u.nom, u.prenom
		FROM membres m
		JOIN utilisateurs u ON u.id = m.utilisateur
		WHERE m.groupe = $1
		ORDER BY m.ajout`, groupe)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	membres := []modeles.Membre{}
	for lignes.Next() {
		var membre modeles.Membre
		if erreur := lignes.Scan(&membre.Groupe, &membre.Utilisateur, &membre.Role, &membre.Fonction, &membre.Ajout,
			&membre.Courriel, &membre.Nom, &membre.Prenom); erreur != nil {
			return nil, erreur
		}
		membres = append(membres, membre)
	}
	return membres, lignes.Err()
}

func (d *Depot) MembresIdentifiants(ctx context.Context, groupe string) ([]string, error) {
	lignes, erreur := d.bd.Query(ctx, `SELECT utilisateur FROM membres WHERE groupe = $1`, groupe)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	identifiants := []string{}
	for lignes.Next() {
		var identifiant string
		if erreur := lignes.Scan(&identifiant); erreur != nil {
			return nil, erreur
		}
		identifiants = append(identifiants, identifiant)
	}
	return identifiants, lignes.Err()
}

func (d *Depot) AjouterMembre(ctx context.Context, groupe, utilisateur, role string) error {
	_, erreur := d.bd.Exec(ctx, `
		INSERT INTO membres (groupe, utilisateur, role) VALUES ($1, $2, $3)
		ON CONFLICT (groupe, utilisateur) DO NOTHING`, groupe, utilisateur, role)
	return erreur
}

func (d *Depot) ModifierRole(ctx context.Context, groupe, utilisateur, role string) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE membres SET role = $3 WHERE groupe = $1 AND utilisateur = $2`, groupe, utilisateur, role)
	return erreur
}

func (d *Depot) ModifierFonction(ctx context.Context, groupe, utilisateur, fonction string) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE membres SET fonction = $3 WHERE groupe = $1 AND utilisateur = $2`, groupe, utilisateur, fonction)
	return erreur
}

func (d *Depot) RetirerMembre(ctx context.Context, groupe, utilisateur string) error {
	_, erreur := d.bd.Exec(ctx,
		`DELETE FROM membres WHERE groupe = $1 AND utilisateur = $2 AND role <> 'proprietaire'`, groupe, utilisateur)
	return erreur
}
