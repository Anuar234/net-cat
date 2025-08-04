package server

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/models"
	"strings"
	"time"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()

	var name string
	reader := bufio.NewReader(conn)

	// name selection loop
	for {
		// send banner and prompt for name
		conn.Write([]byte(Banner + "\n[ENTER YOUR NAME]: "))

		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		name = strings.TrimSpace(input)

		// validate name is not empty
		if name == "" || strings.ContainsFunc(name, func(r rune) bool {
			return r != 10 && (r < 32 || r == 127)
		}) {
			continue
		}

		// check server capacity and name uniqueness
		ClientsMutex.Lock()
		if len(Clients) >= 10 {
			conn.Write([]byte("Server full. Try later.\n"))
			ClientsMutex.Unlock()
			return
		}

		// check if name is already taken
		nameExists := false
		for _, client := range Clients {
			if client.Name == name {
				nameExists = true
				break
			}
		}

		if nameExists {
			conn.Write([]byte("Name already taken. Please choose another name.\n"))
			ClientsMutex.Unlock()
			continue
		}

		// name is unique and valid, add client
		client := models.Client{Conn: conn, Name: name}
		Clients[conn] = client
		ClientsMutex.Unlock()
		break
	}

	// send chat history
	ClientsMutex.Lock()
	for _, m := range MessageHist {
		conn.Write([]byte(m + "\n"))
	}
	ClientsMutex.Unlock()

	// announce user joined
	joinMsg := fmt.Sprintf("%s has joined our chat...", name)
	SaveMessage(joinMsg)
	Broadcast(joinMsg, conn)

	// send prompt for first message
	m := time.Now().Format("2006-01-02 15:04:04")
	conn.Write([]byte(fmt.Sprintf("[%s][%s]:", m, name)))

	// handle messages
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)
		if msg == "" || strings.ContainsFunc(msg, func(r rune) bool {
			return r != 10 && (r < 32 || r == 127)
		}) {
			continue
		}

		//if msg.Message == "" || !utf8.ValidString(msg.Message) ||
		//strings.ContainsFunc(msg.Message, func(r rune) bool {
		// return !unicode.IsGraphic(r)
		// }) {

		formatted := FormatMessage(name, msg)
		SaveMessage(formatted)
		Broadcast(formatted, conn)

		// send prompt for next message
		t := time.Now().Format("2006-01-02 15:04:05")
		conn.Write([]byte(fmt.Sprintf("[%s][%s]:", t, name)))

		// t := time.Now().Format("2006-01-02 15:04:05")
		// return fmt.Sprintf("[%s][%s]: %s", t, name, msg)
	}

	// clean up when client disconnects
	ClientsMutex.Lock()
	delete(Clients, conn)
	ClientsMutex.Unlock()

	leaveMsg := fmt.Sprintf("%s has left our chat...", name)
	SaveMessage(leaveMsg)
	Broadcast(leaveMsg, conn)
}
