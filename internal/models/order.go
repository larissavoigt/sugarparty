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