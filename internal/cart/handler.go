package cart

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	c := New(r)
	id := r.FormValue("id")
	c.Add(id, 1)
	c.Save(w)
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}
