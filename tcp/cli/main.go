package main

import (
	"fmt"
	"log"
	"net"
)

// from proj root -- go run tcp/cli/main.go
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Error connecting: %s\n", err)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte("Hello, server! How are you?")); err != nil {
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
