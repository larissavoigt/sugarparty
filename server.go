package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/luizbranco/sugarparty/internal/admin"
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
	auth.SetPassword(*password)
}

func main() {
	// server static assets files
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux := admin.NewServeMux()
	http.Handle("/admin/", mux)

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}

		id := r.URL.Path[len("/categories/"):]
		if id == "" {
			templates.NotFound(w)
		}

		c, err := db.FindCategory(id)
		if err != nil {
			templates.Error(w, err)
		}
		p, err := db.ActiveProducts(id)
		if err != nil {
			templates.Error(w, err)
		}
		content := struct {
			Category *product.Category
			Products []product.Product
		}{
			c,
			p,
		}
		templates.Render(w, "category", content)
	})

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			c := cart.New(r)
			templates.Render(w, "order", c)
		case "POST":
			http.Redirect(w, r, "/received", http.StatusFound)
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.URL.Path != "/" {
				templates.NotFound(w)
				return
			}
			categories, err := db.ActiveCategories()
			if err != nil {
				templates.Error(w, err)
			} else {
				templates.Render(w, "index", categories)
			}
		} else {
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
