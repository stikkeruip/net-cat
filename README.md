# Net Cat

Net Cat is a custom implementation of the `netcat` utility, written in Go. It allows you to perform TCP/IP-based networking tasks such as creating server-client connections, sending and receiving data, and testing network configurations.

## Features

- Send and receive data over TCP/IP connections.
- Create and manage server-client communication.
- Simple and lightweight implementation.
- Written in Go, making it cross-platform and efficient.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Git installed on your machine

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/stikkeruip/net-cat.git
   ```
2. Navigate to the project directory:
   ```bash
   cd net-cat
   ```
3. Build the project:
   ```bash
   go build -o net-cat
   ```

### Usage

Run the `net-cat` executable with the desired options. Examples:

#### Start a Server

```bash
./net-cat -l -p 8080
```

This starts a server listening on port 8080.

#### Connect to a Server

```bash
./net-cat 127.0.0.1 8080
```

This connects to a server running at `127.0.0.1` on port 8080.

#### Send a Message

Once connected, type a message and press Enter to send it. The message will appear on the server side.

#### Additional Options

For a list of all available options, run:

```bash
./net-cat -h
```

## Contributing

Contributions are welcome! Feel free to fork the repository and submit a pull request. Ensure your changes adhere to the project's coding standards.

1. Fork the project.
2. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add your feature"
   ```
4. Push to the branch:
   ```bash
   git push origin feature/your-feature-name
   ```
5. Open a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments

- Inspired by the original `netcat` utility.
- Thanks to the Go community for their excellent documentation and support.

---

Happy hacking!

