package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Connect() (net.Conn, error){
	conn, err := net.Dial("tcp", "localhost:8800")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			log.Fatal(err)
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			conn.Close()
			log.Fatal(err)
		}
	}
}

func ReceiveMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			log.Fatal(err)
		}

		fmt.Println(message)
	}
}