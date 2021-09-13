package main

import (
	"flag"
	"log"
	"net/http"

	r "github.com/WorldVirus/terminal-chat/room"
	"golang.org/x/net/websocket"
)

var (
	port = flag.String("port", "8080", "port used for ws connection")
)

func main() {
	flag.Parse()

	log.Fatal(server(*port))
}

func server(port string) error {
	room := r.NewRoom()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		room.Handler(ws)
	}))

	s := http.Server{Addr: ":" + port, Handler: mux}
	log.Print("Server start at ws://localhost:", port)
	return s.ListenAndServe()
}
