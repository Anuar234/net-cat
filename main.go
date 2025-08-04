package main

import (
	"fmt"
	"net-cat/client"
	"net-cat/server"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		server.StartServer("8989")
	} else if len(os.Args) == 2 {
		server.StartServer(os.Args[1])
	} else if len(os.Args) == 3 && os.Args[1] == "client" {
		client.StartClient(os.Args[2])
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		fmt.Println("		./TCPChat cleint $host:$port")
	}
}
