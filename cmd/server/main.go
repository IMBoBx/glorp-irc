package main

import (
	"fmt"
	"log"

	"github.com/IMBoBx/glorp-irc/internal/server"
)

func main() {
	addressChan, err := server.StartServer()
	if err != nil {
		log.Fatal(err)
	}

	addr := <-addressChan
	fmt.Println("Server listening on ", addr)
	server.AcceptConnections()

}
