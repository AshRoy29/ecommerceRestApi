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
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Price           int            `json:"price"`
	Size            []string       `json:"size"`
	Description     string         `json:"description"`
	Image           string         `json:"image"`
	Stock           int            `json:"stock"`
	Shipping        bool           `json:"shipping"`
	CreatedAt       time.Time      `json:"-"`
	UpdatedAt       time.Time      `json:"-"`
	ProductCategory map[int]string `json:"categories"`
}

type Size struct {
	ID        int       `json:"id"`
	SizeName  string    `json:"size_name"`
	SizeStock int       `json:"size_stock"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
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
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	AccessLevel int       `json:"access_level"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
