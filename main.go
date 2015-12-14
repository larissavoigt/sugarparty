package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/luizbranco/sugarparty/internal/controllers/admin"
	"github.com/luizbranco/sugarparty/internal/controllers/cart"
	"github.com/luizbranco/sugarparty/internal/middlewares/auth"
	"github.com/luizbranco/sugarparty/internal/models"
	"github.com/luizbranco/sugarparty/internal/views"
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
			views.NotFound(w)
			return
		}

		cat, err := models.FindCategory(id)
		if err != nil {
			views.Error(w, err)
		}
		p, err := models.ActiveProducts(id)
		if err != nil {
			views.Error(w, err)
		}
		content := struct {
			Category *models.Category
			Products []models.Product
			Cart     *models.Cart
		}{
			cat,
			p,
			models.NewCart(r),
		}
		views.Render(w, "category", content)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/" {
			views.NotFound(w)
			return
		}
		cat, err := models.ActiveCategories()
		if err != nil {
			views.Error(w, err)
			return
		}
		content := struct {
			Categories []models.Category
			Cart       *models.Cart
		}{
			cat,
			models.NewCart(r),
		}
		views.Render(w, "index", content)
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
