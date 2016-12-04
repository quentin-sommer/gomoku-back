package main

import (
  "encoding/json"
  "log"
  "math/rand"
  "./ia"
  "./protocol"
  "./referee"
)

const (
  INIT = "INIT"
  START = "START"
  RECONNECTED = "RECONNECTED"
  END = "END"
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
    room.players[1] = &AiPlayer{level : 3}
  }
  return room
}

func (r *Room) addClient(c *Client) {
  c.room = r
  r.clients[c] = true
  if r.state == END {
    log.Println("Going to spectator, game is finished.")
    c.conn.WriteJSON(protocol.SendEndOfGame(r.boardGame, r.turnsPlayed, r.capturedPawns, r.nbTurn % 2))
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
            if r.nbTurn % 2 == i {
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
        playTurnJSON.Map, capturedPawns, end, ok = referee.Exec(playTurnJSON.Map, idx)
        //log.Println("MAP AFTER:\n", playTurnJSON.Map)
        if ok == false {
          // Illegal action, play again
          entity, ok := r.players[r.nbTurn % 2].(*Client)
          if ok == true {
            entity.conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.turnsPlayed, r.capturedPawns, -1))
          }
        } else {
          r.boardGame = playTurnJSON.Map
          r.turnsPlayed[r.nbTurn % 2] += 1
          r.capturedPawns[r.nbTurn % 2] += capturedPawns

          if end == true || r.capturedPawns[r.nbTurn % 2] >= 10 {
            r.state = END
            log.Println("End of game")
            endOfGameJSON := protocol.SendEndOfGame(r.boardGame, r.turnsPlayed, r.capturedPawns, r.nbTurn % 2)
            for client := range r.clients {
              client.conn.WriteJSON(endOfGameJSON)
            }
          } else {
            ia.MinMax(r.boardGame, r.nbTurn % 2, 3)
            r.nbTurn += 1
            if message.client == r.players[0] || message.client == r.players[1] {
              entity, ok := r.players[r.nbTurn % 2].(*Client)
              if ok == true {
                entity.conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.turnsPlayed, r.capturedPawns, -1))
              }
              _, ok = r.players[r.nbTurn % 2].(*AiPlayer)
              if ok == true {
                for ; ; {
                  idxRand := rand.Int() % protocol.MAP_SIZE
                  tmpmap := make([]protocol.MapData, len(r.boardGame))
                  copy(tmpmap, r.boardGame)
                  if tmpmap[idxRand].Empty {
                    tmpmap[idxRand].Empty = false
                    tmpmap[idxRand].Player = r.nbTurn % 2

                    tmpmap, capturedPawns, end, ok = referee.Exec(tmpmap, idxRand)
                    if ok == true {
                      r.boardGame = tmpmap
                      r.turnsPlayed[r.nbTurn % 2] += 1
                      r.capturedPawns[r.nbTurn % 2] += capturedPawns
                      break
                    }
                  }
                }
                r.nbTurn += 1
                entity, ok := r.players[r.nbTurn % 2].(*Client)
                if ok == true {
                  entity.conn.WriteJSON(protocol.SendPlayTurn(r.boardGame, r.turnsPlayed, r.capturedPawns, -1))
                }
              }
              refreshJSON := protocol.SendRefresh(r.boardGame, r.turnsPlayed, r.capturedPawns)
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
