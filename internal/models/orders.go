package models

import "time"

type Status int

const (
	Pending Status = iota
	Started
	Finished
	Delivery
	Canceled
)

type Order struct {
	ID      int
	Name    string
	Email   string
	Message string
	Phone   string
	Status
	Price     float64
	CreatedAt time.Time
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
	switch o.Status {
	case Pending:
		return "Pendente"
	case Started:
		return "Em Andamento"
	case Finished:
		return "Finalizado"
	case Delivery:
		return "Entregue"
	case Canceled:
		return "Cancelado"
	default:
		return "Desconhecido"
	}
}

func (o Order) Timestamp() string {
	return o.CreatedAt.Local().Format("Jan 2, 3:04pm")
}

func CreateOrder(o *Order, c *Cart) error {
	res, err := db.Exec(`
	INSERT INTO orders (name, email, message, phone, price, created_at)
	VALUES(?, ?, ?, ?, ?, ?)`, o.Name, o.Email, o.Message, o.Phone, c.Total(), time.Now())
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	for _, i := range c.Items {
		_, err := db.Exec(`
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES(?, ?, ?, ?)
		`, id, i.Product.ID, i.Quantity, i.Product.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

func FindOrder(id string) (*Order, error) {
	o := &Order{}
	err := db.QueryRow(`
	SELECT id, name, email, phone, message, status, price, created_at
	FROM orders WHERE id=?`, id).Scan(&o.ID, &o.Name, &o.Email, &o.Phone,
		&o.Message, &o.Status, &o.Price, &o.CreatedAt)
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

func AllOrders() (orders []Order, err error) {
	q := `SELECT id, name, email, phone, message, status, price, created_at
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
		err = rows.Scan(&o.ID, &o.Name, &o.Email, &o.Phone, &o.Message, &o.Status,
			&o.Price, &o.CreatedAt)
		if err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	err = rows.Err()
	return orders, err
}
