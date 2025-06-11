package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

const PORT = "8800"

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleConn(conn net.Conn, connections map[string]net.Conn, mu *sync.Mutex) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	line, _ := reader.ReadString('\n')

	if !strings.HasPrefix(line, "NICK ") {
		conn.Write([]byte("Expected NICK <username> as the first message!"))
		return
	}

	nick := strings.TrimSpace(strings.TrimPrefix(line, "NICK "))

	mu.Lock()
	connections[nick] = conn
	mu.Unlock()

	fmt.Println(nick, "joined!")

	for {
		line, _ := reader.ReadString('\n')

		if strings.Compare("EXIT " + nick + "\n", line) == 0 {

			fmt.Println(nick, "left!")
			break
		}

		fmt.Printf("%-12s > %s", nick, line)
	}

}

func init() {
	ln, err := net.Listen("tcp", "localhost:"+PORT)
	handleErr(err)
	defer ln.Close()

	var (
		mu          sync.Mutex
		connections = make(map[string]net.Conn)
	)

	for {
		conn, err := ln.Accept()
		handleErr(err)

		go handleConn(conn, connections, &mu)
	}
}
