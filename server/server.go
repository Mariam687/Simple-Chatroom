package main

import (
	"chat-by-rpc/common"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

// --- SERVICE DEFINITION ---

// ChatService holds the state of our chat (the message history)
type ChatService struct {
	history []string
	mu      sync.Mutex // Mutex to safely handle concurrent access to history
}

// SendMessage is the RPC method called by the client.
// It stores a new message or returns the history.
func (h *ChatService) SendMessage(args common.SendMessageArgs, reply *common.HistoryReply) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check for special commands
	if strings.HasPrefix(args.Message, "/FETCH_HISTORY") {
		//Return the chat history
		reply.History = h.history
		return nil
	}

	// Store the message
	formattedMsg := fmt.Sprintf("[%s]: %s", args.Sender, args.Message)
	h.history = append(h.history, formattedMsg)
	fmt.Printf("New message received: %s\n", formattedMsg)

	// Return the complete chat history after storing the message
	reply.History = h.history
	return nil
}

// --- SERVER MAIN LOGIC ---

func main() {
	// Setup Signal Handling (for Ctrl+C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("Fatal: Listener error: %v", err)
	}
	defer listener.Close()
	fmt.Println("Server running on 127.0.0.1:1234. Waiting for clients...")

	// Launch a goroutine to wait for the Ctrl+C signal
	go func() {
		sig := <-sigs
		fmt.Printf("\nReceived signal: %s. Shutting down...\n", sig)
		listener.Close() // Close the listener to break the main loop
	}()

	// Register Service
	rpc.RegisterName(common.ChatServiceName, &ChatService{})

	// Accept connections forever
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if the error is due to the listener being intentionally closed (graceful shutdown)
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "accept" {
				fmt.Println("Server listener closed. Shutting down gracefully.")
				return // Exit main function
			}
			log.Printf("Warning: Accept error: %v", err)
			continue
		}
		// Serve each connection concurrently
		go rpc.ServeConn(conn)
	}
}
