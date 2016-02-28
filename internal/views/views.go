package views

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

type Views struct {
	*template.Template
	path string
}

var tpl = New("templates")

func New(path string) *Views {
	t := template.New(path)
	template.Must(t.ParseGlob(path + "/*.html"))
	return &Views{t, path}
}

func Layout(p, layout string) *Views {
	t, err := template.New(layout + ".html").ParseFiles(path.Join(p, layout+".html"))
	if err != nil {
		log.Fatal(err)
	}
	return &Views{t, p}
}

func (v *Views) Render(w http.ResponseWriter, name string, data interface{}) {
	err := v.ExecuteTemplate(w, name+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *Views) Yield(w http.ResponseWriter, name string, data interface{}) {
	t, err := template.Must(v.Clone()).ParseFiles(path.Join(v.path, name+".html"))
	if err == nil {
		err = t.Execute(w, data)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Render(w http.ResponseWriter, name string, data interface{}) {
	tpl.Render(w, name, data)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	tpl.Render(w, "404", nil)
}

func Error(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	tpl.Render(w, "500", err)
}
