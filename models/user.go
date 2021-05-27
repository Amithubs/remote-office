package models

import (
	"github.com/volatiletech/null"
	"time"
)
type UserPermission string

const (
	Employee UserPermission="employee"
	Admin UserPermission="admin"
)
type User struct {
	ID       int         `json:"id" db:"id"`
	Name     string      `json:"name" db:"name"`
	Phone    null.String `json:"phone" db:"phone"`
	Email    string      `json:"email" db:"email"`
	Position string      `json:"position" db:"position"`

	ProfileImageID   null.Int `json:"-" db:"profile_image"`
	ProfileImageLink string   `json:"profileImageLink" db:"-"`

	CreatedAt time.Time `json:"-" db:"created_at"`
	Permissions []UserPermission `json:"permission" db:"-"`
	Token string `json:"-" db:"-"`
}
