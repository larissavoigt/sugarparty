package admin

import (
	"net/http"

	"github.com/luizbranco/sugarparty/internal/auth"
	"github.com/luizbranco/sugarparty/internal/db"
	"github.com/luizbranco/sugarparty/internal/product"
	"github.com/luizbranco/sugarparty/internal/templates"
)

var tpl = templates.New("templates/admin")

func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin/categories/", Categories)
	mux.HandleFunc("/admin/products/", Products)
	mux.HandleFunc("/admin/login", Login)
	mux.HandleFunc("/admin/logout", Logout)
	mux.HandleFunc("/admin/", Default)
	return mux
}

func Categories(w http.ResponseWriter, r *http.Request) {
	if !auth.Logged(r) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}
	id := r.URL.Path[len("/admin/categories/"):]
	switch r.Method {
	case "GET":
		switch id {
		case "":
			categories, err := db.AllCategories()
			if err != nil {
				templates.Error(w, err)
			} else {
				tpl.Render(w, "categories", categories)
			}
		case "new":
			tpl.Render(w, "category", nil)
		default:
			c, err := db.FindCategory(id)
			if err != nil {
				templates.Error(w, err)
			} else {
				tpl.Render(w, "category", c)
			}
		}
	case "POST":
		var err error
		c := &product.Category{
			ID:          r.FormValue("id"),
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
		}
		if id == "" {
			err = db.CreateCategory(c)
		} else {
			err = db.UpdateCategory(c)
		}
		if err != nil {
			templates.Error(w, err)
		} else {
			http.Redirect(w, r, "/admin/categories", http.StatusFound)
		}
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func Products(w http.ResponseWriter, r *http.Request) {
	if !auth.Logged(r) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	id := r.URL.Path[len("/admin/products/"):]
	switch r.Method {
	case "GET":
		switch id {
		case "":
			AllProducts(w)
		case "new":
			NewProduct(w)
		default:
			Product(w, id)
		}
	case "POST":
		CreateProduct(w, r, id)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tpl.Render(w, "login", nil)
	case "POST":
		auth.SaveSession(w, r.FormValue("password"))
		http.Redirect(w, r, "/admin/", http.StatusFound)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.DestroySession(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func Default(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/admin/" {
		if auth.Logged(r) {
			tpl.Render(w, "index", nil)
		} else {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		}
	} else {
		templates.NotFound(w)
	}
}
