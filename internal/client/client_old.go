package client

import (
	crand "crypto/rand"
	"fmt"
	"log"
	// "math/rand"
	"net"
	"time"
)

const MAX = 5
const PORT = "8800"

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newConnection(i int, c chan bool) {
	conn, err := net.Dial("tcp", "localhost:"+PORT)
	handleErr(err)
	defer conn.Close()

	conn.Write(fmt.Appendf(nil, "NICK gorb%d\n", i))

	counter := 0
	for range time.Tick(3 * time.Second) {

		conn.Write([]byte(crand.Text() + "\n"))

		if counter++; counter == MAX {
			conn.Write([]byte(fmt.Sprintf("EXIT gorb%d\n", i)))
			c <- true
			break
		}
	}
}

func init() {
	done := make(chan bool, MAX)

	for i := range MAX {
		go newConnection(i, done)
	}

	for range MAX {
		<-done
	}
}
