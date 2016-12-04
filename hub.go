// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "log"
  "encoding/json"
  "./protocol"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
  // Registered clients.
  clients    map[*Client]bool

  // Inbound messages from the clients.
  broadcast  chan *MessageClient

  // Register requests from the clients.
  register   chan *Client

  // Unregister requests from clients.
  unregister chan *Client

  rooms      map[int]*Room
}

func newHub() *Hub {
  return &Hub{
    broadcast:  make(chan *MessageClient),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    clients:    make(map[*Client]bool),
    rooms:            make(map[int]*Room),
  }
}

func (h *Hub) run() {
  for {
    select {
    case client := <-h.register:
    //log.Println("New Connection")
      h.clients[client] = true
    case client := <-h.unregister:
    //log.Println("Deconnection")
      if _, ok := h.rooms[client.room.id]; ok {
        client.room.delClient(client)
        if len(client.room.clients) <= 0 {
          delete(h.rooms, client.room.id)
        }
      }
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
      }
    case message := <-h.broadcast:
      var typeJSON protocol.MessageIdle
      _ = json.Unmarshal(message.broadcast, &typeJSON)

    //Getting the Room ID
      if (typeJSON.Type == protocol.ENTER_ROOM) {
        var roomJSON protocol.MessageEnterRoom
        _ = json.Unmarshal(message.broadcast, &roomJSON)

        //Checking if Room ID exist
        if _, ok := h.rooms[roomJSON.Room]; !ok {
          log.Println("Creating the Room: ", roomJSON.Room)
          h.rooms[roomJSON.Room] = newRoom(roomJSON.Room, roomJSON.AiMode)
          go h.rooms[roomJSON.Room].run()
        }
        h.rooms[roomJSON.Room].addClient(message.client)
        log.Println(h.rooms[roomJSON.Room].players)
      }

    /*for client := range h.clients {
      select {
      case client.send <- message:
      default:
        close(client.send)
        delete(h.clients, client)
      }
    }
    */
    }
  }
}
