# glorp-irc: An IRC server (and client) written in Go

A lightweight IRC (Internet Relay Chat) server and client implementation built with Go, featuring a terminal user interface (TUI) for the client.

## Features

- **IRC Server**: Full-featured IRC server implementation
- **IRC Client**: Terminal-based client with an intuitive interface
- **Easy Setup**: Simple configuration and startup process
- **Cross-platform**: Works on Linux, macOS, and Windows

## Prerequisites

- Go 1.24.3 or later
- Terminal that supports ANSI colors (for the best client experience)

## Installation & Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configuration (Optional)

Create a `.env` file in the root directory to customize the server address:

```bash
echo "SERVER_ADDRESS=0.0.0.0:8800" > .env
```

**Note**: If no `.env` file is provided, the server defaults to `0.0.0.0:8800` (localhost:8800).

### 3. Running the Server

Start the IRC server:

```bash
go run ./cmd/server/
```

The server will start listening on the configured address (default: localhost:8800).

### 4. Running the Client

In a new terminal, start the IRC client:

```bash
go run ./cmd/client/
```

Follow the on-screen instructions to:
- Connect to the server
- Join a chat room
- Start chatting!

## Usage

1. **Start the server** first using the command above
2. **Launch one or more clients** to connect to the server
3. **Join rooms** and start chatting with other connected users
4. Use standard IRC commands within the client interface

## Project Structure

```
glorp-irc/
├── cmd/
│   ├── client/         # Client application entry point
│   ├── server/         # Server application entry point
│   └── main.go         # Main entry point
├── internal/
│   ├── client/         # Client implementation
│   └── server/         # Server implementation
├── releases/          # Release artifacts
├── .env              # Environment configuration (optional)
├── go.mod            # Go module dependencies
└── README.md         # This file
```

## Dependencies

This project uses several Go libraries including:
- **Bubble Tea**: For the terminal user interface
- **Lip Gloss**: For styling the TUI
- **godotenv**: For environment variable management

## Development

### Build from source

- Server:
```bash
go build -o server ./cmd/server/
```

- Client
```bash
go build -o client ./cmd/client/
```
(The client executables are also present in the `releases` directory, if you'd prefer to directly download them instead.)

### Running Tests

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

This project is open source. Please check the repository for license details.

---


**Note**: The asterisk (*) in the project name indicates that the name was decided on a whim, and was meant to be a placeholder name, which is now never getting changed.