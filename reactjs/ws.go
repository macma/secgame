package reactjs

import (
	"fmt"
	"game/db"
	"html"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var blueCounter int
var orangeCounter int
var initTime time.Time

//https://blog.eleven-labs.com/en/construct-structure-go-graphql-api/
type ClickStruct struct {
	Type      int32 //0 for orange; 1 for blue
	ClickTime time.Time
}

type Clicks struct {
	bc int
	oc int
}

var clients = make(map[*websocket.Conn]bool)
var clickbroadcast = make(chan *Clicks)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/orange", orangeHandler).Methods("POST")
	router.HandleFunc("/blue", blueHandler).Methods("POST")
	router.HandleFunc("/ws", wsHandler)
	go echo()

	log.Fatal(http.ListenAndServe(":8844", router))
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

func clickWriter(coord *Clicks) {
	clickbroadcast <- coord
}

func blueHandler(w http.ResponseWriter, r *http.Request) {
	// blueCounter++
	processClick(1, w, r)
}

func orangeHandler(w http.ResponseWriter, r *http.Request) {
	// orangeCounter++
	processClick(0, w, r)
}

func processClick(tp int, w http.ResponseWriter, r *http.Request) {
	if orangeCounter == 0 && blueCounter == 0 {
		initTime = time.Now()
	}
	nowtime := time.Now()
	delta := nowtime.Sub(initTime)
	fmt.Println(fmt.Sprintf("the time difference is:    %v", delta))

	if delta <= 5*time.Second {
		switch tp {
		case 0:
			orangeCounter++
		default:
			blueCounter++
		}
		var clicks Clicks
		clicks.oc = orangeCounter
		clicks.bc = blueCounter
		db.InsertClicks(tp)
		go clickWriter(&clicks)
	} else {
		var tmpl = template.Must(template.ParseFiles("../../reactjs/htemplate.html"))

		var todo Todo
		todo.Name = html.UnescapeString(`var data = [
			{name: 1, uv: 4000, pv: 2400, amt: 2400},
			{name: 2, uv: 4000, pv: 2400, amt: 2500},
			{name: 3, uv: 4000, pv: 2400, amt: 2600},
			{name: 4, uv: 4000, pv: 2400, amt: 2700},
			{name: 5, uv: 4000, pv: 2400, amt: 2800},
	  ];
			  `)
		todo.uv = 4000
		todo.pv = 2400
		todo.amt = 2400
		// todo.Head = head
		tmpl.Execute(w, todo)

	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	blueCounter = 0
	orangeCounter = 0
	db.DropClicks()

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	clients[ws] = true
}

// 3
func echo() {
	for {
		val := <-clickbroadcast
		latlong := fmt.Sprintf("%d@%d@%s", val.bc, val.oc, `[
			{name: 'Page A', uv: 4000, pv: 2400, amt: 2400},
			{name: 'Page B', uv: 3000, pv: 1398, amt: 2210},
			{name: 'Page C', uv: 2000, pv: 9800, amt: 2290},
			{name: 'Page D', uv: 2780, pv: 3908, amt: 2000},
			{name: 'Page E', uv: 1890, pv: 4800, amt: 2181},
			{name: 'Page F', uv: 2390, pv: 3800, amt: 2500},
			{name: 'Page G', uv: 3490, pv: 4300, amt: 2100},
	  ]`)
		// send to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(latlong))
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}

	}
}
