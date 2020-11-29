package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

//encapsulates all methods and states needed to parse the html
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//attach a pointer method which satisfies HandleFunc signature (w, r)
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	//write to response
	t.templ.Execute(w, nil)
}

func main() {

	r := newRoom()
	//pass in a pointer since the pointer implements the Handler interface, which has only the ServeHTTP method! http.HandleFunc doesn't work because it requires a function with the same signature, whicle the method pointer has 3 args (the type itself)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	//starts the room in a separate thread
	go r.run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
