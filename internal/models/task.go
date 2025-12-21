package models

import "time"

type Task struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
