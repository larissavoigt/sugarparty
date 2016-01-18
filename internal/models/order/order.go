package order

import (
	"strconv"
	"time"

	"github.com/larissavoigt/sugarparty/internal/models/cart"
	"github.com/larissavoigt/sugarparty/internal/models/db"
)

type Status int

const (
	Pending Status = iota
	Started
	Finished
	Delivery
	Canceled
)

var StatusNames = []string{
	"Pendente",
	"Em Andamento",
	"Finalizado",
	"Entregue",
	"Cancelado",
}

type Order struct {
	ID      int
	Name    string
	Email   string
	Message string
	Phone   string
	Status
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []OrderItem
}

type OrderItem struct {
	ID          int
	OrderID     int
	ProductID   int
	ProductName string
	Quantity    int
	Price       float64
}

func (o Order) StatusName() string {
	i := int(o.Status)
	if i < len(StatusNames) {
		return StatusNames[i]
	} else {
		return "Desconhecido"
	}
}

func Create(o *Order, c *cart.Cart) (string, error) {
	res, err := db.Exec(`
	INSERT INTO orders (name, email, message, phone, price, created_at,
	updated_at)
	VALUES(?, ?, ?, ?, ?, ?, ?)`, o.Name, o.Email, o.Message, o.Phone, c.Total(),
		time.Now(), time.Now())
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	for _, i := range c.Items {
		_, err := db.Exec(`
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES(?, ?, ?, ?)
		`, id, i.Product.ID, i.Quantity, i.Product.Price)
		if err != nil {
			return "", err
		}
	}
	return strconv.FormatInt(id, 10), nil
}

func Update(id string, status int) error {
	_, err := db.Exec(`
	UPDATE orders SET status=?, updated_at = ? where id = ?`, status, time.Now(),
		id)
	return err
}

func Find(id string) (*Order, error) {
	o := &Order{}
	err := db.QueryRow(`
	SELECT id, name, email, phone, message, status, price, created_at, updated_at
	FROM orders WHERE id=?`, id).Scan(&o.ID, &o.Name, &o.Email, &o.Phone,
		&o.Message, &o.Status, &o.Price, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return o, err
	}
	q := `SELECT p.name, i.quantity, i.price
	FROM order_items i
	INNER JOIN products p
	ON p.id = i.product_id
	WHERE i.order_id = ?
	`

	rows, err := db.Query(q, o.ID)

	if err != nil {
		return o, err
	}
	defer rows.Close()
	for rows.Next() {
		i := OrderItem{}
		err = rows.Scan(&i.ProductName, &i.Quantity, &i.Price)
		if err != nil {
			return o, err
		}
		o.Items = append(o.Items, i)
	}
	err = rows.Err()
	return o, err
}

func All() (orders []Order, err error) {
	q := `SELECT id, name, status, created_at, updated_at
	FROM orders p
	ORDER BY id DESC
	`
	rows, err := db.Query(q)

	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		o := Order{}
		err = rows.Scan(&o.ID, &o.Name, &o.Status, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	err = rows.Err()
	return orders, err
}
