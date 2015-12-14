package home

import (
	"net/http"

	"github.com/luizbranco/sugarparty/internal/models"
	"github.com/luizbranco/sugarparty/internal/views"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
}
