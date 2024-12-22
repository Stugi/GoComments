package model

import "time"

type Comment struct {
	ID       int       `json:"id"`
	NewsID   *int      `json:"news_id"`
	Text     string    `json:"text"`
	ParentID *int      `json:"parent_id"`
	Created  time.Time `json:"created"`
}
