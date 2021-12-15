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
	AccessLevel string    `json:"access_level"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type CartProducts struct {
	ID       []int    `json:"id"`
	Size     []string `json:"size"`
	Price    []int    `json:"price"`
	Quantity []int    `json:"quantity"`
	UserID   int      `json:"user"`
	Total    int      `json:"total"`
}

type BillingInfo struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	City       string    `json:"city"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
