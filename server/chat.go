package server

import (
	"fmt"
	"net"
	"net-cat/models"
	"sync"
	"time"
)

var (
	Clients = make(map[net.Conn]models.Client, 10)

	MessageHist = []string{}

	ClientsMutex = sync.Mutex{}
)

func Broadcast(message string, sender net.Conn) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	for conn := range Clients {
		if conn != sender {
			conn.Write([]byte(message + "\n"))
		}
	}
}

func SaveMessage(msg string) {
	ClientsMutex.Lock()
	MessageHist = append(MessageHist, msg)
	ClientsMutex.Unlock()
}

func FormatMessage(name, msg string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]: %s", t, name, msg)
}
