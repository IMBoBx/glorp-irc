package main

import (
	"log"

	"github.com/IMBoBx/glorp-irc/internal/client"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	conn, err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()


	m := client.InitialModel(conn)
	p := tea.NewProgram(m)

	go client.StartReceiver(conn, p)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
