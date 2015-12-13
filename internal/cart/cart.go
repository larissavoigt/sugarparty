package cart

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

type Cart struct {
	items map[string]int
}

func New(r *http.Request) *Cart {
	c := &Cart{items: make(map[string]int)}
	cookie, err := r.Cookie("cart")
	if err != nil {
		return c
	}
	items := strings.Split(cookie.Value, " ")
	if len(items)%2 != 0 {
		return c
	}

	for i := 0; i < len(items); i += 2 {
		n, err := strconv.Atoi(items[i+1])
		if err == nil {
			c.Add(items[i], n)
		}
	}

	return c
}

func (c *Cart) Add(id string, qty int) {
	c.items[id] = qty
}

func (c *Cart) Remove(id string) {
	delete(c.items, id)
}

func (c *Cart) MarshalText() (text []byte, err error) {
	var buf bytes.Buffer
	for k, v := range c.items {
		_, err = buf.WriteString(k + " " + strconv.Itoa(v) + " ")
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (c *Cart) Save(w http.ResponseWriter) {
	val, _ := c.MarshalText()
	cookie := &http.Cookie{
		Name:     "cart",
		Value:    string(val),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (c *Cart) Total() error {
	return nil
}
