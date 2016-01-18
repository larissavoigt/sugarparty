package home

import (
	"net/http"

	"github.com/larissavoigt/sugarparty/internal/models/cart"
	"github.com/larissavoigt/sugarparty/internal/models/category"
	"github.com/larissavoigt/sugarparty/internal/views"
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
	cat, err := category.Active()
	if err != nil {
		views.Error(w, err)
		return
	}
	content := struct {
		Categories []category.Category
		Cart       *cart.Cart
	}{
		cat,
		cart.New(r),
	}
	views.Render(w, "index", content)
}
