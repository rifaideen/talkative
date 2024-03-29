package talkative

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ChatMessage struct represents a single message sent or received in the chat.
type ChatMessage struct {
	Role    Role   `json:"role"`    // Role of the sender (user or assistant).
	Content string `json:"content"` // Content of the message.
}

// CompletionParams represents the advanced parameters (Optional) to be supplied to the completion request.
type ChatParams struct {
	Format    string                 `json:"format,omitempty"`     // The format to be used in the completion response
	Options   map[string]interface{} `json:"options,omitempty"`    // The additional model parameters  listed in the Modelfile documentation
	Template  string                 `json:"template,omitempty"`   // The prompt template to use (overrides what is defined in the Modelfile)
	Stream    *bool                  `json:"stream,omitempty"`     // Whether to get response in single shot rather than streaming
	KeepAlive string                 `json:"keep_alive,omitempty"` // How long to keep the model will stay loaded into the memory. Default to 5m(inutes)
}

// Callback function type used for handling individual chat responses and errors.
// Takes a pointer to a ChatResponse struct and an error as arguments.
type ChatCallBack func(*ChatResponse, error)

// PlainChatCallBack function type used for handling individual chat responses and errors.
// Takes a string and an error as arguments.
type PlainChatCallBack func(string, error)

// ChatRequest struct represents the request body sent to the Ollama API for chat processing.
type ChatRequest struct {
	Model    string        `json:"model"`    // The model to be used for processing the chat.
	Messages []ChatMessage `json:"messages"` // List of messages to be processed.

	*ChatParams `json:",omitempty"` // The additional parameters for the chat
}

// ChatResponse struct represents the response received from the Ollama API after processing chat messages.
type ChatResponse struct {
	Model       string      `json:"model"`      // The model used for processing.
	Message     ChatMessage `json:"message"`    // The response message.
	CreatedAt   time.Time   `json:"created_at"` // Time the response was created on the server.
	Done        bool        `json:"done"`       // Indicates if processing is complete.
	ChatMetrics             // The metrics associated about the chat
}

type ChatMetrics struct {
	TotalDuration      int `json:"total_duration"`       // Total processing time in milliseconds.
	LoadDuration       int `json:"load_duration"`        // Time spent loading the model (milliseconds).
	PromptEvalCount    int `json:"prompt_eval_count"`    // Number of prompt evaluations performed.
	PromptEvalDuration int `json:"prompt_eval_duration"` // Time spent on prompt evaluation (milliseconds).
	EvalCount          int `json:"eval_count"`           // Number of overall evaluations performed.
	EvalDuration       int `json:"eval_duration"`        // Time spent on overall evaluation (milliseconds).
}

// Initiates a chat process and asynchronously handles responses through a callback function.
//
// This function takes model name, callback function (`cb`) and a variable number of messages (`msgs`) as arguments.
// It performs the following steps:
//  1. Validates the model, callback and message arguments.
//  2. When model is empty, it uses the model to DEFAULT_MODEL to perform the operation.
//  3. Prepares a request body with the messages and model information.
//  4. Sends a POST request to the chat endpoint from this client.
//  5. Handles the response status code and potential errors.
//  6. Launches a goroutine to process the chat response asynchronously.
//  7. Returns a channel (`<-chan bool`) that signals completion of the chat process and any errors encountered.
//
// The callback function (`cb`) is responsible for handling individual chat responses and errors.
// The completion channel (`<-chan bool`) allows the caller to track the progress of the chat process if needed.
//
// Note that the channel (`chDone`) is not explicitly closed in this example. However, the goroutine
// running `processChat` terminates naturally after sending the completion signal (`true`),
// effectively indicating no more data will be received on the channel.
func (c *Client) Chat(model string, cb ChatCallBack, params *ChatParams, msgs ...ChatMessage) (<-chan bool, error) {
	if cb == nil {
		return nil, ErrCallback
	}

	if len(msgs) == 0 {
		return nil, ErrMessage
	}

	if model == "" {
		model = DEFAULT_MODEL
	}

	request := ChatRequest{
		Model:      model,
		Messages:   msgs,
		ChatParams: params,
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
		switch res.StatusCode {
		case http.StatusBadRequest:
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			return nil, fmt.Errorf("%w\n%v", ErrBadRequest, body)
		default:
			return nil, fmt.Errorf("%w: please make sure ollama server is running and url is correct", ErrInvoke)
		}
	}

	chDone := make(chan bool)

	go func() {
		StreamResponse(res.Body, cb)

		chDone <- true
	}()

	return chDone, nil
}

// Initiates a plain chat process and asynchronously handles responses through a callback function.
//
// This method is identical to Chat(), except that it invokes the callback with plain json string without further processing.
func (c *Client) PlainChat(model string, cb PlainChatCallBack, params *ChatParams, msgs ...ChatMessage) (<-chan bool, error) {
	if cb == nil {
		return nil, ErrCallback
	}

	if len(msgs) == 0 {
		return nil, ErrMessage
	}

	if model == "" {
		model = DEFAULT_MODEL
	}

	request := ChatRequest{
		Model:      model,
		Messages:   msgs,
		ChatParams: params,
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
		switch res.StatusCode {
		case http.StatusBadRequest:
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			return nil, fmt.Errorf("%w\n%v", ErrBadRequest, body)
		default:
			return nil, fmt.Errorf("%w: please make sure ollama server is running and url is correct", ErrInvoke)
		}
	}

	chDone := make(chan bool)

	go func() {
		StreamPlainResponse(res.Body, cb)

		chDone <- true
	}()

	return chDone, nil
}
