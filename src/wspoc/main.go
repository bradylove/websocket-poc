package main

import (
	"crypto/rand"
	"encoding/base64"
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
	sessions   map[string]*session
)

type Tick struct {
	Count int `json:"count"`
}

type session struct {
	Id    string `json:"id"`
	index int
}

func startServer(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/stream/", wsHandler).Methods("GET")
	r.HandleFunc("/stream/{sessionId}", wsHandler).Methods("GET")

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

	vars := mux.Vars(r)
	sessionId := vars["sessionId"]
	isNewStream := sessionId == ""

	var sess *session
	if isNewStream {
		sess = &session{Id: generateId()}
		sessions[sess.Id] = sess
		conn.WriteJSON(sess)
		fmt.Println("New websocket connection opened with id", sess.Id)
	}

	if !isNewStream {
		var ok bool
		sess, ok = sessions[sessionId]
		if !ok {
			fmt.Println("Invalid session id:", sessionId)
			defer conn.Close()
			return
		}
		fmt.Println("Resuming websocket stream with id:", sess.Id)
	}

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		for _ = range ticker.C {
			if sess.index == 50 && isNewStream {
				conn.Close()
				return
			}

			count := tickBuffer[sess.index]
			err := conn.WriteJSON(Tick{Count: count})
			if err != nil {
				fmt.Println(err)
				ticker.Stop()
				return
			}

			fmt.Println("Wrote json message with count:", count)

			sess.index++
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

func generateId() string {
	data := make([]byte, 12)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(data)
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

	sessions = make(map[string]*session)

	fmt.Println("Starting server on port", port)
	startServer(port)
}
