package main

import (
	"log"
	"net"
)

func main() {
	// initialize main server, responsible for all incoming commands and state of the rooms
	s := newServer()

	// execute the run coroutine
	go s.run()

	// Start TCP server
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server", err.Error())
	}

	defer listener.Close()
	log.Printf("Started server on :8888")

	// Accept all new clients 
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %s", err.Error())
			continue
		}

		// Initialize every new client
		go s.newClient(conn)
	}
}