package common

// ChatService is the name under which we register the service on the server.
const ChatServiceName = "ChatService"

// SendMessageArgs is the structure for the arguments sent from the client to the server.
// Since the client can send a message or request history, we combine them.
type SendMessageArgs struct {
	Sender  string // Name of the client
	Message string // The message content (or "FETCH_HISTORY")
}

// HistoryReply is the structure for the response sent from the server to the client.
type HistoryReply struct {
	// A long list of all messages in the chat history
	History []string
}
