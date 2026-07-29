package depots

import (
	"context"

	"github.com/jackc/pgx/v5"

	"historykanban/serveur/internal/modeles"
)

func (d *Depot) SousTaches(ctx context.Context, tache string) ([]modeles.SousTache, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT id, tache, libelle, faite, position, creation
		FROM soustaches WHERE tache = $1 ORDER BY position, creation`, tache)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	return lireSousTaches(lignes)
}

func (d *Depot) attacherSousTaches(ctx context.Context, projet string, taches []modeles.Tache) ([]modeles.Tache, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT s.id, s.tache, s.libelle, s.faite, s.position, s.creation
		FROM soustaches s JOIN taches t ON t.id = s.tache
		WHERE t.projet = $1 ORDER BY s.position, s.creation`, projet)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	toutes, erreur := lireSousTaches(lignes)
	if erreur != nil {
		return nil, erreur
	}
	parTache := map[string][]modeles.SousTache{}
	for _, sousTache := range toutes {
		parTache[sousTache.Tache] = append(parTache[sousTache.Tache], sousTache)
	}
	for indice := range taches {
		if liste, presentes := parTache[taches[indice].ID]; presentes {
			taches[indice].SousTaches = liste
		}
	}
	return taches, nil
}

func lireSousTaches(lignes pgx.Rows) ([]modeles.SousTache, error) {
	sousTaches := []modeles.SousTache{}
	for lignes.Next() {
		var sousTache modeles.SousTache
		if erreur := lignes.Scan(&sousTache.ID, &sousTache.Tache, &sousTache.Libelle,
			&sousTache.Faite, &sousTache.Position, &sousTache.Creation); erreur != nil {
			return nil, erreur
		}
		sousTaches = append(sousTaches, sousTache)
	}
	return sousTaches, lignes.Err()
}

func (d *Depot) CreerSousTache(ctx context.Context, tache, libelle string) (*modeles.SousTache, error) {
	var sousTache modeles.SousTache
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO soustaches (tache, libelle, position)
		VALUES ($1, $2, (SELECT coalesce(max(position), -1) + 1 FROM soustaches WHERE tache = $1))
		RETURNING id, tache, libelle, faite, position, creation`, tache, libelle).
		Scan(&sousTache.ID, &sousTache.Tache, &sousTache.Libelle, &sousTache.Faite,
			&sousTache.Position, &sousTache.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &sousTache, nil
}

func (d *Depot) ModifierSousTache(ctx context.Context, id, libelle string, faite bool) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE soustaches SET libelle = $2, faite = $3 WHERE id = $1`, id, libelle, faite)
	return erreur
}

func (d *Depot) SupprimerSousTache(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM soustaches WHERE id = $1`, id)
	return erreur
}

func (d *Depot) ProjetSousTache(ctx context.Context, sousTache string) (string, error) {
	var projet string
	erreur := d.bd.QueryRow(ctx, `
		SELECT t.projet FROM soustaches s JOIN taches t ON t.id = s.tache WHERE s.id = $1`, sousTache).
		Scan(&projet)
	return projet, erreur
}

func (d *Depot) SousTache(ctx context.Context, id string) (*modeles.SousTache, error) {
	var sousTache modeles.SousTache
	erreur := d.bd.QueryRow(ctx, `
		SELECT id, tache, libelle, faite, position, creation FROM soustaches WHERE id = $1`, id).
		Scan(&sousTache.ID, &sousTache.Tache, &sousTache.Libelle, &sousTache.Faite,
			&sousTache.Position, &sousTache.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &sousTache, nil
}
