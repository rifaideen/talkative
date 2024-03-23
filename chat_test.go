package talkative_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rifaideen/talkative"

	"github.com/stretchr/testify/assert"
)

// TestChatResponse tests the chat response handling in the talkative package.
//
// It initializes a mock server to simulate chat responses in NDJSON format and verifies
// if the client correctly concatenates the received messages.
//
// The test involves creating a mock server that sends predefined chat responses, initializing
// a client to connect to this server, and sending a chat message.
//
// The responses are collected and concatenated using a callback function. Finally, the test
// asserts the concatenated response matches the expected output. This ensures the chat
// functionality processes and combines messages accurately.
//
// Parameters:
// - t: A *testing.T object for running assertions.
func TestChatValidation(t *testing.T) {
	message := talkative.ChatMessage{
		Role:    talkative.USER,
		Content: "Hi there!",
	}
	scenario := "not-found"
	server := mockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if scenario == "not-found" {
			w.WriteHeader(404)

			return
		}

		if scenario == "non-json" {
			w.Write([]byte("ok"))

			return
		}

		if scenario == "bad-request" {
			w.WriteHeader(400)
			w.Write([]byte(`{"error": "invalid request"}`))
			return
		}

		// Add more scenarios
	}))

	defer server.Close()

	client, err := talkative.New(server.URL)
	{
		assert.NoError(t, err)
		assert.NotNil(t, client)
	}

	// Assert callback error
	done, err := client.Chat(talkative.DEFAULT_MODEL, nil, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrCallback)
	}

	// Assert message error
	done, err = client.Chat(talkative.DEFAULT_MODEL, func(cr *talkative.ChatResponse, err error) {}, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrMessage)
	}

	// Begin scenario based validation

	scenario = "not-found"
	{
		done, err = client.Chat(talkative.DEFAULT_MODEL, func(cr *talkative.ChatResponse, err error) {}, nil, message)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrInvoke)
	}

	scenario = "non-json"
	{
		done, err = client.Chat(
			talkative.DEFAULT_MODEL,
			func(cr *talkative.ChatResponse, err error) {
				assert.ErrorIs(t, err, talkative.ErrDecoding)
			},
			nil,
			message,
		)

		assert.Nil(t, err)
		assert.NotNil(t, done)

		<-done // wait for completion
	}

	scenario = "bad-request"
	{
		done, err = client.Chat("", func(cr *talkative.ChatResponse, err error) {}, nil, message)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrBadRequest)
	}
}

// TestChatResponse tests the chat response handling in the talkative package.
//
// It initializes a mock server to simulate chat responses in NDJSON format and verifies
// if the client correctly concatenates the received messages.
//
// The test involves creating a mock server that sends predefined chat responses,
// initializing a client to connect to this server, and sending a chat message.
//
// The responses are collected and concatenated using a callback function. Finally, the test asserts the concatenated response matches
// the expected output.
//
// This ensures the chat functionality processes and combines messages accurately.
//
// Parameters:
// - t: A *testing.T object for running assertions.
func TestChatResponse(t *testing.T) {
	server := mockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responses := []talkative.ChatResponse{
			{
				Model: talkative.DEFAULT_MODEL,
				Message: talkative.ChatMessage{
					Role:    talkative.ASSISTANT,
					Content: "Hello",
				},
			},
			{
				Model: talkative.DEFAULT_MODEL,
				Message: talkative.ChatMessage{
					Role:    talkative.ASSISTANT,
					Content: ", ",
				},
			},
			{
				Model: talkative.DEFAULT_MODEL,
				Message: talkative.ChatMessage{
					Role:    talkative.ASSISTANT,
					Content: "It is nice talking to you.",
				},
			},
		}

		w.Header().Add("Content-Type", "application/x-ndjson")
		w.Header().Add("Transfer-Encoding", "chunked")

		flusher, ok := w.(http.Flusher)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Server doesn't support flushing")
			return
		}

		writer := json.NewEncoder(w)

		for _, response := range responses {
			err := writer.Encode(response)

			if err != nil {
				fmt.Println("error encoding response")
				return
			}

			w.Write([]byte("\n"))

			time.Sleep(100 * time.Millisecond)

			flusher.Flush()
		}
	}))

	defer server.Close()

	client, err := talkative.New(server.URL)

	// make sure client is set and no error
	assert.NotNil(t, client)
	assert.NoError(t, err)

	message := talkative.ChatMessage{
		Role:    talkative.USER,
		Content: "Hi there!",
	}

	sb := strings.Builder{}

	done, err := client.Chat(talkative.DEFAULT_MODEL, func(cr *talkative.ChatResponse, err error) {
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			sb.WriteString(cr.Message.Content)
		}
	},
		nil,
		message,
	)

	assert.NotNil(t, done)
	assert.Nil(t, err)

	<-done

	assert.Equal(t, "Hello, It is nice talking to you.", sb.String())
}

