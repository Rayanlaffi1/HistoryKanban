package depots

import (
	"context"

	"historykanban/serveur/internal/modeles"
)

func (d *Depot) ProjetsParGroupe(ctx context.Context, groupe string) ([]modeles.Projet, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT p.id, p.groupe, p.nom, p.description, p.couleur, p.archive, p.createur, p.creation,
			(SELECT count(*) FROM taches WHERE projet = p.id)
		FROM projets p WHERE p.groupe = $1 ORDER BY p.creation`, groupe)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	projets := []modeles.Projet{}
	for lignes.Next() {
		var projet modeles.Projet
		if erreur := lignes.Scan(&projet.ID, &projet.Groupe, &projet.Nom, &projet.Description,
			&projet.Couleur, &projet.Archive, &projet.Createur, &projet.Creation, &projet.NbTaches); erreur != nil {
			return nil, erreur
		}
		projets = append(projets, projet)
	}
	return projets, lignes.Err()
}

func (d *Depot) Projet(ctx context.Context, id string) (*modeles.Projet, error) {
	var projet modeles.Projet
	erreur := d.bd.QueryRow(ctx, `
		SELECT id, groupe, nom, description, couleur, archive, createur, creation
		FROM projets WHERE id = $1`, id).
		Scan(&projet.ID, &projet.Groupe, &projet.Nom, &projet.Description,
			&projet.Couleur, &projet.Archive, &projet.Createur, &projet.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &projet, nil
}

func (d *Depot) CreerProjet(ctx context.Context, projet modeles.Projet) (*modeles.Projet, error) {
	transaction, erreur := d.bd.Begin(ctx)
	if erreur != nil {
		return nil, erreur
	}
	defer transaction.Rollback(ctx)
	erreur = transaction.QueryRow(ctx, `
		INSERT INTO projets (groupe, nom, description, couleur, createur)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, groupe, nom, description, couleur, archive, createur, creation`,
		projet.Groupe, projet.Nom, projet.Description, projet.Couleur, projet.Createur).
		Scan(&projet.ID, &projet.Groupe, &projet.Nom, &projet.Description,
			&projet.Couleur, &projet.Archive, &projet.Createur, &projet.Creation)
	if erreur != nil {
		return nil, erreur
	}
	colonnes := []struct {
		nom     string
		couleur string
	}{
		{"À faire", "#737373"},
		{"En cours", "#525252"},
		{"Terminé", "#404040"},
	}
	for position, colonne := range colonnes {
		_, erreur = transaction.Exec(ctx,
			`INSERT INTO colonnes (projet, nom, couleur, position) VALUES ($1, $2, $3, $4)`,
			projet.ID, colonne.nom, colonne.couleur, position)
		if erreur != nil {
			return nil, erreur
		}
	}
	if erreur = transaction.Commit(ctx); erreur != nil {
		return nil, erreur
	}
	return &projet, nil
}

func (d *Depot) ModifierProjet(ctx context.Context, id, nom, description, couleur string, archive bool) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE projets SET nom = $2, description = $3, couleur = $4, archive = $5 WHERE id = $1`,
		id, nom, description, couleur, archive)
	return erreur
}

func (d *Depot) SupprimerProjet(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM projets WHERE id = $1`, id)
	return erreur
}

func (d *Depot) Colonnes(ctx context.Context, projet string) ([]modeles.Colonne, error) {
	lignes, erreur := d.bd.Query(ctx,
		`SELECT id, projet, nom, couleur, position, limite FROM colonnes WHERE projet = $1 ORDER BY position`, projet)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	colonnes := []modeles.Colonne{}
	for lignes.Next() {
		var colonne modeles.Colonne
		if erreur := lignes.Scan(&colonne.ID, &colonne.Projet, &colonne.Nom, &colonne.Couleur,
			&colonne.Position, &colonne.Limite); erreur != nil {
			return nil, erreur
		}
		colonnes = append(colonnes, colonne)
	}
	return colonnes, lignes.Err()
}

func (d *Depot) CreerColonne(ctx context.Context, colonne modeles.Colonne) (*modeles.Colonne, error) {
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO colonnes (projet, nom, couleur, limite, position)
		VALUES ($1, $2, $3, $4, (SELECT coalesce(max(position), -1) + 1 FROM colonnes WHERE projet = $1))
		RETURNING id, projet, nom, couleur, position, limite`,
		colonne.Projet, colonne.Nom, colonne.Couleur, colonne.Limite).
		Scan(&colonne.ID, &colonne.Projet, &colonne.Nom, &colonne.Couleur, &colonne.Position, &colonne.Limite)
	if erreur != nil {
		return nil, erreur
	}
	return &colonne, nil
}

