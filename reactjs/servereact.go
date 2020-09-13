package reactjs

import (
	"fmt"
	"html/template"
	"net/http"
)

const (
	Port = ":8081"
)

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../../reactjs/websocketjs.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

func serveClient(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("../../reactjs/client.html"))
	tmpl.Execute(w, nil)
}

func ServeHTML() {
	http.Handle("/", http.FileServer(http.Dir("../../reactjs/")))

	http.HandleFunc("/client", serveClient)
	http.HandleFunc("/ws", serveDashboard)
	http.ListenAndServe(Port, nil)
}
