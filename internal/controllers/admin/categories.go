package admin

import (
	"net/http"

	"github.com/larissavoigt/sugarparty/internal/middlewares/auth"
	"github.com/larissavoigt/sugarparty/internal/models/category"
	"github.com/larissavoigt/sugarparty/internal/views"
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
			categories, err := category.All()
			if err != nil {
				views.Error(w, err)
			} else {
				layout.Yield(w, "categories", categories)
			}
		case "new":
			layout.Yield(w, "category", nil)
		default:
			c, err := category.Find(id)
			if err != nil {
				views.Error(w, err)
			} else {
				tpl.Render(w, "category", c)
			}
		}
	case "POST":
		var err error
		c := &category.Category{
			ID:          r.FormValue("id"),
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
		}
		if id == "" {
			err = category.Create(c)
		} else {
			err = category.Update(c)
		}
		if err != nil {
			views.Error(w, err)
		} else {
			http.Redirect(w, r, "/admin/categories", http.StatusFound)
		}
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
