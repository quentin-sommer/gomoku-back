package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"./protocol"
	"encoding/json"
)

var nbSockets = 0
var sockets [2]*websocket.Conn

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	// enable cross origin connnections
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func initGameRoutine() {
	sockets[0].WriteJSON(protocol.SendStartOfGame(0))
	sockets[1].WriteJSON(protocol.SendStartOfGame(1))
}

func playTurn(c *websocket.Conn, msg []byte) ([]byte) {
	c.WriteMessage(1, msg)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		//log.Printf("recv: %s", message)
		// TODO : parse message and check if move is valid
		return message
	}
	return nil
}

func gameRoutine() {
	message, _ := json.Marshal(protocol.SendPlayTurn(initMap()))
	for {
		message = playTurn(sockets[0], message)
		message = playTurn(sockets[1], message)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	switch nbSockets {
	// first connected
	case 0:
		sockets[0] = c
		nbSockets += 1
		log.Printf("%d socket connected\n", nbSockets)
		c.WriteJSON(protocol.SendIdle())
		break
	// all sockets connected
	case 1:
		sockets[1] = c
		nbSockets += 1
		log.Print("starting game")
		initGameRoutine()
		gameRoutine()
		break
	default:
		break
	}
}

// règles bien expliqué http://maximegirou.com/files/projets/b1/gomoku.pdf

// function qui check dans un sens choisi (N NE E SE S SW W NW) pour vérifier la fin du jeu
// attention un gars peut casser une ligne de 5 pions avec une paire

// function qui check la regle "LE DOUBLE-TROIS"

// function qui check s'il peut NIQUER une paire et s'il peut tej les deux entre (prendre plusieurs pair d'un coup)


func initMap() ([]protocol.MapData) {
	myMap := make([]protocol.MapData, 19 * 19)
	for x := 0; x < 19 * 19; x++ {
		myMap[x].Empty = true
		myMap[x].Playable = true
		myMap[x].Team = -1
	}
	return myMap
}

func main() {

	//myMap := initMap()

	for x := 0; x < 19 * 19; x++ {
		//	fmt.Println("x:", x, " sum:", myMap[x].team) // Simple output.

	}

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
