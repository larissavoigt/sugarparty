package db

import "github.com/luizbranco/sugarparty/internal/product"

func CreateCategory(c *product.Category) error {
	_, err := db.Exec(`
	INSERT INTO categories (name, description)
	VALUES(?, ?)`, c.Name, c.Description)
	return err
}

func UpdateCategory(c *product.Category) error {
	_, err := db.Exec(`
	UPDATE categories SET name=?, description=?
	where id = ?`, c.Name, c.Description, c.ID)
	return err
}

func FindCategory(id string) (*product.Category, error) {
	c := &product.Category{}
	err := db.QueryRow(`
	SELECT id, name, description
	FROM categories WHERE id=?`, id).Scan(&c.ID, &c.Name, &c.Description)
	return c, err
}

func AllCategories() (categories []product.Category, err error) {
	rows, err := db.Query(`
	SELECT c.id, c.name, c.description, COUNT(p.category_id) "products"
	FROM categories c
	LEFT JOIN products p
	ON c.id = p.category_id
	GROUP BY c.id
	`)

	if err != nil {
		return categories, err
	}
	defer rows.Close()
	for rows.Next() {
		c := product.Category{}
		err = rows.Scan(&c.ID, &c.Name, &c.Description, &c.Count)
		if err != nil {
			return categories, err
		}
		categories = append(categories, c)
	}
	err = rows.Err()
	return categories, err
}

func ActiveCategories() (categories []product.Category, err error) {
	rows, err := db.Query(`
	SELECT c.id, c.name, c.description, COUNT(p.category_id) "products"
	FROM categories c
	LEFT JOIN products p
	ON c.id = p.category_id AND p.active is true
	GROUP BY c.id
	HAVING products > 0
	`)

	if err != nil {
		return categories, err
	}
	defer rows.Close()
	for rows.Next() {
		c := product.Category{}
		err = rows.Scan(&c.ID, &c.Name, &c.Description, &c.Count)
		if err != nil {
			return categories, err
		}
		categories = append(categories, c)
	}
	err = rows.Err()
	return categories, err
}
