package main

import (
	"log"

	"github.com/IMBoBx/glorp-irc/internal/client"
)

func main() {
	conn, err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go client.ReceiveMessage(conn)
	client.SendMessage(conn)
}