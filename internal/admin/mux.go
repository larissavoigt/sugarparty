package admin

import (
	"net/http"

	"github.com/luizbranco/sugarparty/internal/auth"
	"github.com/luizbranco/sugarparty/internal/templates"
)

var tpl = templates.New("templates/admin")

func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin/categories/", categories)
	mux.HandleFunc("/admin/products/", products)
	mux.HandleFunc("/admin/login", login)
	mux.HandleFunc("/admin/logout", logout)
	mux.HandleFunc("/admin/", index)
	return mux
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tpl.Render(w, "login", nil)
	case "POST":
		if auth.Login(w, r.FormValue("password")) {
			http.Redirect(w, r, "/admin/", http.StatusFound)
		} else {
			tpl.Render(w, "login", "Invalid password")
		}
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/admin/" {
		if auth.Logged(r) {
			tpl.Render(w, "index", nil)
		} else {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		}
	} else {
		templates.NotFound(w)
	}
}
