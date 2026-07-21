package depots

import (
	"context"
	"fmt"
	"strings"

	"historykanban/serveur/internal/modeles"
)

const selectionTaches = `
	SELECT t.id, t.projet, t.colonne, t.lot, t.titre, t.description, t.points, t.urgence, t.echeance, t.commit,
		t.position, t.suppression, t.createur, t.creation, t.modification,
		coalesce(array_agg(DISTINCT a.utilisateur::text) FILTER (WHERE a.utilisateur IS NOT NULL), '{}'),
		coalesce(array_agg(DISTINCT e.etiquette::text) FILTER (WHERE e.etiquette IS NOT NULL), '{}')
	FROM taches t
	LEFT JOIN affectations a ON a.tache = t.id
	LEFT JOIN etiquetages e ON e.tache = t.id`

func (d *Depot) Taches(ctx context.Context, projet string, filtre modeles.Filtre) ([]modeles.Tache, error) {
	conditions := []string{"t.projet = $1", "t.suppression IS NULL"}
	parametres := []any{projet}
	suivant := 2
	ajouter := func(condition string, valeur any) {
		conditions = append(conditions, fmt.Sprintf(condition, suivant))
		parametres = append(parametres, valeur)
		suivant++
	}
	if filtre.Texte != "" {
		ajouter("(t.titre ILIKE '%%' || $%d || '%%' OR t.description ILIKE '%%' || $%[1]d || '%%')", filtre.Texte)
	}
	if filtre.Membre != "" {
		ajouter("EXISTS (SELECT 1 FROM affectations af WHERE af.tache = t.id AND af.utilisateur = $%d)", filtre.Membre)
	}
	if filtre.Etiquette != "" {
		ajouter("EXISTS (SELECT 1 FROM etiquetages et WHERE et.tache = t.id AND et.etiquette = $%d)", filtre.Etiquette)
	}
	if filtre.Lot != "" {
		ajouter("t.lot = $%d", filtre.Lot)
	}
	if filtre.Urgence != "" {
		ajouter("t.urgence = $%d", filtre.Urgence)
	}
	if filtre.PointsMin != nil {
		ajouter("t.points >= $%d", *filtre.PointsMin)
	}
	if filtre.PointsMax != nil {
		ajouter("t.points <= $%d", *filtre.PointsMax)
	}
	switch filtre.Echeance {
	case "depassee":
		conditions = append(conditions, "t.echeance < now()")
	case "semaine":
		conditions = append(conditions, "t.echeance BETWEEN now() AND now() + interval '7 days'")
	case "sans":
		conditions = append(conditions, "t.echeance IS NULL")
	}
	requete := selectionTaches + "\n\tWHERE " + strings.Join(conditions, " AND ") +
		"\n\tGROUP BY t.id\n\tORDER BY t.position, t.creation"
	lignes, erreur := d.bd.Query(ctx, requete, parametres...)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	taches := []modeles.Tache{}
	for lignes.Next() {
		var tache modeles.Tache
		if erreur := lignes.Scan(&tache.ID, &tache.Projet, &tache.Colonne, &tache.Lot, &tache.Titre,
			&tache.Description, &tache.Points, &tache.Urgence, &tache.Echeance, &tache.Commit, &tache.Position,
			&tache.Suppression, &tache.Createur, &tache.Creation, &tache.Modification,
			&tache.Affectations, &tache.Etiquettes); erreur != nil {
			return nil, erreur
		}
		tache.Images = []modeles.Image{}
		tache.SousTaches = []modeles.SousTache{}
		taches = append(taches, tache)
	}
	if erreur := lignes.Err(); erreur != nil {
		return nil, erreur
	}
	taches, erreur = d.attacherImages(ctx, projet, taches)
	if erreur != nil {
		return nil, erreur
	}
	return d.attacherSousTaches(ctx, projet, taches)
}

