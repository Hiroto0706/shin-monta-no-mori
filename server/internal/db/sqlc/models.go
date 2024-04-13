// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Character struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Src       string    `json:"src"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ChildCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  int64     `json:"parent_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Image struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	// 文字ありの画像
	OriginalSrc string `json:"original_src"`
	// 文字無しの画像.オリジナルが文字無しの時もあるので、nullableにする.
	SimpleSrc sql.NullString `json:"simple_src"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
}

type ImageCharactersRelation struct {
	ID          int64 `json:"id"`
	ImageID     int64 `json:"image_id"`
	CharacterID int64 `json:"character_id"`
}

type ImageParentCategoriesRelation struct {
	ID               int64 `json:"id"`
	ImageID          int64 `json:"image_id"`
	ParentCategoryID int64 `json:"parent_category_id"`
}

type Operator struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	HashedPassword string    `json:"hashed_password"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
}

type ParentCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Src       string    `json:"src"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
