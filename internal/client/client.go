package client

import (
	"bufio"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

func Connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:8800")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func StartReceiver(conn net.Conn, p *tea.Program) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			p.Send(IncomingMsg(message))
			p.Quit()
			return
		}
		p.Send(IncomingMsg(message))
	}
}
