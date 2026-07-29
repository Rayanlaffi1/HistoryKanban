package depots

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func depotTest(t *testing.T) (*Depot, context.Context, func()) {
	t.Helper()
	adresse := os.Getenv("HK_TEST_BDURL")
	if adresse == "" {
		t.Skip("HK_TEST_BDURL non defini")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, adresse)
	if err != nil {
		t.Fatalf("connexion base de test: %v", err)
	}
	return Nouveau(pool), ctx, pool.Close
}

func preparerProjetStatistiques(t *testing.T, depot *Depot, ctx context.Context) (groupe string, projet string, todo string, fini string, utilisateur string) {
	t.Helper()
	suffixe := time.Now().UnixNano()
	utilisateur = "00000000-0000-4000-8000-000000000001"
	groupe = "00000000-0000-4000-8000-000000000002"
	projet = "00000000-0000-4000-8000-000000000003"
	todo = "00000000-0000-4000-8000-000000000004"
	fini = "00000000-0000-4000-8000-000000000005"

	_, err := depot.bd.Exec(ctx, `DELETE FROM groupes WHERE id = $1`, groupe)
	if err != nil {
		t.Fatalf("nettoyage groupe: %v", err)
	}
	_, err = depot.bd.Exec(ctx, `DELETE FROM utilisateurs WHERE id = $1`, utilisateur)
	if err != nil {
		t.Fatalf("nettoyage utilisateur: %v", err)
	}
	_, err = depot.bd.Exec(ctx, `
		INSERT INTO utilisateurs (id, courriel, nom, prenom)
		VALUES ($1, $2, 'Test', 'Stats')`, utilisateur, "stats-test+"+time.Unix(0, suffixe).Format("150405.000000000")+"@historykanban.local")
	if err != nil {
		t.Fatalf("creation utilisateur: %v", err)
	}
	_, err = depot.bd.Exec(ctx, `
		INSERT INTO groupes (id, nom, proprietaire) VALUES ($1, 'Groupe stats', $2)`, groupe, utilisateur)
	if err != nil {
		t.Fatalf("creation groupe: %v", err)
	}
	_, err = depot.bd.Exec(ctx, `
		INSERT INTO membres (groupe, utilisateur, role) VALUES ($1, $2, 'proprietaire')`, groupe, utilisateur)
	if err != nil {
		t.Fatalf("creation membre: %v", err)
	}
	_, err = depot.bd.Exec(ctx, `
		INSERT INTO projets (id, groupe, nom, createur) VALUES ($1, $2, 'Projet stats', $3)`, projet, groupe, utilisateur)
	if err != nil {
		t.Fatalf("creation projet: %v", err)
	}
	_, err = depot.bd.Exec(ctx, `
		INSERT INTO colonnes (id, projet, nom, position, terminale) VALUES
		($1, $3, 'A faire', 0, false),
		($2, $3, 'Termine', 1, true)`, todo, fini, projet)
	if err != nil {
		t.Fatalf("creation colonnes: %v", err)
	}
	t.Cleanup(func() {
		_, _ = depot.bd.Exec(ctx, `DELETE FROM groupes WHERE id = $1`, groupe)
		_, _ = depot.bd.Exec(ctx, `DELETE FROM utilisateurs WHERE id = $1`, utilisateur)
	})
	return groupe, projet, todo, fini, utilisateur
}

func TestStatistiquesUtilisentLaDateDeFinEtPasLaDerniereModification(t *testing.T) {
	depot, ctx, fermer := depotTest(t)
	defer fermer()
	_, projet, _, fini, utilisateur := preparerProjetStatistiques(t, depot, ctx)

	terminee := time.Now().AddDate(0, 0, -5).Truncate(time.Second)
	modifieeApresFin := time.Now().Truncate(time.Second)
	_, err := depot.bd.Exec(ctx, `
		INSERT INTO taches (projet, colonne, titre, points, createur, terminee, modification)
		VALUES ($1, $2, 'Tache terminee puis modifiee', 8, $3, $4, $5)`, projet, fini, utilisateur, terminee, modifieeApresFin)
	if err != nil {
		t.Fatalf("creation tache terminee: %v", err)
	}

	stats, err := depot.Statistiques(ctx, []string{projet}, terminee.Add(-24*time.Hour), terminee.Add(24*time.Hour), "day", 100, 0)
	if err != nil {
		t.Fatalf("statistiques: %v", err)
	}
	if stats.Totaux.Terminees != 1 || stats.Totaux.Points != 8 {
		t.Fatalf("statistiques terminees = %d/%d, attendu 1/8", stats.Totaux.Terminees, stats.Totaux.Points)
	}
}

