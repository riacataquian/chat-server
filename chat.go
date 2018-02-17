// Chat is a TCP server that allow clients ("nc") to communicate to each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	ch   chan<- string
	name string
}

var (
	messages    = make(chan string) // all incoming client messages
	entering    = make(chan client) // handle arrival of new clients
	leaving     = make(chan client) // handle leaving clients
	doHeadcount = make(chan bool)   // all changes in client: arrival or leaving
	cleanups    = make(chan bool)   // handle idle client disconnection

	clients = make(map[client]bool)   // client set
	idle    = make(map[net.Conn]bool) // determine if client is idle
)

// disconnectIdle disconnects a client if idle for 5 minutes.
// It performs disconnection after a signal has been received from the cleanups channel.
// After the cleanups, it restarts the timer.
func disconnectIdle(conn net.Conn) {
	for {
		select {
		case <-cleanups:
			if c, ok := idle[conn]; ok && c {
				conn.Close()
			}

			go idleCounter()
		}
	}
}

// idleCounter spawn a goroutine that sends a signal to the cleanups channel
// every 5 minutes.
func idleCounter() {
	go func() {
		time.Sleep(5 * time.Minute) // 5 mins
		cleanups <- true
	}()
}

func broadcast() {
	for {
		// Multiplex using select.
		select {
		case msg := <-messages:
			for c := range clients {
				// Delegate to another goroutine, so this channel would be non-blocking.
				// Otherwise, if client is blocked, this will cause a goroutine leak.
				go send(c, msg)
			}

		case c := <-entering:
			clients[c] = true
			go idleCounter()
			doHeadcount <- true

		case c := <-leaving:
			delete(clients, c)
			close(c.ch)
			doHeadcount <- true

		}
	}
}

// send sends a message to the client.
func send(c client, msg string) {
	c.ch <- msg
}

// Creates a chat server
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
		// Performs headcount for the active clients.
		go headCount()
		// Listen to cleanups channel.
		go disconnectIdle(conn)
	}
}

// headCount lists and reports all active clients.
func headCount() {
	for {
		select {
		case <-doHeadcount:
			var msg string

			msg = "\nPerforming headcount.."

			for c := range clients {
				msg += "\n" + c.name
			}

			messages <- msg
		}
	}
}

// Creates a new client connection
func handleConn(conn net.Conn) {
	defer conn.Close()

	// Returns a new scanner that reads from client's connection.
	input := bufio.NewScanner(conn)

	// Write in client's connection.
	fmt.Fprint(conn, "Chat username: ")

	var who string
	if input.Scan() {
		who = input.Text()
	}

	// create client its own channel
	ch := make(chan string)
	go clientWriter(conn, ch)

	// Inform client of its identity
	ch <- "Joining as " + who

	// Announce all existing clients that `who` joined.
	// Note that this comes first, so newly joined clients
	// won't receive broadcast of its own arrival.
	messages <- who + " joined"

	c := client{ch, who}

	idle[conn] = true

	// Then add the client channel to the clients set
	// by sending an event to the `entering` channel.
	entering <- c

	// Wait for client's input..
	for input.Scan() {
		messages <- who + ": " + input.Text()
		idle[conn] = false
	}

	// This code is unreachable unless input.Scan() returns an error.
	leaving <- c

	// Broadcast leaving of clients.
	messages <- who + " is leaving"
}

// clientWriter writes to client's connection.
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
