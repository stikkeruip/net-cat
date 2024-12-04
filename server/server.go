package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	Conn        net.Conn // The client's connection
	Username    string   // The client's username
	HasUsername bool     // Whether the username has been set
}

const (
	maxClients     = 10
	welcomeMessage = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | '  \\Zq\n_)      \\.___.,|     .\n\\____   )MMMMMP|   .\n     '-'       '--'\n"
)

var (
	clients = make(map[net.Conn]*Client)
	chatLog *os.File
)

func main() {
	start()
}

func start() {
	var err error
	chatLog, err = os.OpenFile("chat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening chat.log:", err)
		return
	}
	listener, err := net.Listen("tcp", "127.0.0.1:8080") // Listen on port 8080
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	listen(listener)
}

func listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("New connection accepted")
		fmt.Printf("New connection from: %s\n", conn.RemoteAddr().String())

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when done

	client := &Client{
		Conn:        conn,
		Username:    "",
		HasUsername: false,
	}

	clients[conn] = client

	welcomeClient(client)

	reader := bufio.NewReader(conn)

	for {
		if !client.HasUsername {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Client disconnected: %v\n", conn.RemoteAddr())
				delete(clients, conn)
				return
			}
			message = strings.TrimSpace(message)

			if validName(client, message) {
				client.Username = message
				client.HasUsername = true
				printChatLog(client)
				fmt.Printf("%s has joined the chat.\n", client.Username)

				broadcastJoin(client)
			} else {
				conn.Write([]byte("[ENTER YOUR NAME]: "))
			}
			continue
		}

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client disconnected: %v\n", conn.RemoteAddr())
			broadcastLeave(client)
			delete(clients, conn)
			return
		}

		message = strings.TrimSpace(message)
		if message != "" {
			broadcastMessage(client, message)
		}
	}
}

func welcomeClient(client *Client) {
	_, err := client.Conn.Write([]byte(welcomeMessage))
	if err != nil {
		fmt.Println("Error sending message to client:", err)
		return
	}

	_, err = client.Conn.Write([]byte("[ENTER YOUR NAME]: "))
	if err != nil {
		fmt.Println("Error sending prompt to client:", err)
		return
	}
}

func broadcastMessage(sender *Client, message string) {

	formattedMessage := fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), sender.Username, message)
	if formattedMessage != "\n" {

		addToLog(formattedMessage)
		for _, client := range clients {
			if !client.HasUsername {
				continue
			}
			_, err := client.Conn.Write([]byte(formattedMessage))
			if err != nil {
				fmt.Printf("Error broadcasting to user, %s", client.Username)
			}
		}
	}
}

func broadcastJoin(sender *Client) {
	formattedMessage := fmt.Sprintf("%s has joined the chat!\n", sender.Username)
	addToLog(formattedMessage)
	for _, client := range clients {
		if !client.HasUsername {
			continue
		}
		_, err := client.Conn.Write([]byte(formattedMessage))
		if err != nil {
			fmt.Printf("Error broadcasting join, %s", client.Username)
		}
	}
}

func broadcastLeave(sender *Client) {
	formattedMessage := fmt.Sprintf("%s has left the chat!\n", sender.Username)
	addToLog(formattedMessage)
	for _, client := range clients {
		if !client.HasUsername {
			continue
		}
		_, err := client.Conn.Write([]byte(formattedMessage))
		if err != nil {
			fmt.Printf("Error broadcasting leave, %s", client.Username)
		}
	}
}

func validName(client *Client, name string) bool {
	// Check length of the new name
	if len(name) < 3 || len(name) > 16 {
		client.Conn.Write([]byte("Username must be 3-16 characters\n"))
		return false
	}

	// Check for empty or whitespace-only names
	if strings.TrimSpace(name) == "" {
		client.Conn.Write([]byte("Username cannot be empty\n"))
		return false
	}

	// Check for taken usernames
	for _, v := range clients {
		if v.Username == name {
			client.Conn.Write([]byte("Username is taken, please try another\n"))
			return false
		}
	}

	// Check for valid characters
	validChars := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validChars.MatchString(name) {
		client.Conn.Write([]byte("Username can only contain letters, numbers, underscores, and hyphens\n"))
		return false
	}

	return true
}

func addToLog(message string) {
	if chatLog == nil {
		fmt.Println("Error: chat log file is not open")
		return
	}

	_, err := chatLog.WriteString(message)
	if err != nil {
		fmt.Println("Error writing to chat log:", err)
	}
	// Ensure the log entry is immediately written to the file
	err = chatLog.Sync()
	if err != nil {
		fmt.Println("Error syncing chat log:", err)
	}
}

func printChatLog(client *Client) {
	// Read the entire chat log file
	content, err := os.ReadFile("chat.log")
	if err != nil {
		fmt.Println("Error reading chat log:", err)
		return
	}

	// Send the content to the client
	_, err = client.Conn.Write(content)
	if err != nil {
		fmt.Println("Error sending chat log to client:", err)
		return
	}
}
