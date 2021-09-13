package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Text string `json:"text"`
}

var (
	url      = flag.String("server", "localhost:8080", "specify url to use")
	username = flag.String("username", "", "specify name for user")
)

func main() {
	flag.Parse()

	ws, err := connect(*url, *username)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	var m Message
	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				log.Println("Error receiving message: ", err.Error())
				break
			}
			log.Println("Message: ", m)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		m := Message{
			Text: text,
		}
		err = websocket.JSON.Send(ws, m)
		if err != nil {
			log.Println("Error sending message: ", err.Error())
			break
		}
	}
}

func connect(url, username string) (*websocket.Conn, error) {
	generateOrigin := mockedIP() + "?=" + username
	return websocket.Dial(url, "", generateOrigin)
}

func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}
