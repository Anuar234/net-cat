package client

import (
	"bufio"
	"fmt"
	"net"
)

func StartClient(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("COnnection error:", err)
		return
	}
	defer conn.Close()

	//recieve from server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	//send to server
	for {
		text, err := ReadInput()
		if err != nil {
			fmt.Println("Input error:", err)
			break
		}
		if text != "" {
			_, err := conn.Write([]byte(text + "\n"))
			if err != nil {
				fmt.Println("Write failed", err)
				break
			}
		}
	}
}
