package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	tickBuffer []int
)

type Tick struct {
	Count int
}

func startServer(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/stream/{streamId}", wsHandler)

	http.Handle("/", r)
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade:", err)
		return
	}
	defer conn.Close()

	fmt.Println("New websocket connection opened")

	var counter int
	ticker := time.NewTicker(200 * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			err := conn.WriteJSON(Tick{Count: counter})
			if err != nil {
				fmt.Println(err)
				ticker.Stop()
				return
			}

			fmt.Println("Wrote json message with count:", counter)

			counter++
		}
	}()

	for {
		var m map[string]interface{}
		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("ERR:", err)
			return
		}

		fmt.Println(m)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9292"
	}

	tickBuffer = make([]int, 10000)
	for i := 0; i < cap(tickBuffer); i++ {
		tickBuffer[i] = i
	}

	fmt.Println("Starting server on port", port)
	startServer(port)
}
