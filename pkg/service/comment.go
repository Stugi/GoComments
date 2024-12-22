package service

import (
	"fmt"
	"strings"
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

func (s *Service) AddComment(comment model.Comment) (int, error) {
	return s.db.AddComment(comment)
}

func (s *Service) GetCommentsByNewsID(newsID int) ([]*model.Comment, error) {
	return s.db.GetCommentsByNews(newsID)
}

var (
	// список запрещенных слов
	bannedWords = []string{"qwerty", "йцукен", "zxvbnm"}
)

func (s *Service) CheckComments() error {
	comments, _ := s.db.GetComments(map[string]any{"status": model.StatusNew}, 10)

	for _, comment := range comments {
		if checkComments(comment.Text) {
			err := s.db.MarkComment(comment.ID, model.StatusBlocked)
			if err != nil {
				fmt.Printf("Error marking comment %d as blocked: %v\n", comment.ID, err)
			}
		} else {
			err := s.db.MarkComment(comment.ID, model.StatusApproved)
			if err != nil {
				fmt.Printf("Error marking comment %d as approved: %v\n", comment.ID, err)
			}
		}
	}
	return nil
}

func checkComments(text string) bool {
	for _, word := range bannedWords {
		if strings.Contains(strings.ToLower(text), word) {
			return true
		}
	}
	return false
}
