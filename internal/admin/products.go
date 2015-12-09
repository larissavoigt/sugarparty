package admin

import (
	"net/http"
	"strconv"

	"github.com/luizbranco/sugarparty/internal/db"
	"github.com/luizbranco/sugarparty/internal/product"
	"github.com/luizbranco/sugarparty/internal/templates"
)

func AllProducts(w http.ResponseWriter) {
	products, err := db.AllProducts()
	if err != nil {
		templates.Error(w, err)
	} else {
		tpl.Render(w, "products", products)
	}
}

func Product(w http.ResponseWriter, id string) {
	p, err := db.FindProduct(id)
	if err != nil {
		templates.Error(w, err)
		return
	}
	c, err := db.AllCategories()
	if err != nil {
		templates.Error(w, err)
		return
	}
	content := struct {
		Categories []product.Category
		Product    *product.Product
	}{
		c,
		p,
	}
	tpl.Render(w, "product", content)
}

func NewProduct(w http.ResponseWriter) {
	c, err := db.AllCategories()
	if err != nil {
		templates.Error(w, err)
		return
	}
	content := struct {
		Categories []product.Category
		Product    *product.Product
	}{
		c,
		&product.Product{},
	}
	tpl.Render(w, "product", content)
}

func CreateProduct(w http.ResponseWriter, r *http.Request, id string) {
	var err error
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		price = 0
	}
	p := &product.Product{
		ID:          r.FormValue("id"),
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       price,
		Active:      r.FormValue("active") != "",
		Category:    product.Category{ID: r.FormValue("category_id")},
	}
	if id == "" {
		err = db.CreateProduct(p)
	} else {
		err = db.UpdateProduct(p)
	}
	if err != nil {
		templates.Error(w, err)
	} else {
		http.Redirect(w, r, "/admin/products/", http.StatusFound)
	}
}
