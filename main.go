package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
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
var nbSockets = 0
var sockets [2]*websocket.Conn

func (messageType MessageType) String() string {
	return messageTypes[messageType - 1]
}

type Message struct {
	Type string
	//	gameMap     mapData
}
type MessageStartOfGame struct {
	Type         string
	PlayerNumber int
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	// enable cross origin connnections
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendMessage(t MessageType) (*Message) {
	return &Message{
		t.String(),
	}
}

func sendStartOfGame(number int) (*MessageStartOfGame) {
	return &MessageStartOfGame{
		"START_OF_GAME",
		number}
}

func initGameRoutine() {
	sockets[0].WriteJSON(sendStartOfGame(0))
	sockets[1].WriteJSON(sendStartOfGame(1))
}

func gameHasEnded() bool {
	return true
}

func playTurn(c *websocket.Conn) bool {
	c.WriteJSON(sendMessage(PLAY_TURN))
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		if !gameHasEnded() {
			return false
		} else {
			return true
		}
	}
	return false
}
func gameRoutine() {
	//	gameHasEnded := false
	for {
		playTurn(sockets[0])
		playTurn(sockets[1])
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
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
		c.WriteJSON(sendMessage(IDLE))
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

	//myMap := initMap()

	for x := 0; x < 19 * 19; x++ {
		//	fmt.Println("x:", x, " sum:", myMap[x].team) // Simple output.

	}

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
