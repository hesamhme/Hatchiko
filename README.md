
# Hatchiko

Hatchiko is a minimal queue manager written in Go, designed to provide lightweight and efficient message queue management. It supports basic queuing functionality with features like topics and a server-client architecture for handling requests over a network.

## Features

- **Queue Management**: Enqueue and dequeue messages with thread-safe operations.
- **Topics**: Organize messages into topics to enable categorized message processing.
- **Server-Client Architecture**: Interact with the queue through a networked server.

## Project Structure

The repository is organized as follows:

- **`queue/`**: Contains the core implementation of queue data structures and related functionalities.
- **`server/`**: Includes server-related code to handle client requests and manage topics.
  - `server.go`: Server implementation for handling client connections and requests.
  - `topic.go`: Topic implementation for organizing messages.
- **`tests/`**: Includes test cases to validate the functionality and reliability of the queue manager.
- **`main.go`**: The entry point of the application, initializing and starting the server.

## Usage

### Prerequisites

- Go 1.18 or higher installed on your system.

### Running the Application

1. Clone the repository:

   ```bash
   git clone https://github.com/hesamhme/Hatchiko.git
   cd Hatchiko
   ```

2. Run the application:

   ```bash
   go run main.go
   ```

3. The server will start on `127.0.0.1:8080`. You can interact with it through a client.

### Testing

To run the test cases:

```bash
go test ./tests/
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contribution

Contributions are welcome! Feel free to fork the repository, submit pull requests, or open issues for suggestions or bug reports.

## Author

[Hesam HME](https://github.com/hesamhme)
