package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get returns one product and error, if any
func (m *DBModel) Get(id int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, price, size, description,  stock, shipping,
				created_at, updated_at from products where id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var product Product

	err := row.Scan(
		&product.ID,
		&product.Title,
		&product.Price,
		pq.Array(&product.Size),
		&product.Description,
		//&product.Image,
		&product.Stock,
		&product.Shipping,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	query = `select
				pc.id, pc.product_id, pc.category_id, c.category_name
			from
				products_category pc
				left join category c on (c.id = pc.category_id)
			where
				pc.product_id = $1
	`

	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	category := make(map[int]string)
	for rows.Next() {
		var pc ProductCategory
		err := rows.Scan(
			&pc.ID,
			&pc.ProductID,
			&pc.CategoryID,
			&pc.Category.CategoryName,
		)
		if err != nil {
			return nil, err
		}
		category[pc.ID] = pc.Category.CategoryName
	}

	product.ProductCategory = category

	return &product, nil
}

// All returns all products and error, if any
func (m *DBModel) All(category ...int) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(category) > 0 {
		where = fmt.Sprintf("where id in (select product_id from products_category where category_id = %d)", category[0])
	}

	query := fmt.Sprintf(`select id, title, price, size, description, stock, shipping,
				created_at, updated_at from products %s order by title`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ID,
			&product.Title,
			&product.Price,
			&product.Size,
			&product.Description,
			//&product.Image,
			&product.Stock,
			&product.Shipping,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// get categories
		categoryQuery := `select
			pc.id, pc.product_id, pc.category_id, c.category_name
		from
			products_category pc
			left join category c on (c.id = pc.category_id)
		where
			pc.product_id = $1
		`

		categoryRows, _ := m.DB.QueryContext(ctx, categoryQuery, product.ID)

		categories := make(map[int]string)
		for categoryRows.Next() {
			var pc ProductCategory
			err := categoryRows.Scan(
				&pc.ID,
				&pc.ProductID,
				&pc.CategoryID,
				&pc.Category.CategoryName,
			)
			if err != nil {
				return nil, err
			}
			categories[pc.CategoryID] = pc.Category.CategoryName
		}
		categoryRows.Close()

		product.ProductCategory = categories
		products = append(products, &product)

	}
	return products, nil
}

func (m *DBModel) InsertProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into products (title, price, size, description, image, stock, shipping, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var newID int
	err := m.DB.QueryRowContext(ctx, stmt,
		product.Title,
		product.Price,
		pq.Array(product.Size),
		product.Description,
		product.Image,
		product.Stock,
		product.Shipping,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	log.Println("New product ID:", newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *DBModel) InsertCategory(pc ProductCategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into products_category (product_id, category_id, created_at, updated_at)
			values ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, stmt,
		pc.ProductID,
		pc.CategoryID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) UpdateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update products set title = $1, price = $2, description = $3, updated_at = $4 
			where id = $5`

	_, err := m.DB.ExecContext(ctx, stmt,
		product.Title,
		product.Price,
		product.Description,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteProduct(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from products where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) GetAllCategory() ([]*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, category_name, created_at, updated_at
			from category order by category_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category

	for rows.Next() {
		var c Category

		err := rows.Scan(
			&c.ID,
			&c.CategoryName,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}

	return categories, nil
}
