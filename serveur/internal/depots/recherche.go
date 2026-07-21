package depots

import (
	"context"
	"sort"

	"historykanban/serveur/internal/modeles"
)

func (d *Depot) Rechercher(ctx context.Context, utilisateur, texte string, limite int) ([]modeles.Resultat, error) {
	lignes, erreur := d.bd.Query(ctx, `
		WITH accessibles AS (
			SELECT p.id, p.nom, p.couleur, p.groupe, g.nom AS groupenom
			FROM projets p
			JOIN groupes g ON g.id = p.groupe
			JOIN membres m ON m.groupe = p.groupe AND m.utilisateur = $1
		),
		correspondances AS (
			SELECT t.id, 'titre' AS origine, t.titre AS extrait, t.modification
			FROM taches t JOIN accessibles a ON a.id = t.projet
			WHERE t.suppression IS NULL AND t.titre ILIKE '%' || $2 || '%'
			UNION
			SELECT t.id, 'description', t.description, t.modification
			FROM taches t JOIN accessibles a ON a.id = t.projet
			WHERE t.suppression IS NULL AND t.description ILIKE '%' || $2 || '%'
			UNION
			SELECT t.id, 'commentaire', c.contenu, c.creation
			FROM commentaires c
			JOIN taches t ON t.id = c.tache
			JOIN accessibles a ON a.id = t.projet
			WHERE t.suppression IS NULL AND c.contenu ILIKE '%' || $2 || '%'
		)
		SELECT DISTINCT ON (t.id)
			t.id, t.titre, t.urgence, t.echeance,
			a.id, a.nom, a.couleur, a.groupe, a.groupenom,
			col.nom, m.origine, m.extrait, m.modification
		FROM correspondances m
		JOIN taches t ON t.id = m.id
		JOIN accessibles a ON a.id = t.projet
		JOIN colonnes col ON col.id = t.colonne
		ORDER BY t.id, CASE m.origine WHEN 'titre' THEN 0 WHEN 'description' THEN 1 ELSE 2 END`,
		utilisateur, texte)
	if erreur != nil {
		return nil, erreur
	}
	defer lignes.Close()
	resultats := []modeles.Resultat{}
	for lignes.Next() {
		var resultat modeles.Resultat
		if erreur := lignes.Scan(&resultat.Tache, &resultat.Titre, &resultat.Urgence, &resultat.Echeance,
			&resultat.Projet, &resultat.ProjetNom, &resultat.ProjetCouleur, &resultat.Groupe, &resultat.GroupeNom,
			&resultat.Colonne, &resultat.Origine, &resultat.Extrait, &resultat.Modification); erreur != nil {
			return nil, erreur
		}
		resultats = append(resultats, resultat)
	}
	if erreur := lignes.Err(); erreur != nil {
		return nil, erreur
	}
	sort.Slice(resultats, func(gauche, droite int) bool {
		if rangOrigine(resultats[gauche].Origine) != rangOrigine(resultats[droite].Origine) {
			return rangOrigine(resultats[gauche].Origine) < rangOrigine(resultats[droite].Origine)
		}
		return resultats[gauche].Modification.After(resultats[droite].Modification)
	})
	if len(resultats) > limite {
		resultats = resultats[:limite]
	}
	return resultats, nil
}

func rangOrigine(origine string) int {
	switch origine {
	case "titre":
		return 0
	case "description":
		return 1
	default:
		return 2
	}
}
