package main

import (
	"flag"
	"log"
	"net/http"
	//"fmt"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	hub := newHub()
	go hub.run()
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