// TestPlainChatResponse tests the plain chat response handling in the talkative package.
//
// It initializes a mock server to simulate chat responses in NDJSON format and verifies
// if the client correctly concatenates the received messages.
//
// The test involves creating a mock server that sends predefined chat responses, initializing
// a client to connect to this server, and sending a chat message.
//
// The responses are collected and concatenated using a callback function. Finally, the test
// asserts the concatenated response matches the expected output. This ensures the chat
// functionality processes and combines messages accurately.
//
// Parameters:
// - t: A *testing.T object for running assertions.
func TestPlainChatValidation(t *testing.T) {
	message := talkative.ChatMessage{
		Role:    talkative.USER,
		Content: "Hi there!",
	}
	scenario := "not-found"
	server := mockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if scenario == "not-found" {
			w.WriteHeader(404)

			return
		}

		if scenario == "non-json" {
			w.Write([]byte("ok"))

			return
		}

		if scenario == "bad-request" {
			w.WriteHeader(400)
			w.Write([]byte(`{"error": "invalid request"}`))
			return
		}

		// Add more scenarios
	}))

	defer server.Close()

	client, err := talkative.New(server.URL)
	{
		assert.NoError(t, err)
		assert.NotNil(t, client)
	}

	// Assert callback error
	done, err := client.PlainChat(talkative.DEFAULT_MODEL, nil, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrCallback)
	}

	// Assert message error
	done, err = client.PlainChat(talkative.DEFAULT_MODEL, func(cr string, err error) {}, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrMessage)
	}

	// Begin scenario based validation

	scenario = "not-found"
	{
		done, err = client.PlainChat(talkative.DEFAULT_MODEL, func(cr string, err error) {}, nil, message)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrInvoke)
	}

	scenario = "non-json"
	{
		done, err = client.PlainChat("", func(cr string, err error) {
			assert.ErrorIs(t, err, talkative.ErrDecoding)
		}, nil, message)

		assert.Nil(t, err)
		assert.NotNil(t, done)

		<-done // wait for completion
	}

	scenario = "bad-request"
	{
		done, err = client.PlainChat(talkative.DEFAULT_MODEL, func(cr string, err error) {}, nil, message)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrBadRequest)
	}
}

// TestChatResponse tests the chat response handling in the talkative package.
//
// It initializes a mock server to simulate chat responses in NDJSON format and verifies
// if the client correctly concatenates the received messages.
//
// The test involves creating a mock server that sends predefined chat responses,
// initializing a client to connect to this server, and sending a chat message.
//
// The responses are collected and concatenated using a callback function. Finally, the test asserts the concatenated response matches
// the expected output.
//
// This ensures the chat functionality processes and combines messages accurately.
//
// Parameters:
// - t: A *testing.T object for running assertions.
func TestPlainChatResponse(t *testing.T) {
	server := mockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responses := []talkative.ChatResponse{
			{
				Model: talkative.DEFAULT_MODEL,
				Message: talkative.ChatMessage{
					Role:    talkative.ASSISTANT,
					Content: "Hello",
				},
			},
			{
				Model: talkative.DEFAULT_MODEL,
				Message: talkative.ChatMessage{
					Role:    talkative.ASSISTANT,
					Content: ", ",
				},
			},
			{
				Model: talkative.DEFAULT_MODEL,
				Message: talkative.ChatMessage{
					Role:    talkative.ASSISTANT,
					Content: "It is nice talking to you.",
				},
			},
		}

		w.Header().Add("Content-Type", "application/x-ndjson")
		w.Header().Add("Transfer-Encoding", "chunked")

		flusher, ok := w.(http.Flusher)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Server doesn't support flushing")
			return
		}

		writer := json.NewEncoder(w)

		for _, response := range responses {
			err := writer.Encode(response)

			if err != nil {
				fmt.Println("error encoding response")
				return
			}

			w.Write([]byte("\n"))

			time.Sleep(100 * time.Millisecond)

			flusher.Flush()
		}
	}))

	defer server.Close()

	client, err := talkative.New(server.URL)

	// make sure client is set and no error
	assert.NotNil(t, client)
	assert.NoError(t, err)

	message := talkative.ChatMessage{
		Role:    talkative.USER,
		Content: "Hi there!",
	}

	sb := strings.Builder{}

	done, err := client.PlainChat(talkative.DEFAULT_MODEL, func(cr string, err error) {
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			response := &talkative.ChatResponse{}
			err := json.Unmarshal([]byte(cr), response)

			if err == io.EOF {
				return
			}

			sb.WriteString(response.Message.Content)
		}
	}, nil, message)

	assert.NotNil(t, done)
	assert.Nil(t, err)

	<-done

	assert.Equal(t, "Hello, It is nice talking to you.", sb.String())
}

// mockServer is a helper function that creates a mock HTTP server for testing purposes.
//
// It takes a handler function as a parameter, which will be used to handle incoming HTTP requests.
// The handler function should have the signature `func(http.ResponseWriter, *http.Request)`.
//
// The function returns a pointer to an httptest.Server, which represents the mock server.
// The server can be closed using the `Close` method when it is no longer needed.
//
// Example usage:
//
//	server := mockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	  // handle the request
//	}))
//	defer server.Close()
//
// Parameters:
// - handler: A handler function that will be used to handle incoming HTTP requests.
//
// Returns:
// - A pointer to an httptest.Server representing the mock server.
func mockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}
