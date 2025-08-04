package server

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/models"
	"strings"
)

func HandleCleint(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte(Banner + "\n[ENTER YOUR NAME]:  "))
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if name == "" {
		conn.Write([]byte("Name cannot be empty. Dissconnecting.\n"))
		return
	}

	client := models.Client{Conn: conn, Name: name}
	ClientsMutex.Lock()
	if len(Clients) >= 10 {
		conn.Write([]byte("Server full. Try later.\n"))
		ClientsMutex.Unlock()
		return
	}

	Clients[conn] = client
	ClientsMutex.Unlock()

	//send history
	ClientsMutex.Lock()
	for _, m := range MessageHist {
		conn.Write([]byte(m + "\n"))
	}
	ClientsMutex.Unlock()

	Broacast(fmt.Sprintf("%s has joined our chat", name), conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		formatted := FormatMessage(name, msg)
		SaveMessage(formatted)
		Broacast(formatted, conn)
	}

	ClientsMutex.Lock()
	delete(Clients, conn)
	ClientsMutex.Unlock()
	Broacast(fmt.Sprintf("%s has left our chat...", name), conn)
}
