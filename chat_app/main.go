package main

import (
	"chat_app/trace"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
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
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the app")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	//pass in a pointer since the pointer implements the Handler interface, which has only the ServeHTTP method! http.HandleFunc doesn't work because it requires a function with the same signature, whicle the method pointer has 3 args (the type itself)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	//starts the room in a separate thread, otherwise either the room or the server will block!
	go r.run()

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

//Server Thread
//Room Thread
//New Client comes in: 1 write to socket thread, 1 read from socket thread
