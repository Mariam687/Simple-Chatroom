# Go RPC Chat Server and Client

This repository contains a simple, persistent chat application built using Go's built-in Remote Procedure Call (`net/rpc`) package. The client can send messages, fetch chat history, and gracefully handle server shutdowns and restarts.

## 🚀 Features

* **Concurrency:** The server handles multiple clients simultaneously using Goroutines (`go rpc.ServeConn(conn)`).
* **Persistent History:** All messages are stored in a server-side list and returned on every message/history request.
* **Graceful Shutdown:** Both the client and server handle `Ctrl+C` (SIGINT) for clean exit.
* **Client Persistence:** The client automatically attempts to reconnect to the server if the connection is lost.
* **Full-Line Input:** Uses `bufio` to correctly read entire messages, including spaces.

## 📁 Project Structure

The project is structured into separate modules for clean separation of concerns:

## 🛠️ Setup and Run

### Prerequisites

You must have Go installed.

### 1. Initialize the Module

Navigate to the root directory of the project (`chat-rpc/`) and ensure your Go module is initialized:

```bash
go mod init chat-rpc
go mod tidy
```

### 2. Start the Server

Open your first terminal (or a split terminal in VS Code). The server will listen on port `1234`.

```bash
go run ./server
```

### 2. Start the Client

```bash
go run ./client
```

## ⚙️ Usage and Commands

| Command/Action | Client Behavior |
| :--- | :--- |
| **`hello world`** | Message is sent to the server, stored, and the full history is returned. |
| **`exit`** | Client gracefully closes its connection and exits the application. |
| **`Ctrl+C` (Client)** | Client gracefully exits via OS signal handler. |
| **`Ctrl+C` (Server)** | Server gracefully shuts down the listener. All active clients will automatically attempt to reconnect. |
| **`/history`** | Requests the full chat history from the server without sending a new message. |

## 📺 Demo

A short demonstration of the client-server interaction and the graceful reconnection feature.

![Demo of the Go RPC Chat Application](assets/chat_demo.gif)

