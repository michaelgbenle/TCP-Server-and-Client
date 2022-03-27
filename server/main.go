package main

import (
	"bufio"

	"log"
	"net"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	openConnect = make(map[net.Conn]bool)
	newConnect  = make(chan net.Conn)
	deadConnect = make(chan net.Conn)
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	logFatal(err)

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			logFatal(err)

			openConnect[conn] = true
			newConnect <- conn

		}

	}()

	for {
		select {
		case conn := <-newConnect:
			broadcastMessage(conn)
		case conn := <-deadConnect:
			for item := range openConnect {
				if item == conn {
					break
				}
			}
			delete(openConnect, conn)
		}
	}

}

func broadcastMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		for item := range openConnect {
			if item != conn {
				item.Write([]byte(message))
			}
		}
	}
	deadConnect <- conn
}
