## terminal-chat

### How start ?


start server `go run server.go --port="8080"`, to select a port using flag `port`

client start `go run ./client/client.go --server="ws://localhost:8080/" --username="test"`, to select a port and username using flags `ports`,  `username`

To get list of active user, using command `/getusers`


### Description: 
* The program should consist of two parts: a client and a server.
* Communication protocol between clients and server - websocket
when starting the client from the command line
* Should be possible to specify the server address and port, as well as the username. eg `./client --server = localhost: 1234 --name = "FooBar"`
when starting the server 
* Should be possible to set the port to listen on. eg `./server --port = 1234 `
* After launch, all clients get into one common "room" on the server
* While the client is working, should be possible to write a message to the general chat right in the terminal window. Messages should be received by all users who are currently online.
all users should be alerted when a new user is added
* Should be possible to send the `/getusers` command to the chat (a message with the text “/ getusers”). in reply