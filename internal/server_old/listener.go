package server_old

import (
	"errors"
	"log"
	"net"
	"strings"
)

type Connection struct {
	User     string
	Conn     net.Conn
	ToRoom   chan string
	FromRoom chan string
}

type Room struct {
	Name   string
	Users  map[string]*Connection
	Buffer chan string
}

type TcpServer struct {
	Listener net.Listener
	Rooms    map[string]*Room
}

var (
	tcp *TcpServer
	max = 32 // limit to 32 users per Room at a time
)

func StartServer() (*TcpServer, error) {
	if tcp != nil {
		return nil, errors.New("server already started")
	}

	ln, err := net.Listen("tcp", "localhost:8800")
	if err != nil {
		log.Fatal(err)
	}

	tcp = &TcpServer{
		Listener: ln,
		Rooms:    make(map[string]*Room),
	}

	go func() {
		for {
			conn, err := tcp.Listener.Accept()
			if err != nil {
				conn.Close()
				continue
			}

			go func(conn net.Conn) {
				data := make([]byte, 1024)
				n, err := conn.Read(data)
				if err != nil {
					conn.Close()
					return
				}

				message := string(data[:n])
				username, room, err := handleFirstMessage(message)
				if err != nil {
					conn.Close()
					return
				}

				CreateRoom(room)

				err = AddConnection(username, room, conn)
				if err != nil {
					conn.Close()
				}
			}(conn)

		}
	}()

	return tcp, nil
}

func handleFirstMessage(message string) (string, string, error) {
	if !strings.HasSuffix(message, "/join ") {
		return "", "", errors.New(`first message must be in format "/join <username> <room-name>". spaces are not allowed in both username and room-name`)
	}

	arr := strings.Split(message, " ")

	if len(arr) != 3 {
		return "", "", errors.New(`invalid args. first message must be in format "/join <username> <room-name>". spaces are not allowed in both username and room-name`)
	}

	username, room := arr[1], strings.TrimSuffix(arr[2], "\n")
	return username, room, nil
}

func CreateRoom(room string) error {
	room = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(room)), " ", "-")

	if tcp.Rooms[room] != nil {
		return errors.New("room with name " + room + " already exists")
	}

	r := &Room{
		Name:   room,
		Buffer: make(chan string, max),
		Users:  make(map[string]*Connection),
	}

	tcp.Rooms[room] = r

	go BroadcastMessage(r)

	return nil
}

func AddConnection(username, room string, conn net.Conn) error {

	room = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(room)), " ", "-")
	username = strings.TrimSpace(username)

	if tcp.Rooms[room].Users[username] != nil {
		return errors.New("username taken")
	}

	tcp.Rooms[room].Users[username] = &Connection{
		User:     username,
		Conn:     conn,
		ToRoom:   make(chan string),
		FromRoom: make(chan string, max),
	}

	tcp.Rooms[room].Users[username].ToRoom <- username + " joined!\n"

	go WatchMessage(username, room)

	return nil
}

/*
func SendMessage(message, username, room string) error {
	msgChan := tcp.Rooms[room].Users[username].ToRoom

	if msgChan == nil {
		return errors.New("connection doesn't exist, try reconnecting")
	}

	if message == "/leave "+username+"\n" {
		msgChan <- username + "left!\n"
		close(msgChan)

		tcp.Rooms[room].Users[username].Conn.Close()

		delete(tcp.Rooms[room].Users, username)
		if len(tcp.Rooms[room].Users) == 0 {
			delete(tcp.Rooms, room)
		}

		return nil
	}

	msgChan <- message

	return nil
}
*/

func WatchMessage(username, room string) error {
	msgChan := tcp.Rooms[room].Users[username].ToRoom

	if msgChan == nil {
		return errors.New("connection doesn't exist, try reconnecting")
	}

	for {
		message := <-msgChan

		if message == "/leave "+username+"\n" {
			msgChan <- username + "left!\n"
			close(msgChan)

			tcp.Rooms[room].Users[username].Conn.Close()

			delete(tcp.Rooms[room].Users, username)
			if len(tcp.Rooms[room].Users) == 0 {
				close(tcp.Rooms[room].Buffer)
				delete(tcp.Rooms, room)
			}

			break
		}

		tcp.Rooms[room].Buffer <- message

	}
	return nil
}

func BroadcastMessage(room *Room) {
	buffer := room.Buffer

	for msg := range buffer {
		for _, conn := range room.Users {
			conn.FromRoom <- msg
		}
	}
}
