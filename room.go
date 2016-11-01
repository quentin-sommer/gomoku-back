package main

import (
	"encoding/json"
	"log"

	"./protocol"
	"./referee"
)

const (
	INIT        = "INIT"
	START       = "START"
	RECONNECTED = "RECONNECTED"
  END = "END"
)

type Room struct {
	clients        map[*Client]bool
	players        [2]*Client
	broadcast      chan *MessageClient
	changingState  chan string
	boardGame      []protocol.MapData
	availablePawns [2]int
	capturedPawns  [2]int
	state          string
	nbTurn         int
	id             int
}

func newRoom(newID int) *Room {
	return &Room{
		clients:       make(map[*Client]bool),
		broadcast:     make(chan *MessageClient),
		changingState: make(chan string),
		state:         INIT,
		nbTurn:        0,
		id:            newID,
	}
}

func (r *Room) addClient(c *Client) {
	c.room = r
	r.clients[c] = true
  if r.state == END {
    c.conn.WriteJSON(protocol.SendEndOfGame(r.boardGame, r.availablePawns, r.capturedPawns, r.nbTurn % 2))
  } else {
  	if r.players[0] == nil {
  		r.players[0] = c
  		if r.state != INIT {
  			r.changingState <- RECONNECTED
  		}
  	} else if r.players[1] == nil {
  		r.players[1] = c
  		if r.state != INIT {
  			r.changingState <- RECONNECTED
  		}
  	} else {
  		c.conn.WriteJSON(protocol.SendStartOfGame(2))
  		c.conn.WriteJSON(protocol.SendRefresh(r.boardGame, r.availablePawns, r.capturedPawns))
  	}
  	if r.state == INIT && r.players[0] != nil && r.players[1] != nil {
  		r.changingState <- START
  	}
  }
}

func (r *Room) delClient(c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
	}
	if r.players[0] == c {
		r.players[0] = nil
	} else if r.players[1] == c {
		r.players[1] = nil
	}
}

func (r *Room) run() {
	r.boardGame, r.availablePawns, r.capturedPawns = protocol.InitGameData()
	// force play middle case by the black player on first turn
	r.boardGame[180].Empty = false
	r.boardGame[180].Playable = false
	r.boardGame[180].Player = 1
	r.availablePawns[1] = 59
	for {
		select {
		case newState := <-r.changingState:
      log.Println("NewState: ", newState)
			switch newState {
			case START:
				r.players[0].conn.WriteJSON(protocol.SendStartOfGame(0))
				r.players[1].conn.WriteJSON(protocol.SendStartOfGame(1))
				r.players[0].conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.availablePawns, r.capturedPawns, -1))
				r.players[1].conn.WriteJSON(protocol.SendRefresh(r.boardGame, r.availablePawns, r.capturedPawns))
				log.Println("Starting Game in a room.")
			case RECONNECTED:
				if r.players[0] != nil {
					r.players[0].conn.WriteJSON(protocol.SendStartOfGame(0))
					r.players[0].conn.WriteJSON(protocol.SendRefresh(r.boardGame, r.availablePawns, r.capturedPawns))
				}
				if r.players[1] != nil {
					r.players[1].conn.WriteJSON(protocol.SendStartOfGame(1))
					r.players[1].conn.WriteJSON(protocol.SendRefresh(r.boardGame, r.availablePawns, r.capturedPawns))
				}
				if r.players[r.nbTurn%2] != nil {
					r.players[r.nbTurn%2].conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.availablePawns, r.capturedPawns, -1))
				}
				log.Println("Reconnected.")
			}
			r.state = newState

		case message := <-r.broadcast:
			var typeJSON protocol.MessageIdle
			_ = json.Unmarshal(message.broadcast, &typeJSON)
			log.Println(typeJSON.Type)

			switch typeJSON.Type {
			case protocol.PLAY_TURN:
				var playTurnJSON protocol.MessagePlayTurn
				_ = json.Unmarshal(message.broadcast, &playTurnJSON)
				idx := referee.GetIndexCasePlayed(r.boardGame, playTurnJSON.Map)
				if idx == -1 {
					log.Println("Error, index played = -1")
				}

        var capturedPawns int
        var end, ok bool
        playTurnJSON.Map, capturedPawns, end, ok = referee.Exec(playTurnJSON.Map, idx)
        log.Println("CapturedPawns: ", capturedPawns)
        log.Println("End: ", end)
        log.Println("OK: ", ok)
        if ok == false && false {
          // Illegal action, play again
          r.players[r.nbTurn % 2].conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.availablePawns, r.capturedPawns, -1))
        } else {
          r.boardGame = playTurnJSON.Map
          r.availablePawns[r.nbTurn % 2] -= 1
          r.capturedPawns[r.nbTurn % 2] += capturedPawns

          if end == true || r.capturedPawns[r.nbTurn % 2] >= 10 {
            r.state = END
            log.Println("End of game")
            endOfGameJSON := protocol.SendEndOfGame(r.boardGame, r.availablePawns, r.capturedPawns, r.nbTurn % 2)
            for client := range r.clients {
              client.conn.WriteJSON(endOfGameJSON)
            }
          } else {

  					r.nbTurn += 1
  					if message.client == r.players[0] || message.client == r.players[1] {
  						if r.players[r.nbTurn % 2] != nil {
  							r.players[r.nbTurn % 2].conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.availablePawns, r.capturedPawns, -1))
  						}
  						refreshJSON := protocol.SendRefresh(r.boardGame, r.availablePawns, r.capturedPawns)
  						for client := range r.clients {
  							if client != r.players[r.nbTurn % 2] {
  								client.conn.WriteJSON(refreshJSON)
  							}
  						}
  					}
          }
				}
			}
		}
	}
}
