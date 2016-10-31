package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"./protocol"
	//"fmt"
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

// règles bien expliqué http://maximegirou.com/files/projets/b1/gomoku.pdf

// function qui check dans un sens choisi (N NE E SE S SW W NW) pour vérifier la fin du jeu
// attention un gars peut casser une ligne de 5 pions avec une paire

func checkLigne(myMap []protocol.MapData, pos int, team int, val int, add int) (int) {
	if (add == -18 && pos % 19 <= 15) || (add == 18 && pos % 19 >= 3) || (add == -20 && pos % 19 >= 3) || (add == 20 && pos % 19 <= 15) {
		if pos < 19 * 19 && pos >= 0 && myMap[pos].Player != team {
			return val
		}
	}
	return checkLigne(myMap, pos + add, team, val + 1, add)
}

func checkEnd(myMap []protocol.MapData, pos int, team int) (bool) {
	var nb int
	// horizontal
	nb = checkLigne(myMap, pos, team, 0, 1)
	nb += checkLigne(myMap, pos, team, 0, -1)
	if nb - 1 == 5 {
		//fmt.Printf("END 5 IN A ROW\n")
		return true
	}

	// vertical
	nb = checkLigne(myMap, pos, team, 0, 19)
	nb += checkLigne(myMap, pos, team, 0, -19)
	if nb - 1 == 5 {
		//fmt.Printf("END 5 IN A ROW\n")
		return true
	}

	// diagonal /
	nb = checkLigne(myMap, pos, team, 0, -18)
	nb += checkLigne(myMap, pos, team, 0, 18)
	if nb - 1 == 5 {
		//fmt.Printf("END 5 IN A ROW\n")
		return true
	}

	// diagonal \
	nb = checkLigne(myMap, pos, team, 0, -20)
	nb += checkLigne(myMap, pos, team, 0, 20)
	if nb - 1 == 5 {
		//fmt.Printf("END 5 IN A ROW\n")
		return true
	}
	return false;
}

func getIndexCasePlayed(oldMap []protocol.MapData, newMap []protocol.MapData) (int) {
	var i int = 0
	for ; i < len(oldMap); i++ {
		if (oldMap[i] != newMap[i]) {
			return i
		}
	}
	return -1
}
// function qui check la regle "LE DOUBLE-TROIS"

// function qui check s'il peut NIQUER une paire et s'il peut tej les deux entre (prendre plusieurs pair d'un coup)

func checkCase(myMap []protocol.MapData, pos int, team int) (bool) {
	if myMap[pos].Player == (team + 1) % 2 {
		return (true)
	}
	return (false)
}

func checkPair(myMap []protocol.MapData, pos int, team int) ([]protocol.MapData, int) {
	var emptyData protocol.MapData
	emptyData.Empty = true
	emptyData.Playable = true
	emptyData.Player = -1
	captured := 0
	if (pos - (19 * 3)) >= 0 {
		// NORD
		if checkCase(myMap, pos - (19 * 1), team) && checkCase(myMap, pos - (19 * 2), team) && checkCase(myMap, pos - (19 * 3), (team + 1) % 2) {
			myMap[pos - (19 * 1)] = emptyData
			myMap[pos - (19 * 2)] = emptyData
			captured += 2
		}
	}
	if (pos - (19 * 3) + 3) >= 0 && pos % 19 <= 15 {
		// NORD EST
		if checkCase(myMap, pos - (19 * 1) + 1, team) && checkCase(myMap, pos - (19 * 2) + 2, team) && checkCase(myMap, pos - (19 * 3) + 3, (team + 1) % 2) {
			myMap[pos - (19 * 1) + 1] = emptyData
			myMap[pos - (19 * 2) + 2] = emptyData
			captured += 2
		}
	}
	if (pos + 3) < 19 * 19 && pos % 19 <= 15 {
		// EST
		if checkCase(myMap, pos + 1, team) && checkCase(myMap, pos + 2, team) && checkCase(myMap, pos + 3, (team + 1) % 2) {
			myMap[pos + 1] = emptyData
			myMap[pos + 2] = emptyData
			captured += 2
		}
	}
	if (pos + (19 * 3) + 3) < 19 * 19 && pos % 19 <= 15 {
		// SUD EST
		if checkCase(myMap, pos + (19 * 1) + 1, team) && checkCase(myMap, pos + (19 * 2) + 2, team) && checkCase(myMap, pos + (19 * 3) + 3, (team + 1) % 2) {
			myMap[pos + (19 * 1) + 1] = emptyData
			myMap[pos + (19 * 2) + 2] = emptyData
			captured += 2
		}
	}
	if (pos + (19 * 3)) < 19 * 19 {
		// SUD
		if checkCase(myMap, pos + (19 * 1), team) && checkCase(myMap, pos + (19 * 2), team) && checkCase(myMap, pos + (19 * 3), (team + 1) % 2) {
			myMap[pos + (19 * 1)] = emptyData
			myMap[pos + (19 * 2)] = emptyData
			captured += 2
		}
	}
	if (pos + (19 * 3) - 3) < 19 * 19 && pos % 19 >= 3 {
		// SUD OUEST
		if checkCase(myMap, pos + (19 * 1) - 1, team) && checkCase(myMap, pos + (19 * 2) - 2, team) && checkCase(myMap, pos + (19 * 3) - 3, (team + 1) % 2) {
			myMap[pos + (19 * 1) - 1] = emptyData
			myMap[pos + (19 * 2) - 2] = emptyData
			captured += 2
		}
	}
	if (pos - 3) >= 0 && pos % 19 >= 3 {
		// OUEST
		if checkCase(myMap, pos - 1, team) && checkCase(myMap, pos - 2, team) && checkCase(myMap, pos - 3, (team + 1) % 2) {
			myMap[pos - 1] = emptyData
			myMap[pos - 2] = emptyData
			captured += 2
		}
	}
	if (pos - (19 * 3) - 3) >= 0 && pos % 19 >= 3 {
		// NORD OUEST
		if checkCase(myMap, pos - (19 * 1) - 1, team) && checkCase(myMap, pos - (19 * 2) - 2, team) && checkCase(myMap, pos - (19 * 3) - 3, (team + 1) % 2) {
			myMap[pos - (19 * 1) - 1] = emptyData
			myMap[pos - (19 * 2) - 2] = emptyData
			captured += 2
		}
	}
	return myMap, captured
}

func main() {

	//myMap := initMap()

	for x := 0; x < 19 * 19; x++ {
		//	fmt.Println("x:", x, " sum:", myMap[x].team) // Simple output.

	}

	flag.Parse()
	log.SetFlags(0)

	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
