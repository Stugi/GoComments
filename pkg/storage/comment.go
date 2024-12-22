package storage

import (
	"context"
	"fmt"
	"strings"
	"stugi/go-comment/pkg/model"
)

func (s *storage) GetComments(filter map[string]any, limit int) ([]*model.Comment, error) {
	var (
		whereClauses []string
		args         []any
	)

	// Обрабатываем фильтр
	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, len(args)+2)) // Нумерация параметров начинается с 1
		args = append(args, value)
	}

	sql :=
		`SELECT 
			id, 
			news_id, 
			text, 
			parent_id, 
			status,
			created
		FROM comment`
	if len(whereClauses) > 0 {
		sql += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	sql += " LIMIT $1"                   // LIMIT всегда параметр $1
	args = append([]any{limit}, args...) // Добавляем LIMIT как первый параметр

	rows, err := s.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(
			&comment.ID,
			&comment.NewsID,
			&comment.Text,
			&comment.ParentID,
			&comment.Status,
			&comment.Created,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return comments, nil
}

func (s *storage) AddComment(comment model.Comment) (int, error) {
	var sql string
	var args []any
	if comment.ParentID == nil {
		sql = `INSERT INTO comment ( text, parent_id, status) VALUES ($1, $2, $3) RETURNING id`
		args = []any{comment.Text, comment.ParentID, comment.Status}
	} else {
		sql = `INSERT INTO comment (news_id, text, status) VALUES ($1, $2, $3) RETURNING id`
		args = []any{comment.NewsID, comment.Text, comment.Status}
	}

	var id int
	err := s.db.QueryRow(
		context.Background(),
		sql,
		args...,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *storage) GetCommentsByNews(newsID int) ([]*model.Comment, error) {
	sql := `
	SELECT 
		id, news_id, text, parent_id, status, created 
	FROM comment 
	WHERE status = $1 AND 
	news_id = $2`
	rows, err := s.db.Query(context.Background(), sql, model.StatusApproved, newsID)
	if err != nil {
		return nil, err
	}

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(
			&comment.ID,
			&comment.NewsID,
			&comment.Text,
			&comment.ParentID,
			&comment.Status,
			&comment.Created,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

// TODO тут бы RWMutex нужен
func (s *storage) MarkComment(commentID int, status int) error {
	sql := `UPDATE comment SET status = $1 WHERE id = $2`
	_, err := s.db.Exec(context.Background(), sql, status, commentID)
	return err
}
