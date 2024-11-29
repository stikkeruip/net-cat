package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	conn := connect()
	go serverResponse(conn)
	handleUserInput(conn)
}

func connect() net.Conn {
	// Connect to the server
	conn, err := net.Dial("tcp", "127.0.0.1:8080") // Replace with server IP/port if not local
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}
	//defer conn.Close()
	fmt.Println("Connected to server")

	return conn
}

func handleUserInput(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			message := scanner.Text() + "\n" // Add a newline for the server to process
			// Clear the user's input after sending
			fmt.Print("\033[1A\033[2K") // ANSI escape code to clear the last line
			_, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Println("Error sending message to server:", err)
				break
			}
		} else {
			fmt.Println("Error reading input or EOF")
			break
		}
	}
}

func serverResponse(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// Read messages from the server
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server")
			os.Exit(0)
		}
		// Display the server's response
		fmt.Print(message)
	}
}
