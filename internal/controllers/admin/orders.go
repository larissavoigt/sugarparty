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

	id := r.URL.Path[len("/admin/orders/"):]
	switch r.Method {
	case "GET":
		if id == "" {
			listOrders(w)
		} else {
			showOrder(w, id)
		}
	case "POST":
		//TODO
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func listOrders(w http.ResponseWriter) {
	orders, err := models.AllOrders()
	if err != nil {
		views.Error(w, err)
	} else {
		tpl.Render(w, "orders", orders)
	}
}

func showOrder(w http.ResponseWriter, id string) {
	o, err := models.FindOrder(id)
	if err != nil {
		views.Error(w, err)
	} else {
		tpl.Render(w, "order", o)
	}
}
