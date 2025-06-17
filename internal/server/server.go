package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Connection struct {
	Username string
	Conn     net.Conn
}

type Room struct {
	Name   string
	Buffer chan string
	Users  map[string]*Connection
	mu     sync.RWMutex
}

type Server struct {
	Listener net.Listener
	Rooms    map[string]*Room
	mu       sync.RWMutex
}

var (
	tcp *Server
	max = 32
)

func StartServer() error {
	if tcp != nil {
		return errors.New("server already started")
	}

	ln, err := net.Listen("tcp", "localhost:8800")
	if err != nil {
		return err
	}

	tcp = &Server{
		Listener: ln,
		Rooms:    make(map[string]*Room),
	}

	return nil
}

func AcceptConnections() {
	for {
		conn, err := tcp.Listener.Accept()
		if err != nil {
			conn.Close()
			continue
		}

		go ReadConnection(conn)
	}
}

func ReadConnection(conn net.Conn) {
	var (
		username string
		room     string
		err      error
	)

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		conn.Write([]byte(err.Error()))
		time.Sleep(time.Second)
		conn.Close()
	}

	username, room, err = handleFirstMessage(message)
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		time.Sleep(time.Second)
		conn.Close()
	}

	if err := CreateRoom(room); err != nil && !strings.Contains(err.Error(), "already exists") {
		conn.Write([]byte(err.Error() + "\n"))
		time.Sleep(time.Second)
		conn.Close()
	}

	err = AddConnection(username, room, conn)
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		time.Sleep(time.Second)
		conn.Close()
	}

}

func handleFirstMessage(message string) (string, string, error) {
	message = strings.TrimSpace(message)
	if !strings.HasPrefix(message, "/join ") {
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
	room = normalizeName(room)

	tcp.mu.Lock()
	defer tcp.mu.Unlock()
	if tcp.Rooms[room] != nil {
		return errors.New("room with name " + room + " already exists")
	}

	r := &Room{
		Name:   room,
		Buffer: make(chan string, max),
		Users:  make(map[string]*Connection),
	}
	tcp.Rooms[room] = r
	go BroadcastMessages(r)

	return nil
}

func BroadcastMessages(room *Room) {
	for message := range room.Buffer {
		room.mu.RLock()
		for _, user := range room.Users {
			user.Conn.Write([]byte(message))
		}
		room.mu.RUnlock()
	}
}

func AddConnection(username, room string, conn net.Conn) error {
	room = normalizeName(room)
	username = normalizeName(username)

	tcp.mu.RLock()
	r := tcp.Rooms[room]
	tcp.mu.RUnlock()
	if r == nil {
		return errors.New("room does not exist")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Users[username] != nil {
		return errors.New("username taken")
	}

	r.Users[username] = &Connection{
		Username: username,
		Conn:     conn,
	}

	tcp.Rooms[room].Buffer <- "\033[95m" + username + " joined!\033[0m\n"

	go WatchMessage(username, room)

	return nil
}

func WatchMessage(username, room string) error {
	tcp.mu.RLock()
	r := tcp.Rooms[room]
	tcp.mu.RUnlock()
	if r == nil {
		return errors.New("room does not exist")
	}

	r.mu.RLock()
	conn := r.Users[username]
	r.mu.RUnlock()
	if conn == nil {
		return errors.New("connection doesn't exist, try reconnecting")
	}

	for {
		reader := bufio.NewReader(conn.Conn)
		message, err := reader.ReadString('\n')

		if err != nil || message == "/leave "+username+"\n" {
			r.mu.Lock()

			r.Buffer <- "\033[91m" + username + " left!\033[0m\n"

			conn.Conn.Close()
			delete(r.Users, username)
			empty := len(r.Users) == 0
			r.mu.Unlock()

			if empty {
				tcp.mu.Lock()
				close(r.Buffer)
				delete(tcp.Rooms, room)
				tcp.mu.Unlock()
			}

			break
		}

		r.Buffer <- fmt.Sprintf("[\033[96m%s\033[0m] %s", username, message)

	}

	return nil
}

func normalizeName(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.ToLower(s)), " ", "-")
}
