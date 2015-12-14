package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/luizbranco/sugarparty/internal/admin"
	"github.com/luizbranco/sugarparty/internal/auth"
	"github.com/luizbranco/sugarparty/internal/cart"
	"github.com/luizbranco/sugarparty/internal/models"
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

	http.Handle("/admin/", admin.NewServeMux())
	http.HandleFunc("/cart/", cart.Handler)

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Path[len("/categories/"):]
		if id == "" {
			templates.NotFound(w)
			return
		}

		cat, err := models.FindCategory(id)
		if err != nil {
			templates.Error(w, err)
		}
		p, err := models.ActiveProducts(id)
		if err != nil {
			templates.Error(w, err)
		}
		content := struct {
			Category *models.Category
			Products []models.Product
			Cart     *cart.Cart
		}{
			cat,
			p,
			cart.New(r),
		}
		templates.Render(w, "category", content)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/" {
			templates.NotFound(w)
			return
		}
		cat, err := models.ActiveCategories()
		if err != nil {
			templates.Error(w, err)
			return
		}
		content := struct {
			Categories []models.Category
			Cart       *cart.Cart
		}{
			cat,
			cart.New(r),
		}
		templates.Render(w, "index", content)
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
