package categories

import (
	"net/http"

	"github.com/luizbranco/sugarparty/internal/models/cart"
	"github.com/luizbranco/sugarparty/internal/models/category"
	"github.com/luizbranco/sugarparty/internal/models/product"
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

	cat, err := category.Find(id)
	if err != nil {
		views.Error(w, err)
	}
	p, err := product.Active(id)
	if err != nil {
		views.Error(w, err)
	}
	content := struct {
		Category *category.Category
		Products []product.Product
		Cart     *cart.Cart
	}{
		cat,
		p,
		cart.New(r),
	}
	views.Render(w, "category", content)
}
