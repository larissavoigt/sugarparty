package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/luizbranco/sugarparty/internal/auth"
	"github.com/luizbranco/sugarparty/internal/cart"
	"github.com/luizbranco/sugarparty/internal/db"
	"github.com/luizbranco/sugarparty/internal/product"
	"github.com/luizbranco/sugarparty/internal/templates"
)

var (
	port     = flag.String("port", "3000", "Server port")
	password = flag.String("password", "", "Admin password")
)

func init() {
	flag.Parse()
}

func main() {
	tpl := templates.New("templates")
	atpl := templates.New("templates/admin")

	// server static assets files
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if auth.Admin(r, *password) {
			atpl.Render(w, "index", nil)
		} else {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		}
	})

	http.HandleFunc("/admin/categories/", func(w http.ResponseWriter, r *http.Request) {
		if !auth.Admin(r, *password) {
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
					tpl.Error(w, err)
				} else {
					atpl.Render(w, "categories", categories)
				}
			case "new":
				atpl.Render(w, "category", nil)
			default:
				c, err := db.FindCategory(id)
				if err != nil {
					tpl.Error(w, err)
				} else {
					atpl.Render(w, "category", c)
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
				tpl.Error(w, err)
			} else {
				http.Redirect(w, r, "/admin/categories", http.StatusFound)
			}
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/admin/products/", func(w http.ResponseWriter, r *http.Request) {
		if !auth.Admin(r, *password) {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}

		id := r.URL.Path[len("/admin/products/"):]
		switch r.Method {
		case "GET":
			switch id {
			case "":
				products, err := db.AllProducts()
				if err != nil {
					tpl.Error(w, err)
				} else {
					atpl.Render(w, "products", products)
				}
			case "new":
				c, err := db.AllCategories()
				if err != nil {
					tpl.Error(w, err)
					return
				}
				content := struct {
					Categories []product.Category
					Product    *product.Product
				}{
					c,
					&product.Product{},
				}
				atpl.Render(w, "product", content)
			default:
				p, err := db.FindProduct(id)
				if err != nil {
					tpl.Error(w, err)
					return
				}
				c, err := db.AllCategories()
				if err != nil {
					tpl.Error(w, err)
					return
				}
				content := struct {
					Categories []product.Category
					Product    *product.Product
				}{
					c,
					p,
				}
				content.Product = p
				atpl.Render(w, "product", content)
			}
		case "POST":
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
				tpl.Error(w, err)
			} else {
				http.Redirect(w, r, "/admin/products", http.StatusFound)
			}
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/admin/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			atpl.Render(w, "login", nil)
		case "POST":
			auth.SaveSession(w, r.FormValue("password"))
			http.Redirect(w, r, "/admin", http.StatusFound)
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/admin/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.DestroySession(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			c := cart.New(r)
			log.Println(c)
			tpl.Render(w, "order", c)
		case "POST":
			http.Redirect(w, r, "/received", http.StatusFound)
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.URL.Path != "/" {
				tpl.NotFound(w)
				return
			}
			categories, err := db.ActiveCategories()
			if err != nil {
				tpl.Error(w, err)
			} else {
				tpl.Render(w, "index", categories)
			}
		} else {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
