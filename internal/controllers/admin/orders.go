package admin

import (
	"net/http"
	"strconv"

	"github.com/larissavoigt/sugarparty/internal/middlewares/auth"
	"github.com/larissavoigt/sugarparty/internal/models/order"
	"github.com/larissavoigt/sugarparty/internal/views"
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
		if id == "" {
			http.Error(w, "", http.StatusBadRequest)
		} else {
			updateOrder(w, r, id)
		}
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func listOrders(w http.ResponseWriter) {
	orders, err := order.All()
	content := struct {
		Orders []order.Order
		Page   string
	}{
		orders,
		"orders",
	}
	if err != nil {
		views.Error(w, err)
	} else {
		layout.Yield(w, "orders", content)
	}
}

func showOrder(w http.ResponseWriter, id string) {
	o, err := order.Find(id)
	content := struct {
		Order       *order.Order
		StatusNames []string
		Page        string
	}{
		o,
		order.StatusNames,
		"orders",
	}
	if err != nil {
		views.Error(w, err)
	} else {
		layout.Yield(w, "order", content)
	}
}

func updateOrder(w http.ResponseWriter, r *http.Request, id string) {
	status := r.FormValue("status")
	i, _ := strconv.Atoi(status)
	err := order.Update(id, i)
	if err != nil {
		views.Error(w, err)
	} else {
		http.Redirect(w, r, "/admin/orders/", http.StatusFound)
	}
}
