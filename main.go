package main

import (
	"fmt"
	"net-cat/server"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		// Default port
		server.StartServer("8989")
	} else if len(os.Args) == 2 {
		// Custom port
		server.StartServer(os.Args[1])
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		fmt.Println("Then connect with: nc localhost $port")
	}
}