func (d *Depot) attacherImages(ctx context.Context, projet string, taches []modeles.Tache) ([]modeles.Tache, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT i.id, i.tache, i.chemin, i.nom, i.taille, i.typecontenu, i.creation
		FROM images i JOIN taches t ON t.id = i.tache
		WHERE t.projet = $1 ORDER BY i.creation`, projet)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	parTache := map[string][]modeles.Image{}
	for lignes.Next() {
		var image modeles.Image
		if erreur := lignes.Scan(&image.ID, &image.Tache, &image.Chemin, &image.Nom, &image.Taille,
			&image.TypeContenu, &image.Creation); erreur != nil {
			return nil, erreur
		}
		parTache[image.Tache] = append(parTache[image.Tache], image)
	}
	if erreur := lignes.Err(); erreur != nil {
		return nil, erreur
	}
	for indice := range taches {
		if images, presentes := parTache[taches[indice].ID]; presentes {
			taches[indice].Images = images
		}
	}
	return taches, nil
}

func (d *Depot) Tache(ctx context.Context, id string) (*modeles.Tache, error) {
	var tache modeles.Tache
	erreur := d.bd.QueryRow(ctx, selectionTaches+"\n\tWHERE t.id = $1\n\tGROUP BY t.id", id).
		Scan(&tache.ID, &tache.Projet, &tache.Colonne, &tache.Lot, &tache.Titre,
			&tache.Description, &tache.Points, &tache.Urgence, &tache.Echeance, &tache.Commit, &tache.Position,
			&tache.Suppression, &tache.Createur, &tache.Creation, &tache.Modification,
			&tache.Affectations, &tache.Etiquettes)
	if erreur != nil {
		return nil, erreur
	}
	tache.Images = []modeles.Image{}
	lignes, erreur := d.bd.Query(ctx,
		`SELECT id, tache, chemin, nom, taille, typecontenu, creation FROM images WHERE tache = $1 ORDER BY creation`, id)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	for lignes.Next() {
		var image modeles.Image
		if erreur := lignes.Scan(&image.ID, &image.Tache, &image.Chemin, &image.Nom, &image.Taille,
			&image.TypeContenu, &image.Creation); erreur != nil {
			return nil, erreur
		}
		tache.Images = append(tache.Images, image)
	}
	if erreur := lignes.Err(); erreur != nil {
		return nil, erreur
	}
	tache.SousTaches, erreur = d.SousTaches(ctx, id)
	if erreur != nil {
		return nil, erreur
	}
	return &tache, nil
}

func (d *Depot) CreerTache(ctx context.Context, tache modeles.Tache) (*modeles.Tache, error) {
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO taches (projet, colonne, lot, titre, description, points, urgence, echeance, createur, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,
			(SELECT coalesce(max(position), -1) + 1 FROM taches WHERE colonne = $2))
		RETURNING id, position, creation, modification`,
		tache.Projet, tache.Colonne, tache.Lot, tache.Titre, tache.Description,
		tache.Points, tache.Urgence, tache.Echeance, tache.Createur).
		Scan(&tache.ID, &tache.Position, &tache.Creation, &tache.Modification)
	if erreur != nil {
		return nil, erreur
	}
	if len(tache.Affectations) > 0 {
		if erreur = d.Affecter(ctx, tache.ID, tache.Affectations); erreur != nil {
			d.SupprimerTache(ctx, tache.ID)
			return nil, erreur
		}
	}
	if len(tache.Etiquettes) > 0 {
		if erreur = d.Etiqueter(ctx, tache.ID, tache.Etiquettes); erreur != nil {
			d.SupprimerTache(ctx, tache.ID)
			return nil, erreur
		}
	}
	return d.Tache(ctx, tache.ID)
}

func (d *Depot) ModifierTache(ctx context.Context, tache modeles.Tache) error {
	_, erreur := d.bd.Exec(ctx, `
		UPDATE taches SET titre = $2, description = $3, points = $4, echeance = $5, lot = $6, commit = $7,
			urgence = $8, modification = now()
		WHERE id = $1`,
		tache.ID, tache.Titre, tache.Description, tache.Points, tache.Echeance, tache.Lot, tache.Commit, tache.Urgence)
	return erreur
}

func (d *Depot) ModifierCommit(ctx context.Context, id, commit string) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE taches SET commit = $2, modification = now() WHERE id = $1`, id, commit)
	return erreur
}

func (d *Depot) SupprimerTache(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx,
		`UPDATE taches SET suppression = now(), modification = now() WHERE id = $1 AND suppression IS NULL`, id)
	return erreur
}

func (d *Depot) RestaurerTache(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `
		UPDATE taches SET suppression = NULL, modification = now(),
			position = (SELECT coalesce(max(position), -1) + 1 FROM taches p
				WHERE p.colonne = taches.colonne AND p.suppression IS NULL)
		WHERE id = $1 AND suppression IS NOT NULL`, id)
	return erreur
}

func (d *Depot) PurgerTache(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM taches WHERE id = $1`, id)
	return erreur
}

func (d *Depot) Corbeille(ctx context.Context, projet string) ([]modeles.Tache, error) {
	requete := selectionTaches + `
	WHERE t.projet = $1 AND t.suppression IS NOT NULL
	GROUP BY t.id
	ORDER BY t.suppression DESC`
	lignes, erreur := d.bd.Query(ctx, requete, projet)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	taches := []modeles.Tache{}
	for lignes.Next() {
		var tache modeles.Tache
		if erreur := lignes.Scan(&tache.ID, &tache.Projet, &tache.Colonne, &tache.Lot, &tache.Titre,
			&tache.Description, &tache.Points, &tache.Urgence, &tache.Echeance, &tache.Commit, &tache.Position,
			&tache.Suppression, &tache.Createur, &tache.Creation, &tache.Modification,
			&tache.Affectations, &tache.Etiquettes); erreur != nil {
			return nil, erreur
		}
		tache.Images = []modeles.Image{}
		taches = append(taches, tache)
	}
	return taches, lignes.Err()
}

