package server

import (
	"log"
	"net"
)

func StartServer(port string) {
	address := ":" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	log.Println("Listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error acccepting:", err)
			continue
		}
		go HandleCleint(conn)
	}
}
