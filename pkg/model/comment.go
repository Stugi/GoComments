package model

import "time"

type Comment struct {
	ID       int       `json:"id"`
	NewsID   *int      `json:"news_id"`
	Text     string    `json:"text"`
	ParentID *int      `json:"parent_id"`
	Status   int       `json:"status"`
	Created  time.Time `json:"created"`
}

const (
	StatusNew      = 0
	StatusApproved = 1
	StatusBlocked  = 2
)
