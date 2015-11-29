package product

type Category struct {
	ID          string
	Name        string
	Description string
	Count       int
}

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Active      bool
	Category
}
