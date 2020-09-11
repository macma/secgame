package reactjs

import (
	"fmt"
	"html/template"
	"net/http"
)

const (
	Port = ":8080"
)

func serveStaticClient(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../../reactjs/websocketjs.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

func ServeHTML() {
	http.HandleFunc("/ws", serveStaticClient)
	http.ListenAndServe(Port, nil)
}
