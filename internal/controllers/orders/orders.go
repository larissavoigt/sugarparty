package orders

import (
	"net/http"

	"github.com/luizbranco/sugarparty/internal/models"
	"github.com/luizbranco/sugarparty/internal/views"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	c := models.NewCart(r)
	c.Ready = true
	content := struct {
		Cart  *models.Cart
		Error error
	}{
		c,
		nil,
	}
	switch r.Method {
	case "GET":
		id := r.URL.Path[len("/orders/"):]
		if id == "confirmation" {
			views.Render(w, "confirmation", nil)
		} else {
			views.Render(w, "order", content)
		}
	case "POST":
		o := &models.Order{
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Phone:   r.FormValue("phone"),
			Message: r.FormValue("message"),
		}
		err := models.CreateOrder(o, c)
		if err == nil {
			c.Destroy(w)
			http.Redirect(w, r, "/orders/confirmation", http.StatusFound)
		} else {
			content.Error = err
			views.Render(w, "order", content)
		}
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
