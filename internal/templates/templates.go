package templates

import (
	"html/template"
	"net/http"
)

type Templates struct {
	*template.Template
}

func New(path string) *Templates {
	t := template.New(path)
	template.Must(t.ParseGlob(path + "/*.html"))
	return &Templates{t}
}

func (t *Templates) Render(w http.ResponseWriter, name string, data interface{}) {
	err := t.ExecuteTemplate(w, name+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Templates) NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	t.Render(w, "404", nil)
}

func (t *Templates) Error(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	t.Render(w, "500", err)
}
