package stockage

import (
	"context"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Stockage struct {
	client       *minio.Client
	seau         string
	urlPublique  string
}

func Nouveau(hote, cle, secret, seau, urlPublique string) (*Stockage, error) {
	client, erreur := minio.New(hote, &minio.Options{
		Creds:  credentials.NewStaticV4(cle, secret, ""),
		Secure: false,
	})
	if erreur != nil {
		return nil, erreur
	}
	return &Stockage{
		client:      client,
		seau:        seau,
		urlPublique: strings.TrimRight(urlPublique, "/"),
	}, nil
}

func (s *Stockage) Televerser(ctx context.Context, chemin string, contenu io.Reader, taille int64, typeContenu string) error {
	_, erreur := s.client.PutObject(ctx, s.seau, chemin, contenu, taille, minio.PutObjectOptions{ContentType: typeContenu})
	return erreur
}

func (s *Stockage) Supprimer(ctx context.Context, chemin string) error {
	return s.client.RemoveObject(ctx, s.seau, chemin, minio.RemoveObjectOptions{})
}

func (s *Stockage) URL(chemin string) string {
	return s.urlPublique + "/" + s.seau + "/" + chemin
}
