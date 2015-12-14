package cart

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	c := New(r)
	id := r.FormValue("id")
	url := r.URL.Path[len("/cart/"):]
	switch url {
	case "":
		c.Add(id, 1)
	case "decrease":
		c.Add(id, -1)
	case "remove":
		c.Remove(id)
	}
	c.Save(w)
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}