func TestStatistiquesListentLesTachesTermineesParPeriodeTousProjets(t *testing.T) {
	depot, ctx, fermer := depotTest(t)
	defer fermer()
	_, premierProjet, _, premiereFin, utilisateur := preparerProjetStatistiques(t, depot, ctx)

	secondProjet := "00000000-0000-4000-8000-000000000006"
	secondeFin := "00000000-0000-4000-8000-000000000007"
	_, err := depot.bd.Exec(ctx, `
		INSERT INTO projets (id, groupe, nom, couleur, createur) VALUES ($1, $2, 'Second projet stats', '#2563eb', $3)
		ON CONFLICT (id) DO NOTHING`, secondProjet, "00000000-0000-4000-8000-000000000002", utilisateur)
	if err != nil {
		t.Fatalf("creation second projet: %v", err)
	}
	_, err = depot.bd.Exec(ctx, `
		INSERT INTO colonnes (id, projet, nom, position, terminale) VALUES ($1, $2, 'Termine', 0, true)
		ON CONFLICT (id) DO NOTHING`, secondeFin, secondProjet)
	if err != nil {
		t.Fatalf("creation colonne second projet: %v", err)
	}

	periode := time.Now().AddDate(0, 0, -2).Truncate(time.Second)
	_, err = depot.bd.Exec(ctx, `
		INSERT INTO taches (projet, colonne, titre, points, urgence, createur, terminee)
		VALUES
		($1, $2, 'Tache premier projet', 5, 'normale', $5, $6),
		($3, $4, 'Tache second projet', 3, 'urgente', $5, $6)`, premierProjet, premiereFin, secondProjet, secondeFin, utilisateur, periode)
	if err != nil {
		t.Fatalf("creation taches terminees: %v", err)
	}

	stats, err := depot.Statistiques(ctx, []string{premierProjet, secondProjet}, periode.Add(-24*time.Hour), periode.Add(24*time.Hour), "day", 100, 0)
	if err != nil {
		t.Fatalf("statistiques: %v", err)
	}
	if len(stats.Terminees) != 1 {
		t.Fatalf("periodes terminees = %d, attendu 1", len(stats.Terminees))
	}
	periodeTerminee := stats.Terminees[0]
	if periodeTerminee.Points != 8 || periodeTerminee.Total != 2 || len(periodeTerminee.Taches) != 2 {
		t.Fatalf("periode terminee = %+v, attendu 8 points et 2 taches", periodeTerminee)
	}
	if periodeTerminee.Taches[0].Titre != "Tache premier projet" || periodeTerminee.Taches[0].ProjetNom == "" {
		t.Fatalf("premiere tache terminee invalide: %+v", periodeTerminee.Taches[0])
	}
	if periodeTerminee.Taches[1].Titre != "Tache second projet" || periodeTerminee.Taches[1].Projet != secondProjet {
		t.Fatalf("seconde tache terminee invalide: %+v", periodeTerminee.Taches[1])
	}
}

func TestDeplacerTacheRenseigneLaDateDeFinEnColonneTerminale(t *testing.T) {
	depot, ctx, fermer := depotTest(t)
	defer fermer()
	_, projet, todo, fini, utilisateur := preparerProjetStatistiques(t, depot, ctx)

	var tache string
	err := depot.bd.QueryRow(ctx, `
		INSERT INTO taches (projet, colonne, titre, points, createur)
		VALUES ($1, $2, 'Tache a terminer', 3, $3)
		RETURNING id`, projet, todo, utilisateur).Scan(&tache)
	if err != nil {
		t.Fatalf("creation tache: %v", err)
	}

	if err := depot.DeplacerTache(ctx, tache, fini, 0); err != nil {
		t.Fatalf("deplacement vers termine: %v", err)
	}
	var terminee *time.Time
	if err := depot.bd.QueryRow(ctx, `SELECT terminee FROM taches WHERE id = $1`, tache).Scan(&terminee); err != nil {
		t.Fatalf("lecture date de fin: %v", err)
	}
	if terminee == nil {
		t.Fatal("date de fin absente apres deplacement en colonne terminale")
	}
}

