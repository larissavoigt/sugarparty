package mail

import (
	"bytes"
	"errors"
	"net/smtp"
	"text/template"

	"github.com/larissavoigt/sugarparty/internal/models/order"
)

var (
	auth  smtp.Auth
	addr  string
	from  string
	to    string
	ready bool
	t     *template.Template
)

func init() {
	t = template.Must(template.ParseFiles("templates/mail.txt"))
}

type message struct {
	To      string
	Subject string
	Order   *order.Order
}

func (m message) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	err := t.Execute(buf, m)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func Config(recipient, username, password, host string) {
	addr = host + ":587"
	from = username
	to = recipient
	auth = smtp.PlainAuth("", username, password, host)
	ready = true
}

func NotifyOrder(id string) error {
	if ready {
		o, err := order.Find(id)
		if err != nil {
			return err
		}
		msg := message{
			To:      to,
			Subject: "Novo Pedido!",
			Order:   o,
		}

		return smtp.SendMail(addr, auth, from, []string{to}, msg.Bytes())
	} else {
		return errors.New("Email hasn't been configurated yet")
	}
}
