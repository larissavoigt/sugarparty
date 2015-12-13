package admin

import (
	"net/http"
	"strconv"

	"github.com/luizbranco/sugarparty/internal/auth"
	"github.com/luizbranco/sugarparty/internal/db"
	"github.com/luizbranco/sugarparty/internal/product"
	"github.com/luizbranco/sugarparty/internal/templates"
)

func products(w http.ResponseWriter, r *http.Request) {
	if !auth.Logged(r) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	id := r.URL.Path[len("/admin/products/"):]
	switch r.Method {
	case "GET":
		switch id {
		case "":
			indexProducts(w)
		case "new":
			newProduct(w)
		default:
			showProduct(w, id)
		}
	case "POST":
		createProduct(w, r, id)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func indexProducts(w http.ResponseWriter) {
	products, err := db.AllProducts()
	if err != nil {
		templates.Error(w, err)
	} else {
		tpl.Render(w, "products", products)
	}
}

func showProduct(w http.ResponseWriter, id string) {
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

func newProduct(w http.ResponseWriter) {
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

func createProduct(w http.ResponseWriter, r *http.Request, id string) {
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
