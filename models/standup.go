package models

import "time"

type StandUp struct {
	ID int `json:"id" db:"id"`
	UserID int `json:"-" db:"user_id"`
	Data string `json:"data" db:"data"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

