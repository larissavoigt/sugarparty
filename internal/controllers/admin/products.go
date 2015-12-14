package admin

import (
	"net/http"
	"strconv"

	"github.com/luizbranco/sugarparty/internal/middlewares/auth"
	"github.com/luizbranco/sugarparty/internal/models"
	"github.com/luizbranco/sugarparty/internal/views"
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
	products, err := models.AllProducts()
	if err != nil {
		views.Error(w, err)
	} else {
		tpl.Render(w, "products", products)
	}
}

func showProduct(w http.ResponseWriter, id string) {
	p, err := models.FindProduct(id)
	if err != nil {
		views.Error(w, err)
		return
	}
	c, err := models.AllCategories()
	if err != nil {
		views.Error(w, err)
		return
	}
	content := struct {
		Categories []models.Category
		Product    *models.Product
	}{
		c,
		p,
	}
	tpl.Render(w, "product", content)
}

func newProduct(w http.ResponseWriter) {
	c, err := models.AllCategories()
	if err != nil {
		views.Error(w, err)
		return
	}
	content := struct {
		Categories []models.Category
		Product    *models.Product
	}{
		c,
		&models.Product{},
	}
	tpl.Render(w, "product", content)
}

func createProduct(w http.ResponseWriter, r *http.Request, id string) {
	var err error
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		price = 0
	}
	p := &models.Product{
		ID:          r.FormValue("id"),
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       price,
		Active:      r.FormValue("active") != "",
		Category:    models.Category{ID: r.FormValue("category_id")},
	}
	if id == "" {
		err = models.CreateProduct(p)
	} else {
		err = models.UpdateProduct(p)
	}
	if err != nil {
		views.Error(w, err)
	} else {
		http.Redirect(w, r, "/admin/products/", http.StatusFound)
	}
}
