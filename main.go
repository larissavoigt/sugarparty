package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/luizbranco/sugarparty/internal/controllers/admin"
	"github.com/luizbranco/sugarparty/internal/controllers/cart"
	"github.com/luizbranco/sugarparty/internal/controllers/categories"
	"github.com/luizbranco/sugarparty/internal/controllers/home"
	"github.com/luizbranco/sugarparty/internal/controllers/orders"
	"github.com/luizbranco/sugarparty/internal/mail"
	"github.com/luizbranco/sugarparty/internal/middlewares/auth"
)

var (
	port     = flag.String("port", "3000", "Server port")
	password = flag.String("admin-password", "", "Admin password")
	mailTo   = flag.String("mail-recipient", "", "Mail recipient")
	mailUser = flag.String("mail-username", "", "Mail username")
	mailPass = flag.String("mail-password", "", "Mail password")
	mailHost = flag.String("mail-host", "", "Mail host")
)

func init() {
	flag.Parse()
	auth.SetPassword(*password)
	mail.Config(*mailTo, *mailUser, *mailPass, *mailHost)
}

func main() {
	// server static assets files
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.Handle("/admin/", admin.NewServeMux())
	http.HandleFunc("/cart/", cart.Handler)
	http.HandleFunc("/categories/", categories.Handler)
	http.HandleFunc("/orders/", orders.Handler)
	http.HandleFunc("/", home.Handler)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
