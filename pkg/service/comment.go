package service

import (
	"context"
	"stugi/go-comment/pkg/model"
	"stugi/go-comment/pkg/storage"
)

type Service struct {
	db storage.DB
}

func New(db storage.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) AddComment(ctx context.Context, comment model.Comment) (int, error) {
	return s.db.AddComment(comment)
}

func (s *Service) GetCommentsByNewsID(ctx context.Context, newsID int) ([]*model.Comment, error) {
	return s.db.GetCommentsByNews(newsID)
}
