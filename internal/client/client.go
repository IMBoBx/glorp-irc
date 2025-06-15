// package client

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// func Connect() (net.Conn, error) {
// 	conn, err := net.Dial("tcp", "localhost:8800")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return conn, nil
// }

// func SendMessage(conn net.Conn) {
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		message, err := reader.ReadString('\n')
// 		fmt.Print("\033[2K\r")
// 		if err != nil {
// 			conn.Close()
// 			log.Fatal(err)
// 		}

// 		_, err = conn.Write([]byte(message))
// 		if err != nil {
// 			conn.Close()
// 			log.Fatal(err)
// 		}
// 	}
// }

// func ReceiveMessage(conn net.Conn) tea.Msg {
// 	reader := bufio.NewReader(conn)

// 	for {
// 		message, err := reader.ReadString('\n')
// 		if err != nil {
// 			conn.Close()
// 			log.Fatal(err)
// 		}

// 		fmt.Print(message)
// 		// m.incoming <- message
// 	}

// 	// return func() tea.Msg {
// 	// 	reader := bufio.NewReader(conn)
// 	// 	msg, err := reader.ReadString('\n')
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	return incomingMsg(msg)
// 	// }
// }


// internal/client/client.go
package client

import (
	"bufio"
	"log"
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
			p.Quit()
			log.Fatal(err)
		}
		p.Send(IncomingMsg(message))
	}
}
