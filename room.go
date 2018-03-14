package main

import (
	"encoding/json"
	"github.com/quentin-sommer/gomoku-back/ia"
	"github.com/quentin-sommer/gomoku-back/protocol"
	"github.com/quentin-sommer/gomoku-back/referee"
	"log"
)

const (
	INIT        = "INIT"
	START       = "START"
	RECONNECTED = "RECONNECTED"
	END         = "END"
)

type Player interface {
}

type Room struct {
	clients       map[*Client]bool
	players       [2]Player
	broadcast     chan *MessageClient
	changingState chan string
	boardGame     []protocol.MapData
	turnsPlayed   [2]int
	capturedPawns [2]int
	state         string
	nbTurn        int
	id            int
	AiMode        bool
}

func newRoom(newID int, AiMode bool) *Room {
	room := &Room{
		clients:       make(map[*Client]bool),
		broadcast:     make(chan *MessageClient),
		changingState: make(chan string),
		state:         INIT,
		nbTurn:        0,
		id:            newID,
		AiMode:        AiMode,
	}
	if AiMode == true {
		room.players[1] = &AiPlayer{level: 2}
	}
	return room
}

func (r *Room) addClient(c *Client) {
	c.room = r
	r.clients[c] = true
	if r.state == END {
		log.Println("Going to spectator, game is finished.")
		c.conn.WriteJSON(protocol.SendEndOfGame(r.boardGame, r.turnsPlayed, r.capturedPawns, r.nbTurn%2))
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
			c.conn.WriteJSON(protocol.SendRefresh(r.boardGame, r.turnsPlayed, r.capturedPawns))
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

func (r *Room) checkEndRoom(end bool) bool {
	if end == true || r.capturedPawns[r.nbTurn%2] >= 10 {
		r.state = END
		log.Println("End of game")
		endOfGameJSON := protocol.SendEndOfGame(r.boardGame, r.turnsPlayed, r.capturedPawns, r.nbTurn%2)
		for client := range r.clients {
			client.conn.WriteJSON(endOfGameJSON)
		}
		return true
	} else {
		return false
	}
}

func (r *Room) run() {
	r.boardGame, r.turnsPlayed, r.capturedPawns = protocol.InitGameData()
	// force play middle case by the black player on first turn
	r.boardGame[180].Empty = false
	r.boardGame[180].Playable = false
	r.boardGame[180].Player = 1
	r.turnsPlayed[1] = 1
	for {
		select {
		case newState := <-r.changingState:
			log.Println("NewState: ", newState)
			switch newState {
			case START:
				entity, ok := r.players[0].(*Client)
				if ok == true {
					entity.conn.WriteJSON(protocol.SendStartOfGame(0))
					entity.conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.turnsPlayed, r.capturedPawns, -1))
				}
				entity, ok = r.players[1].(*Client)
				if ok == true {
					entity.conn.WriteJSON(protocol.SendStartOfGame(1))
					entity.conn.WriteJSON(protocol.SendRefresh(r.boardGame, r.turnsPlayed, r.capturedPawns))
				}
				log.Println("Starting Game in a room.")
			case RECONNECTED:
				for i := 0; i < len(r.players); i++ {
					entity, ok := r.players[i].(*Client)
					if ok == true {
						entity.conn.WriteJSON(protocol.SendStartOfGame(i))
						entity.conn.WriteJSON(protocol.SendRefresh(r.boardGame, r.turnsPlayed, r.capturedPawns))
						if r.nbTurn%2 == i {
							entity.conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.turnsPlayed, r.capturedPawns, -1))
						}
					}
				}
				log.Println("Reconnected.")
			}
			r.state = newState

		case message := <-r.broadcast:
			var typeJSON protocol.MessageIdle
			_ = json.Unmarshal(message.broadcast, &typeJSON)
			log.Println(typeJSON.Type)

			switch typeJSON.Type {
			case protocol.SET_AI_LEVEL:
				var setLevelJSON protocol.MessageSetAiLevel
				_ = json.Unmarshal(message.broadcast, &setLevelJSON)
				for i := 0; i < len(r.players); i++ {
					ai, ok := r.players[i].(*AiPlayer)
					if ok == true {
						log.Println("Setting ai : ", i, " at level ", setLevelJSON.Level)
						ai.level = setLevelJSON.Level
					}
				}
			case protocol.PLAY_TURN:
				var playTurnJSON protocol.MessagePlayTurn
				_ = json.Unmarshal(message.broadcast, &playTurnJSON)
				idx := referee.GetIndexCasePlayed(r.boardGame, playTurnJSON.Map)
				if idx == -1 {
					log.Println("Error, index played = -1")
				}

				var capturedPawns int
				var end, ok bool
				//log.Println("MAP BEFORE:\n", playTurnJSON.Map)
				capturedPawns, end, ok = referee.Exec(playTurnJSON.Map, idx)
				//fmt.Println("Je suis Player ", r.nbTurn % 2, playTurnJSON.Map[idx].Player)
				//fmt.Println("eval of play ", ia.Eval(&ia.MinMaxStruct{playTurnJSON.Map, int8(r.nbTurn % 2), 0, end, idx}))
				//log.Println("MAP AFTER:\n", playTurnJSON.Map)
				if ok == false {
					// Illegal action, play again
					entity, ok := r.players[r.nbTurn%2].(*Client)
					if ok == true {
						entity.conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.turnsPlayed, r.capturedPawns, -1))
					}
				} else {
					r.boardGame = playTurnJSON.Map
					r.turnsPlayed[r.nbTurn%2] += 1
					r.capturedPawns[r.nbTurn%2] += capturedPawns

					if !r.checkEndRoom(end) {
						r.nbTurn += 1
						if message.client == r.players[0] || message.client == r.players[1] {
							entity, ok := r.players[r.nbTurn%2].(*Client)
							if ok == true {
								entity.conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.turnsPlayed, r.capturedPawns, -1))
							}
							ai, ok := r.players[r.nbTurn%2].(*AiPlayer)
							if ok == true {
								refreshJSON := protocol.SendRefresh(r.boardGame, r.turnsPlayed, r.capturedPawns)
								// refresh before ia turn
								for client := range r.clients {
									client.conn.WriteJSON(refreshJSON)
								}
								idx := ia.MinMax(r.boardGame, int8(r.nbTurn%2), ai.level)
								r.boardGame[idx].Empty = false
								r.boardGame[idx].Playable = false
								r.boardGame[idx].Player = int8(r.nbTurn % 2)
								captured, end, _ := referee.Exec(r.boardGame, idx)
								r.turnsPlayed[r.nbTurn%2] += 1
								r.capturedPawns[r.nbTurn%2] += captured
								if !r.checkEndRoom(end) {
									r.nbTurn += 1
									entity, ok := r.players[r.nbTurn%2].(*Client)
									if ok == true {
										entity.conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.turnsPlayed, r.capturedPawns, -1))
										suggestedMove := ia.MinMax(r.boardGame, int8(r.nbTurn%2), 3)
										entity.conn.WriteJSON(protocol.SendSuggestedMove(suggestedMove))
									}
								}
							}
							refreshJSON := protocol.SendRefresh(r.boardGame, r.turnsPlayed, r.capturedPawns)
							for client := range r.clients {
								if client != r.players[r.nbTurn%2] {
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
