package talkative

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CompletionRequest represents a request for completion.
type CompletionRequest struct {
	Model  string   `json:"model"`  // The model to use for completion.
	Prompt string   `json:"prompt"` // The prompt for completion.
	Images []string `json:"images"` // The images associated with the completion.
}

// CompletionMessage represents the message structure for initiating a completion request.
type CompletionMessage struct {
	Prompt string   `json:"prompt"` // The text prompt to be completed.
	Images []string `json:"images"` // A list of image URLs associated with the prompt.
}

// CompletionResponse represents the response received after a completion request.
//
// It also embeds CompletionMetrics which includes upon completion
type CompletionResponse struct {
	Model     string `json:"model"`      // The model used for the completion.
	Response  string `json:"response"`   // The generated response based on the prompt.
	CreatedAt string `json:"created_at"` // The timestamp when the response was created.
	Done      bool   `json:"done"`       // A boolean indicating if the completion process is finished.

	CompletionMetrics // embeds CompletionMetrics
}

// CompletionMetrics struct encapsulates various metrics related to the completion process.
// It includes total processing time, model loading time, counts and durations of prompt and overall evaluations,
// and the context encoding of the conversation used in the response.
type CompletionMetrics struct {
	TotalDuration      int   `json:"total_duration"`       // Total processing time in milliseconds.
	LoadDuration       int   `json:"load_duration"`        // Time spent loading the model (milliseconds).
	PromptEvalCount    int   `json:"prompt_eval_count"`    // Number of prompt evaluations performed.
	PromptEvalDuration int   `json:"prompt_eval_duration"` // Time spent on prompt evaluation (milliseconds).
	EvalCount          int   `json:"eval_count"`           // Number of overall evaluations performed.
	EvalDuration       int   `json:"eval_duration"`        // Time spent on overall evaluation (milliseconds).
	Context            []int `json:"context"`              // Encoding of the conversation used in this response.
}

// CompletionCallback defines a function type that is used as a callback for handling completion responses.
// It takes a pointer to a CompletionResponse and an error as arguments.
//
// Parameters:
// - *CompletionResponse: A pointer to the CompletionResponse received after a completion request.
// - error: An error that might have occurred during the completion process.
type CompletionCallback func(*CompletionResponse, error)

// PlainCompletionCallback defines a function type that is used as a callback for handling plain completion responses.
// It takes string and an error as arguments.
//
// Parameters:
// - string: Json encoded plain string received after a completion request.
// - error: An error that might have occurred during the completion process.
type PlainCompletionCallback func(string, error)

// Completion initiates a completion request to the server and returns a channel that signals when the operation is done.
//
// This method takes a CompletionCallback function and a CompletionMessage as arguments. The callback function is invoked
// with the completion response and any error that occurred during the request. The CompletionMessage contains the prompt
// and any associated images for the completion request.
//
// Parameters:
// - cb CompletionCallback: The callback function to be called upon completion of the request. It must not be nil.
// - msg *CompletionMessage: A pointer to the CompletionMessage containing the prompt and images for the completion. It must not be nil.
//
// Returns:
// - <-chan bool: A channel that signals when the completion operation is done. The channel receives a boolean value of true upon completion.
// - error: An error if the callback function or the message is nil, if there's an error encoding the request, or if the server responds with an error.
//
// The method constructs a CompletionRequest from the provided message, encodes it into JSON, and sends it to the server.
// It handles HTTP response status codes, specifically checking for a BadRequest (400) to return any server-side error messages.
// Upon a successful request, it starts a goroutine to stream the response and invoke the provided callback function, signaling completion through the returned channel.
func (c *Client) Completion(model string, cb CompletionCallback, msg *CompletionMessage) (<-chan bool, error) {
	if cb == nil {
		return nil, ErrCallback
	}

	if msg == nil {
		return nil, ErrMessage
	}

	if model == "" {
		model = DEFAULT_MODEL
	}

	request := CompletionRequest{
		Model:  model,
		Prompt: msg.Prompt,
		Images: msg.Images,
	}
	body := &bytes.Buffer{}

	if err := json.NewEncoder(body).Encode(request); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrEncoding, err)
	}

	res, err := c.client.Post(c.urls["completion"], "application/json", body)

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

	chDone := make(chan bool, 1)

	go func() {
		StreamResponse(res.Body, cb)

		chDone <- true
	}()

	return chDone, nil
}

// Completion initiates a plain completion request to the server and returns a channel that signals when the operation is done.
//
// This method is identical to Completion(), except that it invokes the callback with plain json string without further processing.
func (c *Client) PlainCompletion(model string, cb PlainCompletionCallback, msg *CompletionMessage) (<-chan bool, error) {
	if cb == nil {
		return nil, ErrCallback
	}

	if msg == nil {
		return nil, ErrMessage
	}

	if model == "" {
		model = DEFAULT_MODEL
	}

	request := CompletionRequest{
		Model:  model,
		Prompt: msg.Prompt,
		Images: msg.Images,
	}
	body := &bytes.Buffer{}

	if err := json.NewEncoder(body).Encode(request); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrEncoding, err)
	}

	res, err := c.client.Post(c.urls["completion"], "application/json", body)

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

	chDone := make(chan bool, 1)

	go func() {
		StreamPlainResponse(res.Body, cb)

		chDone <- true
	}()

	return chDone, nil
}
