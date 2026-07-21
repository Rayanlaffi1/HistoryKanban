package depots

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"historykanban/serveur/internal/modeles"
)

func (d *Depot) CreerNotification(ctx context.Context, utilisateur, categorie string, contenu any) (*modeles.Notification, error) {
	octets, erreur := json.Marshal(contenu)
	if erreur != nil {
		return nil, erreur
	}
	var notification modeles.Notification
	erreur = d.bd.QueryRow(ctx, `
		INSERT INTO notifications (utilisateur, type, contenu) VALUES ($1, $2, $3)
		RETURNING id, utilisateur, type, contenu, lue, creation`,
		utilisateur, categorie, octets).
		Scan(&notification.ID, &notification.Utilisateur, &notification.Type,
			&notification.Contenu, &notification.Lue, &notification.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &notification, nil
}

func (d *Depot) Notifications(ctx context.Context, utilisateur string) ([]modeles.Notification, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT id, utilisateur, type, contenu, lue, creation
		FROM notifications WHERE utilisateur = $1
		ORDER BY creation DESC LIMIT 50`, utilisateur)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	notifications := []modeles.Notification{}
	for lignes.Next() {
		var notification modeles.Notification
		if erreur := lignes.Scan(&notification.ID, &notification.Utilisateur, &notification.Type,
			&notification.Contenu, &notification.Lue, &notification.Creation); erreur != nil {
			return nil, erreur
		}
		notifications = append(notifications, notification)
	}
	return notifications, lignes.Err()
}

func (d *Depot) MarquerLue(ctx context.Context, id, utilisateur string) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE notifications SET lue = true WHERE id = $1 AND utilisateur = $2`, id, utilisateur)
	return erreur
}

func (d *Depot) MarquerToutesLues(ctx context.Context, utilisateur string) error {
	_, erreur := d.bd.Exec(ctx, `UPDATE notifications SET lue = true WHERE utilisateur = $1`, utilisateur)
	return erreur
}

func (d *Depot) Preferences(ctx context.Context, utilisateur string) (*modeles.Preferences, error) {
	preferences := modeles.Preferences{Utilisateur: utilisateur, Courriels: true, Types: map[string]bool{}}
	var octets []byte
	erreur := d.bd.QueryRow(ctx,
		`SELECT courriels, types FROM preferences WHERE utilisateur = $1`, utilisateur).
		Scan(&preferences.Courriels, &octets)
	if erreur == pgx.ErrNoRows {
		return &preferences, nil
	}
	if erreur != nil {
		return nil, erreur
	}
	if erreur = json.Unmarshal(octets, &preferences.Types); erreur != nil {
		preferences.Types = map[string]bool{}
	}
	return &preferences, nil
}

func (d *Depot) EnregistrerPreferences(ctx context.Context, preferences modeles.Preferences) error {
	octets, erreur := json.Marshal(preferences.Types)
	if erreur != nil {
		return erreur
	}
	_, erreur = d.bd.Exec(ctx, `
		INSERT INTO preferences (utilisateur, courriels, types) VALUES ($1, $2, $3)
		ON CONFLICT (utilisateur) DO UPDATE SET courriels = EXCLUDED.courriels, types = EXCLUDED.types`,
		preferences.Utilisateur, preferences.Courriels, octets)
	return erreur
}

func (d *Depot) CourrielsUtilisateurs(ctx context.Context, identifiants []string, categorie string) ([]string, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT u.courriel FROM utilisateurs u
		LEFT JOIN preferences p ON p.utilisateur = u.id
		WHERE u.id = ANY($1)
			AND coalesce(p.courriels, true)
			AND coalesce((p.types ->> $2)::boolean, true)`,
		identifiants, categorie)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	courriels := []string{}
	for lignes.Next() {
		var courriel string
		if erreur := lignes.Scan(&courriel); erreur != nil {
			return nil, erreur
		}
		courriels = append(courriels, courriel)
	}
	return courriels, lignes.Err()
}