func (d *Depot) TachesExpirees(ctx context.Context, jours int) ([]string, []string, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT t.id, coalesce(array_agg(i.chemin) FILTER (WHERE i.chemin IS NOT NULL), '{}')
		FROM taches t LEFT JOIN images i ON i.tache = t.id
		WHERE t.suppression IS NOT NULL AND t.suppression < now() - make_interval(days => $1)
		GROUP BY t.id`, jours)
	if erreur != nil {
		return nil, nil, erreur
	}
	defer lignes.Close()
	identifiants := []string{}
	chemins := []string{}
	for lignes.Next() {
		var identifiant string
		var fichiers []string
		if erreur := lignes.Scan(&identifiant, &fichiers); erreur != nil {
			return nil, nil, erreur
		}
		identifiants = append(identifiants, identifiant)
		chemins = append(chemins, fichiers...)
	}
	return identifiants, chemins, lignes.Err()
}

func (d *Depot) DeplacerTache(ctx context.Context, id, colonne string, position int) error {
	transaction, erreur := d.bd.Begin(ctx)
	if erreur != nil {
		return erreur
	}
	defer transaction.Rollback(ctx)
	if _, erreur = transaction.Exec(ctx,
		`UPDATE taches SET position = position + 1 WHERE colonne = $1 AND position >= $2`,
		colonne, position); erreur != nil {
		return erreur
	}
	if _, erreur = transaction.Exec(ctx,
		`UPDATE taches SET colonne = $2, position = $3, modification = now() WHERE id = $1`,
		id, colonne, position); erreur != nil {
		return erreur
	}
	return transaction.Commit(ctx)
}

func distincts(valeurs []string) []string {
	vus := map[string]bool{}
	resultat := []string{}
	for _, valeur := range valeurs {
		if valeur == "" || vus[valeur] {
			continue
		}
		vus[valeur] = true
		resultat = append(resultat, valeur)
	}
	return resultat
}

func (d *Depot) Affecter(ctx context.Context, tache string, utilisateurs []string) error {
	attendus := distincts(utilisateurs)
	if len(attendus) > 0 {
		var valides int
		erreur := d.bd.QueryRow(ctx, `
			SELECT count(DISTINCT m.utilisateur)
			FROM taches t
			JOIN projets p ON p.id = t.projet
			JOIN membres m ON m.groupe = p.groupe
			WHERE t.id = $1 AND m.utilisateur::text = ANY($2)`, tache, attendus).Scan(&valides)
		if erreur != nil {
			return erreur
		}
		if valides != len(attendus) {
			return fmt.Errorf("un utilisateur affecte n'appartient pas au groupe du projet")
		}
	}
	transaction, erreur := d.bd.Begin(ctx)
	if erreur != nil {
		return erreur
	}
	defer transaction.Rollback(ctx)
	if _, erreur = transaction.Exec(ctx, `DELETE FROM affectations WHERE tache = $1`, tache); erreur != nil {
		return erreur
	}
	for _, utilisateur := range utilisateurs {
		if _, erreur = transaction.Exec(ctx,
			`INSERT INTO affectations (tache, utilisateur) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			tache, utilisateur); erreur != nil {
			return erreur
		}
	}
	return transaction.Commit(ctx)
}

