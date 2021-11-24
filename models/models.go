package models

import (
	"database/sql"
	"time"
)

type Models struct {
	DB DBModel
}

// NewModels returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Product struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	//Image           string         `json:"image"`
	CreatedAt       time.Time      `json:"-"`
	UpdatedAt       time.Time      `json:"-"`
	ProductCategory map[int]string `json:"categories"`
}

type Category struct {
	ID           int       `json:"-"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type ProductCategory struct {
	ID         int       `json:"-"`
	ProductID  int       `json:"-"`
	CategoryID int       `json:"-"`
	Category   Category  `json:"category"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
