package admin

import (
	"net/http"

	"github.com/luizbranco/sugarparty/internal/auth"
	"github.com/luizbranco/sugarparty/internal/db"
	"github.com/luizbranco/sugarparty/internal/product"
	"github.com/luizbranco/sugarparty/internal/templates"
)

func categories(w http.ResponseWriter, r *http.Request) {
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
