// Chat is a TCP server that allow clients ("nc") to communicate to each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string

var (
	messages = make(chan string) // all incoming client messages
	entering = make(chan client) // handle arrival of new clients
	leaving  = make(chan client) // handle leaving clients
)

func broadcast() {
	clients := make(map[client]bool) // client set

	for {
		// Multiplex using select.
		select {
		case msg := <-messages:
			for c := range clients {
				c <- msg
			}

		case c := <-entering:
			clients[c] = true

		case c := <-leaving:
			delete(clients, c)
			close(c)

		}
	}
}

// main creates a TCP chat server.
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	// Its important to delegate this to a separate goroutine,
	// since this infinitely loops, this will block
	// unless an error is encountered.
	go broadcast()

	for {
		// Waiting for incoming request..
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle new request.
		go handleConn(conn)
	}
}

// Creates a new client connection
func handleConn(conn net.Conn) {
	defer conn.Close()

	// create client its own channel
	ch := make(chan string)
	go clientWriter(conn, ch)

	// Inform client of its identity
	who := conn.RemoteAddr().String()
	ch <- "You are " + who

	// Announce all existing clients that `who` joined.
	// Note that this comes first, so newly joined clients
	// won't receive broadcast of its own arrival.
	messages <- who + " joined"

	// Then add the client channel to the clients set
	// by sending an event to the `entering` channel.
	entering <- ch

	// Wait for client's input..
	input := bufio.NewScanner(conn)
	for input.Scan() {
		// NOTE: ignoring err from input.Text()
		messages <- who + ":" + input.Text()
	}

	// This code is unreachable unless input.Scan() returns an error.
	leaving <- ch

	// Broadcast leaving of clients.
	messages <- who + " is leaving"
}

// clientWriter writes to client's connection.
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
