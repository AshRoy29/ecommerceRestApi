package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get returns one movie and error, if any
func (m *DBModel) Get(id int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, price, description, image,
				created_at, updated_at from products where id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var product Product

	err := row.Scan(
		&product.ID,
		&product.Title,
		&product.Price,
		&product.Description,
		//&product.Image,
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
func (m *DBModel) All() ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, price, description, image,
				created_at, updated_at from products order by title
	`

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
			&product.Description,
			//&product.Image,
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

func (m *DBModel) InsertProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into products (title, price, description, created_at, updated_at) 
			values ($1, $2, $3, $4, $5)`

	_, err := m.DB.ExecContext(ctx, stmt,
		product.Title,
		product.Price,
		product.Description,
		product.CreatedAt,
		product.UpdatedAt,
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
