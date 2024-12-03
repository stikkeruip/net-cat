package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn := connect()
	if conn == nil {
		return
	}
	go serverResponse(conn)
	handleUserInput(conn)
}

func connect() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}
	fmt.Println("Connected to server")
	return conn
}

func handleUserInput(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			message := scanner.Text() // Read user input

			// Send the message to the server
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println("Error sending message to server:", err)
				break
			}

			// Remove the previous line (user's input)
			fmt.Print("\033[F\033[K")
		} else {
			fmt.Println("Error reading input or EOF")
			break
		}
	}
}

func serverResponse(conn net.Conn) {
	for {
		// Copy data from the server to stdout
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Println("Disconnected from server")
			os.Exit(0)
		}
	}
}
