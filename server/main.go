package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	maxClients     = 10
	welcomeMessage = `Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    ` + "`" + `.       | '  \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
[ENTER YOUR NAME]:`
)

func main() {
	start()
}

func start() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080") // Listen on port 8080
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	for {
		// Step 2: Accept an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("New connection accepted")

		// Display connection details
		fmt.Printf("New connection from: %s\n", conn.RemoteAddr().String())

		// Close the connection (no further communication for simplicity)
		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when done
	reader := bufio.NewReader(conn)
	user := conn.RemoteAddr().String()

	for {
		// Read message from the client
		message, err := reader.ReadString('\n') // Read until newline
		if err != nil {
			fmt.Printf("Client disconnected: %s\n", user)
			break
		}

		fmt.Printf("Message received from %s: %s", user, message)

		// Send the message back to the client
		_, err = conn.Write([]byte(user + ": " + message))
		if err != nil {
			fmt.Println("Error sending message to client:", err)
			break
		}
	}
}