func (d *Depot) ModifierColonne(ctx context.Context, id, nom, couleur string, limite *int) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE colonnes SET nom = $2, couleur = $3, limite = $4 WHERE id = $1`, id, nom, couleur, limite)
	return erreur
}

func (d *Depot) SupprimerColonne(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM colonnes WHERE id = $1`, id)
	return erreur
}

func (d *Depot) ReordonnerColonnes(ctx context.Context, projet string, ordre []string) error {
	transaction, erreur := d.bd.Begin(ctx)
	if erreur != nil {
		return erreur
	}
	defer transaction.Rollback(ctx)
	for position, id := range ordre {
		if _, erreur = transaction.Exec(ctx,
			`UPDATE colonnes SET position = $3 WHERE id = $1 AND projet = $2`, id, projet, position); erreur != nil {
			return erreur
		}
	}
	return transaction.Commit(ctx)
}

func (d *Depot) ProjetColonne(ctx context.Context, colonne string) (string, error) {
	var projet string
	erreur := d.bd.QueryRow(ctx, `SELECT projet FROM colonnes WHERE id = $1`, colonne).Scan(&projet)
	return projet, erreur
}

func (d *Depot) ProjetEtiquette(ctx context.Context, etiquette string) (string, error) {
	var projet string
	erreur := d.bd.QueryRow(ctx, `SELECT projet FROM etiquettes WHERE id = $1`, etiquette).Scan(&projet)
	return projet, erreur
}

func (d *Depot) ProjetLot(ctx context.Context, lot string) (string, error) {
	var projet string
	erreur := d.bd.QueryRow(ctx, `SELECT projet FROM lots WHERE id = $1`, lot).Scan(&projet)
	return projet, erreur
}

func (d *Depot) Etiquettes(ctx context.Context, projet string) ([]modeles.Etiquette, error) {
	lignes, erreur := d.bd.Query(ctx,
		`SELECT id, projet, nom, couleur FROM etiquettes WHERE projet = $1 ORDER BY nom`, projet)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	etiquettes := []modeles.Etiquette{}
	for lignes.Next() {
		var etiquette modeles.Etiquette
		if erreur := lignes.Scan(&etiquette.ID, &etiquette.Projet, &etiquette.Nom, &etiquette.Couleur); erreur != nil {
			return nil, erreur
		}
		etiquettes = append(etiquettes, etiquette)
	}
	return etiquettes, lignes.Err()
}

func (d *Depot) CreerEtiquette(ctx context.Context, etiquette modeles.Etiquette) (*modeles.Etiquette, error) {
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO etiquettes (projet, nom, couleur) VALUES ($1, $2, $3)
		RETURNING id, projet, nom, couleur`,
		etiquette.Projet, etiquette.Nom, etiquette.Couleur).
		Scan(&etiquette.ID, &etiquette.Projet, &etiquette.Nom, &etiquette.Couleur)
	if erreur != nil {
		return nil, erreur
	}
	return &etiquette, nil
}

func (d *Depot) ModifierEtiquette(ctx context.Context, id, nom, couleur string) error {
	_, erreur := d.bd.Exec(ctx, `UPDATE etiquettes SET nom = $2, couleur = $3 WHERE id = $1`, id, nom, couleur)
	return erreur
}

func (d *Depot) SupprimerEtiquette(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM etiquettes WHERE id = $1`, id)
	return erreur
}

func (d *Depot) Lots(ctx context.Context, projet string) ([]modeles.Lot, error) {
	lignes, erreur := d.bd.Query(ctx,
		`SELECT id, projet, nom, couleur, echeance, creation FROM lots WHERE projet = $1 ORDER BY creation`, projet)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	lots := []modeles.Lot{}
	for lignes.Next() {
		var lot modeles.Lot
		if erreur := lignes.Scan(&lot.ID, &lot.Projet, &lot.Nom, &lot.Couleur, &lot.Echeance, &lot.Creation); erreur != nil {
			return nil, erreur
		}
		lots = append(lots, lot)
	}
	return lots, lignes.Err()
}

func (d *Depot) CreerLot(ctx context.Context, lot modeles.Lot) (*modeles.Lot, error) {
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO lots (projet, nom, couleur, echeance) VALUES ($1, $2, $3, $4)
		RETURNING id, projet, nom, couleur, echeance, creation`,
		lot.Projet, lot.Nom, lot.Couleur, lot.Echeance).
		Scan(&lot.ID, &lot.Projet, &lot.Nom, &lot.Couleur, &lot.Echeance, &lot.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &lot, nil
}

func (d *Depot) ModifierLot(ctx context.Context, lot modeles.Lot) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE lots SET nom = $2, couleur = $3, echeance = $4 WHERE id = $1`,
		lot.ID, lot.Nom, lot.Couleur, lot.Echeance)
	return erreur
}

func (d *Depot) SupprimerLot(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM lots WHERE id = $1`, id)
	return erreur
}
