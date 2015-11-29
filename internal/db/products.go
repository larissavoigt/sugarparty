package db

import "github.com/luizbranco/sugarparty/internal/product"

func CreateProduct(p *product.Product) error {
	_, err := db.Exec(`
	INSERT INTO products (name, description, price, active, category_id)
	VALUES(?, ?, ?, ?, ?)`, p.Name, p.Description, p.Price, p.Active, p.Category.ID)
	return err
}

func UpdateProduct(p *product.Product) error {
	_, err := db.Exec(`
	UPDATE products SET name=?, description=?, price=?, active=?, category_id=?
	where id = ?`, p.Name, p.Description, p.Price, p.Active, p.Category.ID, p.ID)
	return err
}

func FindProduct(id string) (*product.Product, error) {
	p := &product.Product{}
	err := db.QueryRow(`
	SELECT id, name, description, price, active, category_id
	FROM products WHERE id=?`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Active, &p.Category.ID)
	return p, err
}

func AllProducts() (products []product.Product, err error) {
	rows, err := db.Query(`
	SELECT p.id, p.name, p.description, p.price, p.active, c.name
	FROM products p
	INNER JOIN categories c
	ON p.category_id = c.id
	`)

	if err != nil {
		return products, err
	}
	defer rows.Close()
	for rows.Next() {
		p := product.Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Active, &p.Category.Name)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}
	err = rows.Err()
	return products, err
}
