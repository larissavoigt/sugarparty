package categories

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
}
