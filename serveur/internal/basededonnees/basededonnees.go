package basededonnees

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schema string

func Connecter(ctx context.Context, url string) (*pgxpool.Pool, error) {
	var reserve *pgxpool.Pool
	var erreur error
	for tentative := 0; tentative < 20; tentative++ {
		reserve, erreur = pgxpool.New(ctx, url)
		if erreur == nil {
			erreur = reserve.Ping(ctx)
			if erreur == nil {
				break
			}
			reserve.Close()
		}
		time.Sleep(3 * time.Second)
	}
	if erreur != nil {
		return nil, fmt.Errorf("connexion à la base impossible : %w", erreur)
	}
	if _, erreur = reserve.Exec(ctx, schema); erreur != nil {
		return nil, fmt.Errorf("application du schéma impossible : %w", erreur)
	}
	return reserve, nil
}