func (d *Depot) Etiqueter(ctx context.Context, tache string, etiquettes []string) error {
	attendues := distincts(etiquettes)
	if len(attendues) > 0 {
		var valides int
		erreur := d.bd.QueryRow(ctx, `
			SELECT count(DISTINCT e.id)
			FROM taches t
			JOIN etiquettes e ON e.projet = t.projet
			WHERE t.id = $1 AND e.id::text = ANY($2)`, tache, attendues).Scan(&valides)
		if erreur != nil {
			return erreur
		}
		if valides != len(attendues) {
			return fmt.Errorf("une etiquette n'appartient pas a ce projet")
		}
	}
	transaction, erreur := d.bd.Begin(ctx)
	if erreur != nil {
		return erreur
	}
	defer transaction.Rollback(ctx)
	if _, erreur = transaction.Exec(ctx, `DELETE FROM etiquetages WHERE tache = $1`, tache); erreur != nil {
		return erreur
	}
	for _, etiquette := range etiquettes {
		if _, erreur = transaction.Exec(ctx,
			`INSERT INTO etiquetages (tache, etiquette) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			tache, etiquette); erreur != nil {
			return erreur
		}
	}
	return transaction.Commit(ctx)
}

func (d *Depot) ProjetTache(ctx context.Context, tache string) (string, error) {
	var projet string
	erreur := d.bd.QueryRow(ctx, `SELECT projet FROM taches WHERE id = $1`, tache).Scan(&projet)
	return projet, erreur
}

func (d *Depot) AjouterImage(ctx context.Context, image modeles.Image) (*modeles.Image, error) {
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO images (tache, chemin, nom, taille, typecontenu) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, tache, chemin, nom, taille, typecontenu, creation`,
		image.Tache, image.Chemin, image.Nom, image.Taille, image.TypeContenu).
		Scan(&image.ID, &image.Tache, &image.Chemin, &image.Nom, &image.Taille, &image.TypeContenu, &image.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &image, nil
}

func (d *Depot) Image(ctx context.Context, id string) (*modeles.Image, error) {
	var image modeles.Image
	erreur := d.bd.QueryRow(ctx,
		`SELECT id, tache, chemin, nom, taille, typecontenu, creation FROM images WHERE id = $1`, id).
		Scan(&image.ID, &image.Tache, &image.Chemin, &image.Nom, &image.Taille, &image.TypeContenu, &image.Creation)
	if erreur != nil {
		return nil, erreur
	}
	return &image, nil
}

func (d *Depot) SupprimerImage(ctx context.Context, id string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM images WHERE id = $1`, id)
	return erreur
}

func (d *Depot) NomColonne(ctx context.Context, colonne string) (string, error) {
	var nom string
	erreur := d.bd.QueryRow(ctx, `SELECT nom FROM colonnes WHERE id = $1`, colonne).Scan(&nom)
	return nom, erreur
}

func (d *Depot) CreerActivite(ctx context.Context, tache, utilisateur, categorie, detail string) error {
	_, erreur := d.bd.Exec(ctx,
		`INSERT INTO activites (tache, utilisateur, type, detail) VALUES ($1, $2, $3, $4)`,
		tache, utilisateur, categorie, detail)
	return erreur
}

func (d *Depot) Activites(ctx context.Context, tache string) ([]modeles.Activite, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT a.id, a.tache, coalesce(a.utilisateur::text, ''), a.type, a.detail, a.creation,
			coalesce(u.nom, ''), coalesce(u.prenom, '')
		FROM activites a
		LEFT JOIN utilisateurs u ON u.id = a.utilisateur
		WHERE a.tache = $1
		ORDER BY a.creation DESC
		LIMIT 100`, tache)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	activites := []modeles.Activite{}
	for lignes.Next() {
		var activite modeles.Activite
		if erreur := lignes.Scan(&activite.ID, &activite.Tache, &activite.Utilisateur, &activite.Type,
			&activite.Detail, &activite.Creation, &activite.Nom, &activite.Prenom); erreur != nil {
			return nil, erreur
		}
		activites = append(activites, activite)
	}
	return activites, lignes.Err()
}

func (d *Depot) Commentaires(ctx context.Context, tache string) ([]modeles.Commentaire, error) {
	lignes, erreur := d.bd.Query(ctx, `
		SELECT c.id, c.tache, c.auteur, c.contenu, c.creation, u.nom, u.prenom
		FROM commentaires c JOIN utilisateurs u ON u.id = c.auteur
		WHERE c.tache = $1 ORDER BY c.creation`, tache)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	commentaires := []modeles.Commentaire{}
	for lignes.Next() {
		var commentaire modeles.Commentaire
		if erreur := lignes.Scan(&commentaire.ID, &commentaire.Tache, &commentaire.Auteur,
			&commentaire.Contenu, &commentaire.Creation, &commentaire.Nom, &commentaire.Prenom); erreur != nil {
			return nil, erreur
		}
		commentaires = append(commentaires, commentaire)
	}
	return commentaires, lignes.Err()
}

func (d *Depot) CreerCommentaire(ctx context.Context, commentaire modeles.Commentaire) (*modeles.Commentaire, error) {
	erreur := d.bd.QueryRow(ctx, `
		INSERT INTO commentaires (tache, auteur, contenu) VALUES ($1, $2, $3)
		RETURNING id, creation`,
		commentaire.Tache, commentaire.Auteur, commentaire.Contenu).
		Scan(&commentaire.ID, &commentaire.Creation)
	if erreur != nil {
		return nil, erreur
	}
	erreur = d.bd.QueryRow(ctx, `SELECT nom, prenom FROM utilisateurs WHERE id = $1`, commentaire.Auteur).
		Scan(&commentaire.Nom, &commentaire.Prenom)
	return &commentaire, erreur
}

func (d *Depot) SupprimerCommentaire(ctx context.Context, id, auteur string) error {
	_, erreur := d.bd.Exec(ctx, `DELETE FROM commentaires WHERE id = $1 AND auteur = $2`, id, auteur)
	return erreur
}
