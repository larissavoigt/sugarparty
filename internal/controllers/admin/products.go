package admin

import (
	"net/http"
	"strconv"

	"github.com/larissavoigt/sugarparty/internal/middlewares/auth"
	"github.com/larissavoigt/sugarparty/internal/models/category"
	"github.com/larissavoigt/sugarparty/internal/models/product"
	"github.com/larissavoigt/sugarparty/internal/views"
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
			listProducts(w)
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

func listProducts(w http.ResponseWriter) {
	products, err := product.All()
	if err != nil {
		views.Error(w, err)
	} else {
		layout.Yield(w, "products", products)
	}
}

func showProduct(w http.ResponseWriter, id string) {
	p, err := product.Find(id)
	if err != nil {
		views.Error(w, err)
		return
	}
	c, err := category.All()
	if err != nil {
		views.Error(w, err)
		return
	}
	content := struct {
		Categories []category.Category
		Product    *product.Product
	}{
		c,
		p,
	}
	layout.Yield(w, "product", content)
}

func newProduct(w http.ResponseWriter) {
	c, err := category.All()
	if err != nil {
		views.Error(w, err)
		return
	}
	content := struct {
		Categories []category.Category
		Product    *product.Product
	}{
		c,
		&product.Product{},
	}
	layout.Yield(w, "product", content)
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
		Category:    category.Category{ID: r.FormValue("category_id")},
	}
	if id == "" {
		err = product.Create(p)
	} else {
		err = product.Update(p)
	}
	if err != nil {
		views.Error(w, err)
	} else {
		http.Redirect(w, r, "/admin/products/", http.StatusFound)
	}
}
