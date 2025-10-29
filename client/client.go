package main

import (
	"bufio"
	"chat-by-rpc/common"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/signal" //To handle OS signals
	"strings"
	"syscall" //For specific signal types (like SIGINT)
)

// printHistory clears the screen and prints the full chat history
func printHistory(history []string) {
	// ANSI Escape Codes to clear terminal screen and move cursor home
	fmt.Print("\033[H\033[2J")
	fmt.Println("--- CHAT HISTORY ---")
	for _, msg := range history {
		fmt.Println(msg)
	}
	fmt.Println("--------------------")
}

func main() {
	// 1. Get Client Name
	var clientName string
	fmt.Print("Enter your name: ")
	fmt.Scanln(&clientName)

	// 2. Setup Signal Handling (for Ctrl+C)
	sigs := make(chan os.Signal, 1)
	// We want to catch the interrupt signal (Ctrl+C)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//When the user presses Ctrl+C, the OS sends SIGINT,
	// signal.Notify pushes the signal into the sigs channel, the goroutine wakes up, and execution continues past <-sigs

	// Launch a goroutine that waits for the signal
	go func() {
		<-sigs // Blocks until a signal is received

		fmt.Println("\n[Client] Interrupt signal received (Ctrl+C). Shutting down...")
		// Since we used 'defer client.Close()' below, the connection will be closed.
		// We can now exit the program immediately.
		os.Exit(0) // Exit cleanly
	}()
	// Note: The signal handler is set up early so it's always running.

	// 3. Connect to server
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		// Handle connection errors if the server is down
		log.Fatalf("Error connecting to server. Is the server running? %v", err)
	}
	defer client.Close() // This ensures the connection is closed when main exits (including via os.Exit(0))

	// Print the welcome message and initial prompt
	fmt.Printf("Connected as %s. Type 'exit' to quit.\n", clientName)
	reader := bufio.NewReader(os.Stdin)

	// 4. Main Communication Loop
	for {
		fmt.Printf("[%s]> ", clientName)

		// Use ReadString to read the entire line, including spaces
		input, err := reader.ReadString('\n')
		if err != nil {
			// This catches unexpected terminal closure or other IO issues (like EOF)
			log.Printf("Input error: %v. Shutting down.", err)
			break
		}

		message := strings.TrimSpace(input)

		// 5. Handle "exit" command
		if strings.ToLower(message) == "exit" {
			fmt.Println("Client shutting down...")
			return // Exits the main function, triggering defer client.Close()
		}

		// Set the message to "FETCH_HISTORY" if the input is empty or a special command
		if message == "" || strings.ToLower(message) == "/history" {
			message = "/FETCH_HISTORY"
		}

		// 6. Build arguments and reply structures
		args := common.SendMessageArgs{
			Sender:  clientName,
			Message: message,
		}
		var reply common.HistoryReply

		// 7. Call the Remote Procedure
		err = client.Call(common.ChatServiceName+".SendMessage", args, &reply)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "shut down") {
				// Handle server going down gracefully
				log.Fatalf("Server connection lost: %v. Exiting.", err)
			}
		}

		// 8. Display the full chat history
		printHistory(reply.History)
	}
}
