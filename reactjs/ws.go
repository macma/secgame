package reactjs

import (
	"encoding/json"
	"fmt"
	"game/db"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
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
	router.HandleFunc("/query", queryHandler).Methods("POST")
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
	processClick(1, w, r)
}

type ClickGraph struct {
	id     string
	orange int
	blue   int
	black  int
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	s, _ := ioutil.ReadAll(r.Body)
	strQuery := string(s)
	strQuery = strQuery[1 : len(strQuery)-1]
	fmt.Println(fmt.Sprintf("the body is %s", strQuery))

	clicks := db.SelectClicks()

	cgs := make([]ClickGraph, 11)
	for i := range cgs {
		cgs[i].id = fmt.Sprintf("%.1f", float64(i)/2)
	}
	for i, a := range clicks {
		if i != 0 {
			delta := a.Created_at.Sub(clicks[0].Created_at)
			sec := int64((delta/time.Millisecond)/500) + 1
			//fmt.Println(int64(delta / time.Millisecond))
			switch a.Color {
			case 0:
				cgs[sec].orange++
			default:
				cgs[sec].blue++
			}
		} else {
			switch a.Color {
			case 0:
				cgs[1].orange++
			default:
				cgs[1].blue++
			}
		}
	}
	for i, a := range cgs {
		cgs[i].black = a.blue - a.orange
	}

	gqobj := getGraphqlObject(cgs, strQuery)
	bjson, _ := json.Marshal(gqobj)
	sjson := string(bjson)
	sjson = sjson[10 : len(sjson)-1]
	fmt.Println("the value get from db is:" + sjson)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(sjson))

}

func getGraphqlObject(cgs []ClickGraph, strQuery string) interface{} {
	var ClickObject = graphql.NewObject(graphql.ObjectConfig{
		Name: "Click",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"orange": &graphql.Field{
				Type: graphql.String,
			},
			"blue": &graphql.Field{
				Type: graphql.String,
			},
			"black": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	var RootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"clicks": &graphql.Field{
				Type: graphql.NewList(ClickObject),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					clicks := []map[string]string{}
					for _, a := range cgs {
						obj := map[string]string{
							"id":     a.id,
							"black":  fmt.Sprintf("%d", a.black),
							"orange": fmt.Sprintf("%d", a.orange),
							"blue":   fmt.Sprintf("%d", a.blue),
						}
						clicks = append(clicks, obj)
					}
					return clicks, nil
				},
			},
		},
	})
	schemaConfig := graphql.SchemaConfig{Query: RootQuery}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	// query := `
	// 	{
	// 		clicks
	// 		{
	// 			id
	// 			black
	// 			orange
	// 			blue
	// 		}
	// 	}
	// `
	params := graphql.Params{Schema: schema, RequestString: strQuery}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	// rJSON, _ := json.Marshal(r.Data)
	return r.Data
	// fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}

}
func orangeHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("success"))
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
				{"name": "0", "uv": 4000, "pv": 2400, "amt": 10},
				{"name": "0.5", "uv": 3000, "pv": 1398, "amt": 20},
				{"name": "1", "uv": -2000, "pv": 9800, "amt": 30},
				{"name": "1.5", "uv": 2780, "pv": 3908, "amt": 2000},
				{"name": "2", "uv": 1890, "pv": 4800, "amt": -2181},
				{"name": "2.5", "uv": 2390, "pv": 3800, "amt": 2500},
				{"name": "3", "uv": 3490, "pv": 4300, "amt": 2100},
				{"name": "3.5", "uv": 2780, "pv": 3908, "amt": 2000},
				{"name": "4", "uv": 1890, "pv": 4800, "amt": -2181},
				{"name": "4.5", "uv": 2390, "pv": 3800, "amt": 2500},
				{"name": "5", "uv": 3490, "pv": 4300, "amt": 2100}]`)
		// latlong := fmt.Sprintf("%d@%d@%s", val.bc, val.oc, `{"data":{"hello":"world"}}`)
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
