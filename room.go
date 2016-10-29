package main

import (
  "./protocol"
  "log"
  "encoding/json"
)

const
(
  INIT = "INIT"
  START = "START"
  RECONNECTED = "RECONNECTED"
)

type Room struct {
  clients map[*Client]bool
  players [2]*Client
  broadcast chan *MessageClient
  changingState chan string
  boardGame []protocol.MapData
  state string
  nbTurn int
  id int
}

func newRoom(newID int) (*Room) {
  return &Room{
    clients:        make(map[*Client]bool),
    broadcast:      make(chan *MessageClient),
    changingState:  make(chan string),
    state:          INIT,
    nbTurn:         0,
    id:             newID,
  }
}

func (r *Room) addClient(c *Client) {
  c.room = r
  r.clients[c] = true
  if (r.players[0] == nil) {
    r.players[0] = c
    if (r.state != INIT) {
      r.changingState <- RECONNECTED
    }
  } else if (r.players[1] == nil) {
    r.players[1] = c
    if (r.state != INIT) {
      r.changingState <- RECONNECTED
    }
  } else {
    c.conn.WriteJSON(protocol.SendStartOfGame(2))
    c.conn.WriteJSON(protocol.SendRefresh(r.boardGame))
  }
  if (r.state == INIT && r.players[0] != nil && r.players[1] != nil) {
    r.changingState <- START
  }
}

func (r *Room) delClient(c *Client) {
  if _, ok := r.clients[c]; ok {
    delete(r.clients, c)
  }
  if (r.players[0] == c) {
    r.players[0] = nil
  } else if (r.players[1] == c) {
    r.players[1] = nil
  }
}

func (r *Room) run() {
  r.boardGame = initMap()
  for {
    select {
    case newState := <- r.changingState:
      switch newState {
      case START:
        r.players[0].conn.WriteJSON(protocol.SendStartOfGame(0))
        r.players[1].conn.WriteJSON(protocol.SendStartOfGame(1))
        r.players[0].conn.WriteJSON(protocol.SendPlayTurn(r.boardGame))
        log.Println("Starting Game in a room.")
      case RECONNECTED:
        if (r.players[0] != nil) {
          r.players[0].conn.WriteJSON(protocol.SendStartOfGame(0))
          r.players[0].conn.WriteJSON(protocol.SendRefresh(r.boardGame))
        }
        if (r.players[1] != nil) {
          r.players[1].conn.WriteJSON(protocol.SendStartOfGame(1))
          r.players[1].conn.WriteJSON(protocol.SendRefresh(r.boardGame))
        }
        if (r.players[r.nbTurn % 2] != nil) {
          r.players[r.nbTurn % 2].conn.WriteJSON(protocol.SendPlayTurn(r.boardGame))
        }
        log.Println("Reconnected.")
      }
      r.state = newState

    case message := <- r.broadcast:
      var typeJSON protocol.MessageIdle
			_ = json.Unmarshal(message.broadcast, &typeJSON)
      log.Println(typeJSON.Type)

      switch typeJSON.Type {
      case protocol.PLAY_TURN:
        var playTurnJSON protocol.MessagePlayTurn
        _ = json.Unmarshal(message.broadcast, &playTurnJSON)
        // Si le tour est bon +1 au tour
        idx := getIndexCasePlayed(r.boardGame, playTurnJSON.Map)

        checkEnd(playTurnJSON.Map, idx, playTurnJSON.Map[idx].Player)

        playTurnJSON.Map, _ = checkPair(playTurnJSON.Map, idx, playTurnJSON.Map[idx].Player)

        if (true) {
          r.boardGame = playTurnJSON.Map
          r.nbTurn += 1

          if (message.client == r.players[0] || message.client == r.players[1]) {
            if (r.players[r.nbTurn % 2] != nil) {
              r.players[r.nbTurn % 2].conn.WriteJSON(protocol.SendPlayTurn(r.boardGame))
            }
            for client := range r.clients {
              client.conn.WriteJSON(protocol.SendRefresh(r.boardGame))
            }
          }
        }
      }
    }
  }
}
