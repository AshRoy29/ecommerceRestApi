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

	query := `select id, title, price, size, description, image, stock, shipping,
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
		&product.Image,
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

	query := fmt.Sprintf(`select id, title, price, size, description, image, stock, shipping,
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
			pq.Array(&product.Size),
			&product.Description,
			&product.Image,
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

	stmt := `update products set title = $1, price = $2, size = $3, description = $4, image = $5, stock = $6, shipping = $7, updated_at = $8 
			where id = $9`

	_, err := m.DB.ExecContext(ctx, stmt,
		product.Title,
		product.Price,
		pq.Array(product.Size),
		product.Description,
		product.Image,
		product.Stock,
		product.Shipping,
		product.UpdatedAt,
		product.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) UpdateCategory(pc ProductCategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update products_category set category_id = $1, updated_at = $2
			where product_id = $3`

	_, err := m.DB.ExecContext(ctx, stmt,
		pc.CategoryID,
		pc.UpdatedAt,
		pc.ProductID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteProduct(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var image string

	query := `select image from products where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&image,
	)

	stmt := `delete from products where id = $1`

	_, err = m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return "", err
	}
	return image, nil
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

func (m *DBModel) NewUser(u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into users (first_name, last_name, phone, email, password, access_level, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.FirstName,
		u.LastName,
		u.Phone,
		u.Email,
		u.Password,
		"user",
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) ValidUser(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, phone, email, password, access_level
			from users
			where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	//var password string
	//var level string
	var user User

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
	//&user.AccessLevel,
	)
	if err != nil {
		return nil, err
	}

	return &user, err

}

func (m *DBModel) CheckEmail(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select email from users where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var userEmail string

	err := row.Scan(
		&userEmail,
	)
	if err != nil {
		return "", err
	}

	return userEmail, nil
}

func (m *DBModel) CartOrders(cp CartProducts) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into orders (product_id, product_size, product_price, quantity, user_id, total)
			values ($1, $2, $3, $4, $5, $6) returning user_id, id`

	var userID int
	var orderID int

	err := m.DB.QueryRowContext(ctx, stmt,
		pq.Array(cp.ProductID),
		pq.Array(cp.Size),
		pq.Array(cp.Price),
		pq.Array(cp.Quantity),
		cp.UserID,
		cp.Total,
	).Scan(&userID, &orderID)
	if err != nil {
		return 0, 0, err
	}

	return userID, orderID, nil
}

func (m *DBModel) BillingInfo(b BillingInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into billing_info (name, phone, address, postal_code, city, user_id, created_at, updated_at, order_id)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := m.DB.ExecContext(ctx, stmt,
		b.Name,
		b.Phone,
		b.Address,
		b.PostalCode,
		b.City,
		b.UserID,
		time.Now(),
		time.Now(),
		b.OrderID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) AllOrders() ([]*CartProducts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select
						o.id, o.product_id, o.product_size, o.product_price, o.quantity, o.total,
						bi.name, bi.phone, bi.address, bi.postal_code, bi.city, bi.user_id, bi.created_at,
						u.first_name, u.last_name, u.phone, u.email
					from orders o
					left join billing_info bi on (bi.order_id = o.id)
					left join users u on (u.id = bi.user_id)`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*CartProducts

	for rows.Next() {
		var order CartProducts

		err := rows.Scan(
			&order.ID,
			pq.Array(&order.ProductID),
			pq.Array(&order.Size),
			pq.Array(&order.Price),
			pq.Array(&order.Quantity),
			&order.Total,
			&order.BillingInfo.Name,
			&order.BillingInfo.Phone,
			&order.BillingInfo.Address,
			&order.BillingInfo.PostalCode,
			&order.BillingInfo.City,
			&order.BillingInfo.UserID,
			&order.BillingInfo.CreatedAt,
			&order.User.FirstName,
			&order.User.LastName,
			&order.User.Phone,
			&order.User.Email,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}