func TestDeplacerTacheConserveLaDateDeFinApresReouverture(t *testing.T) {
	depot, ctx, fermer := depotTest(t)
	defer fermer()
	_, projet, todo, fini, utilisateur := preparerProjetStatistiques(t, depot, ctx)

	var tache string
	err := depot.bd.QueryRow(ctx, `
		INSERT INTO taches (projet, colonne, titre, points, createur)
		VALUES ($1, $2, 'Tache a rouvrir', 3, $3)
		RETURNING id`, projet, todo, utilisateur).Scan(&tache)
	if err != nil {
		t.Fatalf("creation tache: %v", err)
	}

	// Premiere completion : la date de fin est renseignee.
	if err := depot.DeplacerTache(ctx, tache, fini, 0); err != nil {
		t.Fatalf("deplacement vers termine: %v", err)
	}
	var premiereFin *time.Time
	if err := depot.bd.QueryRow(ctx, `SELECT terminee FROM taches WHERE id = $1`, tache).Scan(&premiereFin); err != nil {
		t.Fatalf("lecture date de fin: %v", err)
	}
	if premiereFin == nil {
		t.Fatal("date de fin absente apres deplacement en colonne terminale")
	}

	// Reouverture : la tache quitte la colonne terminale mais garde son historique.
	if err := depot.DeplacerTache(ctx, tache, todo, 0); err != nil {
		t.Fatalf("deplacement hors termine: %v", err)
	}
	var apresReouverture *time.Time
	if err := depot.bd.QueryRow(ctx, `SELECT terminee FROM taches WHERE id = $1`, tache).Scan(&apresReouverture); err != nil {
		t.Fatalf("lecture date de fin apres reouverture: %v", err)
	}
	if apresReouverture == nil {
		t.Fatal("date de fin effacee apres reouverture : l'historique de completion est perdu")
	}
	if !apresReouverture.Equal(*premiereFin) {
		t.Fatalf("date de fin modifiee a la reouverture : %v puis %v", premiereFin, apresReouverture)
	}

	// Retour en colonne terminale : la premiere date de completion est conservee.
	if err := depot.DeplacerTache(ctx, tache, fini, 0); err != nil {
		t.Fatalf("re-deplacement vers termine: %v", err)
	}
	var apresReprise *time.Time
	if err := depot.bd.QueryRow(ctx, `SELECT terminee FROM taches WHERE id = $1`, tache).Scan(&apresReprise); err != nil {
		t.Fatalf("lecture date de fin apres reprise: %v", err)
	}
	if apresReprise == nil || !apresReprise.Equal(*premiereFin) {
		t.Fatalf("premiere date de completion non conservee au retour en colonne terminale : %v puis %v", premiereFin, apresReprise)
	}
}

func TestStatistiquesUtilisentLaColonneTerminaleExpliciteMalgreUneColonneSuivante(t *testing.T) {
	depot, ctx, fermer := depotTest(t)
	defer fermer()
	_, projet, _, fini, utilisateur := preparerProjetStatistiques(t, depot, ctx)

	// Une colonne « Archive » est ajoutee apres « Termine » : elle occupe desormais
	// la position maximale sans etre la colonne terminale.
	archive := "00000000-0000-4000-8000-000000000008"
	if _, err := depot.bd.Exec(ctx, `
		INSERT INTO colonnes (id, projet, nom, position, terminale) VALUES ($1, $2, 'Archive', 2, false)
		ON CONFLICT (id) DO NOTHING`, archive, projet); err != nil {
		t.Fatalf("creation colonne archive: %v", err)
	}

	// La tache reste dans « Termine » (colonne terminale) et non dans « Archive ».
	terminee := time.Now().AddDate(0, 0, -1).Truncate(time.Second)
	if _, err := depot.bd.Exec(ctx, `
		INSERT INTO taches (projet, colonne, titre, points, createur, terminee)
		VALUES ($1, $2, 'Tache terminee non archivee', 5, $3, $4)`, projet, fini, utilisateur, terminee); err != nil {
		t.Fatalf("creation tache terminee: %v", err)
	}

	stats, err := depot.Statistiques(ctx, []string{projet}, terminee.Add(-24*time.Hour), terminee.Add(24*time.Hour), "day", 100, 0)
	if err != nil {
		t.Fatalf("statistiques: %v", err)
	}
	// Avec un reperage par max(position), la tache serait ignoree car « Archive » est
	// plus a droite ; la colonne terminale explicite doit au contraire la comptabiliser.
	if stats.Totaux.Terminees != 1 || stats.Totaux.Points != 5 {
		t.Fatalf("statistiques terminees = %d/%d, attendu 1/5 via la colonne terminale explicite", stats.Totaux.Terminees, stats.Totaux.Points)
	}
}
