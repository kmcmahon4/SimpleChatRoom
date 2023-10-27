package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var chatParticipants = make(map[net.Conn]bool)

func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listen.Close()
	fmt.Println("Awesome! The chat server is running on localhost:8080")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		chatParticipants[conn] = true

		go handleClient(conn)
	}
}

func broadcastMessage(sender net.Conn, message string) {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return
	}

	for chatParticipant := range chatParticipants {
		if chatParticipant == sender {
			continue
		}
		_, err := chatParticipant.Write([]byte(message))

		if err != nil {
			fmt.Println(err)
			chatParticipant.Close()
			delete(chatParticipants, chatParticipant)
		}
	}
}

func removeParticipant(conn net.Conn) {
	delete(chatParticipants, conn)
}

func handleClient(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	message := make(chan string)

	go readInput(conn, message)

	for {
		select {
		case msg := <-message:
			broadcastMessage(conn, msg)
		}
	}
}

func readInput(conn net.Conn, message chan string) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message <- msg
	}

	removeParticipant(conn)

	err := conn.Close()
	if err != nil {
		return
	}
}
