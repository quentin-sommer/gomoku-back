package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"encoding/json"
)

type mapData struct {
	empty, playable bool
	team            int
}
type MessageType int

const (
	IDLE = iota + 1
	START_OF_GAME
	PLAY_TURN
	END_OF_GAME
)

var messageTypes = [...]string{
	"IDLE",
	"START_OF_GAME",
	"PLAY_TURN",
	"END_OF_GAME",
}

func (messageType MessageType) String() string {
	return messageTypes[messageType - 1]
}

type Message struct {
	Type string
	//	gameMap     mapData
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	// enable cross origin connnections
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendTestMessage(t MessageType) ([]byte) {
	msg := &Message{
		t.String(),
	}

	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("sendTestMessage:", err)
		return nil
	}
	return b
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		msg := sendTestMessage(START_OF_GAME)
		err = c.WriteMessage(mt, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func initMap() ([]mapData) {

	myMap := make([]mapData, 19 * 19)
	for x := 0; x < 19 * 19; x++ {
		myMap[x].empty = true
		myMap[x].playable = true
		myMap[x].team = -1
	}
	return myMap
}

func main() {

	myMap := initMap()

	for  x := 0; x < 19*19; x++ {
		fmt.Println("x:", x, " sum:", myMap[x].team) // Simple output.

	}

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
