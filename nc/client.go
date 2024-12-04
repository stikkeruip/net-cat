package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please provide the server IP and port as arguments")
		os.Exit(0)
	}
	IP := os.Args[1]
	intPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error while converting string to")
		os.Exit(0)
	}
	if intPort < 1024 || intPort > 65535 {
		fmt.Println("Port needs to be between 1024 and 65535.")
		os.Exit(1)
	}
	port := os.Args[2]

	conn := connect(IP, port)
	if conn == nil {
		return
	}
	go serverResponse(conn)
	handleUserInput(conn)
}

func connect(IP, port string) net.Conn {
	conn, err := net.Dial("tcp", IP+":"+port)
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

			// Check if the message is the quit command
			if message == "/quit" {
				fmt.Println("Disconnecting from the server...")
				conn.Write([]byte(message + "\n"))
				conn.Close()
				os.Exit(0)
			}

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
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		if err == io.EOF {
			fmt.Println("\nServer closed the connection. Exiting...")
		} else {
			fmt.Println("\nLost connection to server:", err)
		}
		os.Exit(0)
	}
}
