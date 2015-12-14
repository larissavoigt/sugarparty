package admin

import (
	"net/http"

	"github.com/luizbranco/sugarparty/internal/middlewares/auth"
	"github.com/luizbranco/sugarparty/internal/models"
	"github.com/luizbranco/sugarparty/internal/views"
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
			categories, err := models.AllCategories()
			if err != nil {
				views.Error(w, err)
			} else {
				tpl.Render(w, "categories", categories)
			}
		case "new":
			tpl.Render(w, "category", nil)
		default:
			c, err := models.FindCategory(id)
			if err != nil {
				views.Error(w, err)
			} else {
				tpl.Render(w, "category", c)
			}
		}
	case "POST":
		var err error
		c := &models.Category{
			ID:          r.FormValue("id"),
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
		}
		if id == "" {
			err = models.CreateCategory(c)
		} else {
			err = models.UpdateCategory(c)
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
