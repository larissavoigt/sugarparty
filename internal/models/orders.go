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
	CreatedAt time.Time
	Items     []OrderItem
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
	Price     float64
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
	_, err := db.Exec(`
	INSERT INTO orders (name, email, message, phone, created_at)
	VALUES(?, ?, ?, ?, ?)`, o.Name, o.Email, o.Message, o.Phone, time.Now())
	return err
}

func AllOrders() (orders []Order, err error) {
	q := `SELECT id, name, email, phone, message, status, created_at
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
		err = rows.Scan(&o.ID, &o.Name, &o.Email, &o.Phone, &o.Message, &o.Status, &o.CreatedAt)
		if err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	err = rows.Err()
	return orders, err
}
