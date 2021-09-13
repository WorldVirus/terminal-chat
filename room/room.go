package room

import (
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

type Message struct {
	Text string `json:"text"`
}

type Room struct {
	clients          map[string]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan Message
}

const (
	allActiveUsersCommand = "/getusers"
)

func NewRoom() *Room {
	return &Room{
		clients:          make(map[string]*websocket.Conn),
		addClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		broadcastChan:    make(chan Message),
	}
}

func (r *Room) Handler(ws *websocket.Conn) {
	go r.run()

	r.addClientChan <- ws

	for {
		var message Message
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			r.broadcastChan <- Message{err.Error()}
			r.removeClient(ws)
			return
		}
		r.broadcastChan <- message
	}
}

func (r *Room) run() {
	for {
		select {
		case conn := <-r.addClientChan:
			r.addClient(conn)
		case conn := <-r.removeClientChan:
			r.removeClient(conn)
		case message := <-r.broadcastChan:
			r.broadcastMessage(message)
		}
	}
}

func (r *Room) removeClient(conn *websocket.Conn) {
	delete(r.clients, conn.LocalAddr().String())
}

func (r *Room) addClient(conn *websocket.Conn) {

	userOrigin := conn.RemoteAddr().String()
	indexUserName := strings.Index(userOrigin, "=")
	if indexUserName != -1 {
		m := Message{
			Text: userOrigin[indexUserName+1:] + " joined",
		}

		for _, conn := range r.clients {
			err := websocket.JSON.Send(conn, &m)
			if err != nil {
				log.Println("Error broadcasting message: ", err)
				return
			}
		}

		log.Println(conn.RemoteAddr().String())
		r.clients[userOrigin] = conn

		return
	}

	log.Println("Username is empty")
}

func (r *Room) broadcastMessage(message Message) {
	var userList string

	if message.Text == allActiveUsersCommand {
		for username := range r.clients {
			index := strings.Index(username, "=")

			userList += username[index+1:] + " "
		}
		message.Text = userList
	}

	for _, conn := range r.clients {
		err := websocket.JSON.Send(conn, message)
		if err != nil {
			log.Println("Error broadcasting message: ", err)
			return
		}
	}
}
