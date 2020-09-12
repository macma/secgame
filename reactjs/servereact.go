package reactjs

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
)

const (
	Port = ":8081"
)

func serveStaticClient(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../../reactjs/websocketjs.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

type Todo struct {
	Name string
	uv   int
	pv   int
	amt  int
}

func htemplateHandler(w http.ResponseWriter, r *http.Request) {

	var tmpl = template.Must(template.ParseFiles("../../reactjs/htemplate.html"))

	var todo Todo
	todo.Name = html.UnescapeString(`var data = [
  ];
		  `)
	// todo.uv = 4000
	// todo.pv = 2400
	// todo.amt = 2400
	// todo.Head = head
	tmpl.Execute(w, todo)

}

func ServeHTML() {
	http.Handle("/", http.FileServer(http.Dir("../../reactjs/")))

	http.HandleFunc("/htemplate", htemplateHandler)
	http.HandleFunc("/ws", serveStaticClient)
	http.ListenAndServe(Port, nil)
}
