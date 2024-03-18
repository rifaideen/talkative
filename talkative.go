package talkative

import (
	"errors"
	"net/http"
	"strings"
)

// Define an enum-like type to represent different user roles in the chat system.
type Role string

const (
	// User role for regular users.
	USER Role = "user"

	// Assistant role for AI assistants or chatbots.
	ASSISTANT Role = "assistant"
)

// Pre-defined errors used throughout the code for consistency.
var (
	ErrUrl      = errors.New("url cannot be empty")         // Error for missing URL
	ErrCallback = errors.New("callback cannot be empty")    // Error for missing callback function.
	ErrMessage  = errors.New("messages cannot be empty")    // Error for empty message list.
	ErrInvoke   = errors.New("unable to invoke ollama api") // Error for failing to call the Ollama API.
	ErrEncoding = errors.New("unable to encode")            // Error for problems encoding data to JSON.
	ErrDecoding = errors.New("unable to decode")            // Error for problems encoding data to JSON.
)

// Client struct holds information for interacting with the Ollama API.
type Client struct {
	urls   map[string]string // Stores endpoint URLs for the Ollama API.
	client *http.Client      // Holds an http.Client instance for making HTTP requests.
}

// Message struct represents a single message sent or received in the chat.
type Message struct {
	Role    Role   `json:"role"`    // Role of the sender (user or assistant).
	Content string `json:"content"` // Content of the message.
}

// New function creates a new Client instance for interacting with the Ollama API.
// Takes the base URL of the Ollama API as an argument.
func New(url string) (*Client, error) {
	url = strings.Trim(url, " ")

	if url == "" {
		return nil, ErrUrl
	}

	client := &http.Client{} // Create a new HTTP client instance.

	return &Client{
		urls: map[string]string{
			"chat": url + "/api/chat", // Define the chat endpoint URL based on the provided base URL.
		},
		client: client,
	}, nil
}
