package storage

import (
	"context"
	"stugi/go-comment/pkg/cache"
	"stugi/go-comment/pkg/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

const connection = "postgres://postgres:postgres@localhost:5432/gonews?sslmode=disable"

type storage struct {
	db    *pgxpool.Pool
	cache *cache.Impl
}

type DB interface {
	AddComment(model.Comment) (int, error)
	GetCommentsByNews(int) ([]*model.Comment, error)
	MarkComment(int, int) error
	GetComments(map[string]any, int) ([]*model.Comment, error)
}

func New(cache *cache.Impl) (DB, error) {
	db, err := pgxpool.Connect(context.Background(), connection)
	if err != nil {
		return nil, err
	}
	s := &storage{
		db:    db,
		cache: cache,
	}
	return s, nil
}
