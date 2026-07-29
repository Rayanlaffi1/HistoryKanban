package depots

import (
	"context"
	"sort"
	"time"

	"historykanban/serveur/internal/modeles"
)

const cteTerminees = `
	WITH terminees AS (
		SELECT t.id, t.points, t.terminee, t.createur
		FROM taches t
		JOIN colonnes c ON c.id = t.colonne
		WHERE t.projet = ANY($1) AND t.suppression IS NULL
			AND t.terminee BETWEEN $2 AND $3
			AND c.position = (SELECT max(cc.position) FROM colonnes cc WHERE cc.projet = t.projet)
	)`

// Les agregats temporels sont regroupes sur la journee civile locale via le littéral
// SQL « AT TIME ZONE 'Europe/Paris' » : une tache terminee le 28/07 a 23h30 UTC compte
// pour le 29/07 a Paris. Postgres embarque la base des fuseaux, donc le passage heure
// d'ete / heure d'hiver est pris en charge automatiquement.

// Plafond de la liste detaillee des taches terminees pour eviter une reponse non bornee.
const (
	limiteTermineesParDefaut = 100
	limiteTermineesMax       = 500
)

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

func (d *Depot) Statistiques(ctx context.Context, projets []string, debut, fin time.Time, granularite string, limite, offset int) (*modeles.Statistiques, error) {
	statistiques := &modeles.Statistiques{
		Classement: []modeles.LigneClassement{},
		Serie:      []modeles.PointSerie{},
		Terminees:  []modeles.PeriodeTermineesStatistiques{},
	}
	if len(projets) == 0 {
		return statistiques, nil
	}
	if limite <= 0 {
		limite = limiteTermineesParDefaut
	}
	if limite > limiteTermineesMax {
		limite = limiteTermineesMax
	}
	if offset < 0 {
		offset = 0
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

	// Les points d'une tache sont repartis a parts egales entre ses affectes (a defaut,
	// credites a sa creatrice ou son createur). count(*) OVER vaut 1 quand la tache n'a aucune
	// affectation (jointure a gauche => une seule ligne). On evite ainsi le double comptage :
	// la somme des points du classement egale la somme des points des taches terminees.
	lignes, erreur := d.bd.Query(ctx, cteTerminees+`,
		contributions AS (
			SELECT te.id,
				coalesce(a.utilisateur, te.createur) AS utilisateur,
				te.points::double precision / count(*) OVER (PARTITION BY te.id) AS points
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
		var points float64
		var taches int
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
		WHERE t.projet = ANY($1) AND t.suppression IS NULL AND t.creation BETWEEN $2 AND $3
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
		WHERE t.projet = ANY($1) AND t.suppression IS NULL AND c.creation BETWEEN $2 AND $3
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
		SELECT date_trunc($4, terminee AT TIME ZONE 'Europe/Paris') AS periode, coalesce(sum(points), 0), count(*)
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

	lignes, erreur = d.bd.Query(ctx, cteTerminees+`
		SELECT date_trunc($4, te.terminee AT TIME ZONE 'Europe/Paris') AS periode,
			te.id, p.id, p.nom, p.couleur, t.titre, te.points, t.urgence, te.terminee
		FROM terminees te
		JOIN taches t ON t.id = te.id
		JOIN projets p ON p.id = t.projet
		ORDER BY periode DESC, p.nom, te.terminee DESC, t.titre
		LIMIT $5 OFFSET $6`,
		projets, debut, fin, granularite, limite, offset)
	if erreur != nil {
		return nil, erreur
	}
	parPeriode := map[time.Time]int{}
	for lignes.Next() {
		var periode time.Time
		var tache modeles.TacheTermineeStatistiques
		if erreur := lignes.Scan(&periode, &tache.ID, &tache.Projet, &tache.ProjetNom, &tache.ProjetCouleur,
			&tache.Titre, &tache.Points, &tache.Urgence, &tache.Terminee); erreur != nil {
			lignes.Close()
			return nil, erreur
		}
		indice, present := parPeriode[periode]
		if !present {
			statistiques.Terminees = append(statistiques.Terminees, modeles.PeriodeTermineesStatistiques{
				Periode: periode,
				Taches:  []modeles.TacheTermineeStatistiques{},
			})
			indice = len(statistiques.Terminees) - 1
			parPeriode[periode] = indice
		}
		periodeTerminee := &statistiques.Terminees[indice]
		periodeTerminee.Taches = append(periodeTerminee.Taches, tache)
		periodeTerminee.Total++
		periodeTerminee.Points += tache.Points
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
		WHERE t.projet = ANY($1) AND t.suppression IS NULL`, projets).
		Scan(&statistiques.Totaux.Total, &statistiques.Totaux.TotalPoints, &statistiques.Totaux.EnRetard)
	if erreur != nil {
		return nil, erreur
	}
	return statistiques, nil
}
