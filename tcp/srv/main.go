package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"runtime"
	"sync"
)

// from proj root -- go run tcp/srv/main.go
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Error listening: %s", err.Error())
	}
	defer listener.Close()
	log.Printf("Server listening on 127.0.0.1:8080")

	errs := make(chan error)
	maxWorkers := runtime.NumCPU()
	workers := make(chan net.Conn, maxWorkers)
	var wg sync.WaitGroup

	go func() {
		for err := range errs {
			fmt.Printf("%v", err)
		}
	}()

	for range maxWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for conn := range workers {
				handleClient(conn, errs)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s", err.Error())
		}
		go handleClient(conn, errs)
	}
}

func handleClient(conn net.Conn, errs chan error) {
	defer conn.Close()

	buf := make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		errs <- err
		return
	}

	log.Printf("Received request: %s", string(buf))

	// This is where I could use a pipe operator. Assuming type conv is a func,
	// randResp() |> []byte |> conn.Write
	if _, err := conn.Write([]byte(randResp())); err != nil {
		errs <- err
		return
	}
}

func randResp() string {
	if v := rand.Intn(2); v == 0 { // 0 || 1
		return "Not too bad, client.."
	}
	return "Not so great, client.."
}
