package reactjs

import (
	"fmt"
	"game/db"
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
	processClick(1)
}

func orangeHandler(w http.ResponseWriter, r *http.Request) {
	// orangeCounter++
	processClick(0)
}

func processClick(tp int) {
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
		latlong := fmt.Sprintf("%d %d", val.bc, val.oc)
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
