package main

import (
	"fmt"
	"log"

	"github.com/IMBoBx/glorp-irc/internal/server"
)

func main() {
	err := server.StartServer()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server listening on localhost:8800")
	server.AcceptConnections()
	
}
