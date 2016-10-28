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

/*var upgrader = websocket.Upgrader{
	// enable cross origin connnections
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}*/

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
		log.Println("New turn:")
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
func checkCase(myMap []protocol.MapData, pos int, team int) (bool) {
	if myMap[pos].Team == (team + 1) % 2 {
		return (true)
	}
	return (false)
}

func checkPair(myMap []protocol.MapData, pos int, team int) ([]protocol.MapData, int) {
	var emptyData protocol.MapData
	emptyData.Empty = true
	emptyData.Playable = true
	emptyData.Team = -1
	captured := 0
	if (pos - (19 * 3)) >= 0 { // NORD
		if checkCase(myMap, pos - (19 * 1), team) && checkCase(myMap, pos - (19 * 2), team) && checkCase(myMap, pos - (19 * 3), (team + 1) % 2 ) {
			myMap[pos - (19 * 1)] =  emptyData
			myMap[pos - (19 * 2)] =  emptyData
			captured += 2
		}
	}
	if (pos - (19 * 3) + 3) >= 0 && pos % 19 <= 15 { // NORD EST
		if checkCase(myMap, pos - (19 * 1) + 1, team) && checkCase(myMap, pos - (19 * 2) + 2, team) && checkCase(myMap, pos - (19 * 3) + 3, (team + 1) % 2 ) {
			myMap[pos - (19 * 1) + 1] =  emptyData
			myMap[pos - (19 * 2) + 2] =  emptyData
			captured += 2
		}
	}
	if (pos + 3) < 19 * 19 && pos % 19 <= 15 { // EST
		if checkCase(myMap, pos + 1, team) && checkCase(myMap, pos + 2, team) && checkCase(myMap, pos + 3, (team + 1) % 2 ) {
			myMap[pos + 1] =  emptyData
			myMap[pos + 2] =  emptyData
			captured += 2
		}
	}
	if (pos + (19 * 3) + 3) < 19 * 19 && pos % 19 <= 15 { // SUD EST
		if checkCase(myMap, pos + (19 * 1) + 1, team) && checkCase(myMap, pos + (19 * 2) + 2, team) && checkCase(myMap, pos + (19 * 3) + 3, (team + 1) % 2 ) {
			myMap[pos + (19 * 1) + 1] =  emptyData
			myMap[pos + (19 * 2) + 2] =  emptyData
			captured += 2
		}
	}
	if (pos + (19 * 3)) < 19 * 19 { // SUD
		if checkCase(myMap, pos + (19 * 1), team) && checkCase(myMap, pos + (19 * 2), team) && checkCase(myMap, pos + (19 * 3), (team + 1) % 2 ) {
			myMap[pos + (19 * 1)] =  emptyData
			myMap[pos + (19 * 2)] =  emptyData
			captured += 2
		}
	}
	if (pos + (19 * 3) - 3) < 19 * 19 && pos % 19 >= 3 { // SUD OUEST
		if checkCase(myMap, pos + (19 * 1) - 1, team) && checkCase(myMap, pos + (19 * 2) - 2, team) && checkCase(myMap, pos + (19 * 3) - 3, (team + 1) % 2 ) {
			myMap[pos + (19 * 1) - 1] =  emptyData
			myMap[pos + (19 * 2) - 2] =  emptyData
			captured += 2
		}
	}
	if (pos - 3) >= 0 && pos % 19 >= 3 { // OUEST
		if checkCase(myMap, pos - 1, team) && checkCase(myMap, pos - 2, team) && checkCase(myMap, pos - 3, (team + 1) % 2 ) {
			myMap[pos - 1] =  emptyData
			myMap[pos - 2] =  emptyData
			captured += 2
		}
	}
	if (pos - (19 * 3) - 3) >= 0 && pos % 19 >= 3 { // NORD OUEST
		if checkCase(myMap, pos - (19 * 1) - 1, team) && checkCase(myMap, pos - (19 * 2) - 2, team) && checkCase(myMap, pos - (19 * 3) - 3, (team + 1) % 2 ) {
			myMap[pos - (19 * 1) - 1] =  emptyData
			myMap[pos - (19 * 2) - 2] =  emptyData
			captured += 2
		}
	}
	return myMap, captured
}

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
	//http.HandleFunc("/ws", wsHandler)

	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
