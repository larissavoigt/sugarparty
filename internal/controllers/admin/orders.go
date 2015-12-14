package admin

import (
	"net/http"

	"github.com/luizbranco/sugarparty/internal/middlewares/auth"
	"github.com/luizbranco/sugarparty/internal/models"
	"github.com/luizbranco/sugarparty/internal/views"
)

func orders(w http.ResponseWriter, r *http.Request) {
	if !auth.Logged(r) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	switch r.Method {
	case "GET":
		orders, err := models.AllOrders()
		if err != nil {
			views.Error(w, err)
		} else {
			tpl.Render(w, "orders", orders)
		}
	case "POST":
		//TODO
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
