package product

import (
	"fmt"
	"strings"

	"github.com/luizbranco/sugarparty/internal/models/category"
	"github.com/luizbranco/sugarparty/internal/models/db"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Active      bool
	Category    category.Category
}

func Create(p *Product) error {
	_, err := db.Exec(`
	INSERT INTO products (name, description, price, active, category_id)
	VALUES(?, ?, ?, ?, ?)`, p.Name, p.Description, p.Price, p.Active, p.Category.ID)
	return err
}

func Update(p *Product) error {
	_, err := db.Exec(`
	UPDATE products SET name=?, description=?, price=?, active=?, category_id=?
	where id = ?`, p.Name, p.Description, p.Price, p.Active, p.Category.ID, p.ID)
	return err
}

func Find(id string) (*Product, error) {
	p := &Product{}
	err := db.QueryRow(`
	SELECT id, name, description, price, active, category_id
	FROM products WHERE id=?`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Active, &p.Category.ID)
	return p, err
}

func FindAll(ids []interface{}) ([]Product, error) {
	args := strings.Repeat("?,", len(ids))
	q := fmt.Sprintf(`SELECT p.id, p.name, p.description, p.price, p.active, c.name
	FROM products p
	INNER JOIN categories c
	ON p.category_id = c.id
	WHERE p.id IN (%s)
	`, args[:len(args)-1])
	return scan(q, ids...)
}

func All() ([]Product, error) {
	return scan(`
	SELECT p.id, p.name, p.description, p.price, p.active, c.name
	FROM products p
	INNER JOIN categories c
	ON p.category_id = c.id
	`)
}

func Active(id string) ([]Product, error) {
	return scan(`
	SELECT id, name, description, price, active, ''
	FROM products
	WHERE category_id = ? AND active IS TRUE
	`, id)
}

func scan(query string, args ...interface{}) (products []Product, err error) {
	rows, err := db.Query(query, args...)

	if err != nil {
		return products, err
	}
	defer rows.Close()
	for rows.Next() {
		p := Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Active, &p.Category.Name)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}
	err = rows.Err()
	return products, err
}
