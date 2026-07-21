package depots

import (
	"context"
	"sort"
	"time"

	"historykanban/serveur/internal/modeles"
)

const cteTerminees = `
	WITH terminees AS (
		SELECT t.id, t.points, t.modification, t.createur
		FROM taches t
		JOIN colonnes c ON c.id = t.colonne
		WHERE t.projet = ANY($1)
			AND t.modification BETWEEN $2 AND $3
			AND c.position = (SELECT max(cc.position) FROM colonnes cc WHERE cc.projet = t.projet)
	)`

func (d *Depot) IdentifiantsProjets(ctx context.Context, groupe string) ([]string, error) {
	lignes, erreur := d.bd.Query(ctx, `SELECT id FROM projets WHERE groupe = $1`, groupe)
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

func (d *Depot) Statistiques(ctx context.Context, projets []string, debut, fin time.Time, granularite string) (*modeles.Statistiques, error) {
	statistiques := &modeles.Statistiques{
		Classement: []modeles.LigneClassement{},
		Serie:      []modeles.PointSerie{},
	}
	if len(projets) == 0 {
		return statistiques, nil
	}
	parLigne := map[string]*modeles.LigneClassement{}
	assurer := func(utilisateur, nom, prenom string) *modeles.LigneClassement {
		ligne, present := parLigne[utilisateur]
		if !present {
			ligne = &modeles.LigneClassement{Utilisateur: utilisateur, Nom: nom, Prenom: prenom}
			parLigne[utilisateur] = ligne
		}
		return ligne
	}

	lignes, erreur := d.bd.Query(ctx, cteTerminees+`,
		contributions AS (
			SELECT te.id, te.points, coalesce(a.utilisateur, te.createur) AS utilisateur
			FROM terminees te
			LEFT JOIN affectations a ON a.tache = te.id
		)
		SELECT co.utilisateur, u.nom, u.prenom, coalesce(sum(co.points), 0), count(DISTINCT co.id)
		FROM contributions co
		JOIN utilisateurs u ON u.id = co.utilisateur
		GROUP BY co.utilisateur, u.nom, u.prenom`,
		projets, debut, fin)
	if erreur != nil {
		return nil, erreur
	}
	for lignes.Next() {
		var utilisateur, nom, prenom string
		var points, taches int
		if erreur := lignes.Scan(&utilisateur, &nom, &prenom, &points, &taches); erreur != nil {
			lignes.Close()
			return nil, erreur
		}
		ligne := assurer(utilisateur, nom, prenom)
		ligne.Points = points
		ligne.Taches = taches
	}
	lignes.Close()
	if erreur := lignes.Err(); erreur != nil {
		return nil, erreur
	}

	lignes, erreur = d.bd.Query(ctx, `
		SELECT t.createur, u.nom, u.prenom, count(*)
		FROM taches t
		JOIN utilisateurs u ON u.id = t.createur
		WHERE t.projet = ANY($1) AND t.creation BETWEEN $2 AND $3
		GROUP BY t.createur, u.nom, u.prenom`,
		projets, debut, fin)
	if erreur != nil {
		return nil, erreur
	}
	for lignes.Next() {
		var utilisateur, nom, prenom string
		var creees int
		if erreur := lignes.Scan(&utilisateur, &nom, &prenom, &creees); erreur != nil {
			lignes.Close()
			return nil, erreur
		}
		assurer(utilisateur, nom, prenom).Creees = creees
	}
	lignes.Close()
	if erreur := lignes.Err(); erreur != nil {
		return nil, erreur
	}

	lignes, erreur = d.bd.Query(ctx, `
		SELECT c.auteur, u.nom, u.prenom, count(*)
		FROM commentaires c
		JOIN taches t ON t.id = c.tache
		JOIN utilisateurs u ON u.id = c.auteur
		WHERE t.projet = ANY($1) AND c.creation BETWEEN $2 AND $3
		GROUP BY c.auteur, u.nom, u.prenom`,
		projets, debut, fin)
	if erreur != nil {
		return nil, erreur
	}
	for lignes.Next() {
		var utilisateur, nom, prenom string
		var commentaires int
		if erreur := lignes.Scan(&utilisateur, &nom, &prenom, &commentaires); erreur != nil {
			lignes.Close()
			return nil, erreur
		}
		assurer(utilisateur, nom, prenom).Commentaires = commentaires
	}
	lignes.Close()
	if erreur := lignes.Err(); erreur != nil {
		return nil, erreur
	}

	for _, ligne := range parLigne {
		statistiques.Classement = append(statistiques.Classement, *ligne)
	}
	sort.Slice(statistiques.Classement, func(gauche, droite int) bool {
		premier, second := statistiques.Classement[gauche], statistiques.Classement[droite]
		if premier.Points != second.Points {
			return premier.Points > second.Points
		}
		if premier.Taches != second.Taches {
			return premier.Taches > second.Taches
		}
		return premier.Creees > second.Creees
	})

	lignes, erreur = d.bd.Query(ctx, cteTerminees+`
		SELECT date_trunc($4, modification) AS periode, coalesce(sum(points), 0), count(*)
		FROM terminees
		GROUP BY periode
		ORDER BY periode`,
		projets, debut, fin, granularite)
	if erreur != nil {
		return nil, erreur
	}
	for lignes.Next() {
		var point modeles.PointSerie
		if erreur := lignes.Scan(&point.Periode, &point.Points, &point.Taches); erreur != nil {
			lignes.Close()
			return nil, erreur
		}
		statistiques.Serie = append(statistiques.Serie, point)
		statistiques.Totaux.Points += point.Points
		statistiques.Totaux.Terminees += point.Taches
	}
	lignes.Close()
	if erreur := lignes.Err(); erreur != nil {
		return nil, erreur
	}

	erreur = d.bd.QueryRow(ctx, `
		SELECT count(*),
			coalesce(sum(t.points), 0),
			count(*) FILTER (WHERE t.echeance < now()
				AND c.position <> (SELECT max(cc.position) FROM colonnes cc WHERE cc.projet = t.projet))
		FROM taches t
		JOIN colonnes c ON c.id = t.colonne
		WHERE t.projet = ANY($1)`, projets).
		Scan(&statistiques.Totaux.Total, &statistiques.Totaux.TotalPoints, &statistiques.Totaux.EnRetard)
	if erreur != nil {
		return nil, erreur
	}
	return statistiques, nil
}
