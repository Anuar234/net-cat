package client

import (
	"bufio"
	"os"
	"strings"
)

//readInput reads user input from terminal and trims it

func ReadInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}
