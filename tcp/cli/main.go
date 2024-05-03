package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// from proj root -- go run tcp/cli/main.go
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Error connecting: %s\n", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("Enter message to send to server: ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}

		msg = msg[:len(msg)-1]

		if _, err = conn.Write([]byte(msg)); err != nil {
			log.Fatalf("Error sending request: %s\n", err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("Error reading response: %s\n", err)
		}
		// include only the n read bytes in resp printout
		fmt.Printf("Server response: %s\n", string(buf[:n]))
	}
}
