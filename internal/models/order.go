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

func CreateOrder(o *Order, c *Cart) error {
	_, err := db.Exec(`
	INSERT INTO orders (name, email, message, phone, created_at)
	VALUES(?, ?, ?, ?, ?)`, o.Name, o.Email, o.Message, o.Phone, time.Now())
	return err
}
