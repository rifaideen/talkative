package talkative

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
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
	ErrCallback = errors.New("callback cannot be empty")    // Error for missing callback function.
	ErrMessage  = errors.New("messages cannot be empty")    // Error for empty message list.
	ErrInvoke   = errors.New("unable to invoke ollama api") // Error for failing to call the Ollama API.
	ErrEncoding = errors.New("unable to encode")            // Error for problems encoding data to JSON.
)

// Callback function type used for handling individual chat responses and errors.
// Takes a pointer to a ChatResponse struct and an error as arguments.
type ChatCallBack func(*ChatResponse, error)

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

// ChatRequest struct represents the request body sent to the Ollama API for chat processing.
type ChatRequest struct {
	Model    string    `json:"model"`    // The model to be used for processing the chat.
	Messages []Message `json:"messages"` // List of messages to be processed.
}

// ChatResponse struct represents the response received from the Ollama API after processing chat messages.
type ChatResponse struct {
	Model              string    `json:"model"`                // The model used for processing.
	Message            Message   `json:"message"`              // The response message.
	CreatedAt          time.Time `json:"created_at"`           // Time the response was created on the server.
	Done               bool      `json:"done"`                 // Indicates if processing is complete.
	TotalDuration      int       `json:"total_duration"`       // Total processing time in milliseconds.
	LoadDuration       int       `json:"load_duration"`        // Time spent loading the model (milliseconds).
	PromptEvalCount    int       `json:"prompt_eval_count"`    // Number of prompt evaluations performed.
	PromptEvalDuration int       `json:"prompt_eval_duration"` // Time spent on prompt evaluation (milliseconds).
	EvalCount          int       `json:"eval_count"`           // Number of overall evaluations performed.
	EvalDuration       int       `json:"eval_duration"`        // Time spent on overall evaluation (milliseconds).
}

// New function creates a new Client instance for interacting with the Ollama API.
// Takes the base URL of the Ollama API as an argument.
func New(url string) *Client {
	client := &http.Client{} // Create a new HTTP client instance.

	return &Client{
		urls: map[string]string{
			"chat": url + "/api/chat", // Define the chat endpoint URL based on the provided base URL.
		},
		client: client,
	}
}

// Initiates a chat process and asynchronously handles responses through a callback function.
//
// This function takes a callback function (`cb`) and a variable number of messages (`msgs`) as arguments.
// It performs the following steps:
//  1. Validates the callback and message arguments.
//  2. Prepares a request body with the messages and model information.
//  3. Sends a POST request to the chat endpoint from this client.
//  4. Handles the response status code and potential errors.
//  5. Launches a goroutine to process the chat response asynchronously.
//  6. Returns a channel (`<-chan bool`) that signals completion of the chat process and any errors encountered.
//
// The callback function (`cb`) is responsible for handling individual chat responses and errors.
// The completion channel (`<-chan bool`) allows the caller to track the progress of the chat process if needed.
//
// Note that the channel (`chDone`) is not explicitly closed in this example. However, the goroutine
// running `processChat` terminates naturally after sending the completion signal (`true`),
// effectively indicating no more data will be received on the channel.
func (c *Client) Chat(cb ChatCallBack, msgs ...Message) (<-chan bool, error) {
	if cb == nil {
		return nil, ErrCallback
	}

	if len(msgs) == 0 {
		return nil, ErrMessage
	}

	request := ChatRequest{
		Model:    "llama2",
		Messages: msgs,
	}
	body := &bytes.Buffer{}

	if err := json.NewEncoder(body).Encode(request); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrEncoding, err)
	}

	res, err := c.client.Post(c.urls["chat"], "application/json", body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: please make sure ollama server is running and url is correct", ErrInvoke)
	}

	chDone := make(chan bool)

	go func() {
		c.processChat(res.Body, cb)

		chDone <- true
	}()

	return chDone, nil
}

// Processes the chat response from the server asynchronously.
//
// This function takes an io.ReadCloser object (`body`) representing the response body
// and a callback function (`cb`) for handling individual responses and errors.
// It iterates through the response, decoding each message and invoking the callback for processing.
//
// In case of errors during decoding or processing, the callback is invoked with the error
// and processing stops. The function closes the response body before exiting.
func (c *Client) processChat(body io.ReadCloser, cb ChatCallBack) {
	defer body.Close()

	for {
		var response ChatResponse

		err := json.NewDecoder(body).Decode(&response)

		if err == io.EOF {
			return
		}

		if err != nil {
			cb(nil, err)

			return
		}

		cb(&response, nil)
	}
}
