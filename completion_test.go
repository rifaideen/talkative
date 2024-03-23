package talkative_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/rifaideen/talkative"

	"github.com/stretchr/testify/assert"
)

func TestCompletionValidation(t *testing.T) {
	message := &talkative.CompletionMessage{
		Prompt: "Hi there!",
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

	done, err := client.Completion("", nil, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrCallback)
	}

	// Assert no message error
	done, err = client.Completion(talkative.DEFAULT_MODEL, func(cr *talkative.CompletionResponse, err error) {}, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrMessage)
	}

	// Begin scenario based validation

	scenario = "not-found"
	{
		done, err = client.Completion(talkative.DEFAULT_MODEL, func(cr *talkative.CompletionResponse, err error) {}, message)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrInvoke)
	}

	scenario = "non-json"
	{
		done, err = client.Completion(talkative.DEFAULT_MODEL, func(cr *talkative.CompletionResponse, err error) {
			assert.ErrorIs(t, err, talkative.ErrDecoding)
		}, message)

		assert.Nil(t, err)
		assert.NotNil(t, done)

		<-done // wait for completion
	}

	scenario = "bad-request"
	{
		done, err = client.Completion("", func(cr *talkative.CompletionResponse, err error) {}, message)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrBadRequest)
	}
}

func TestCompletionResponse(t *testing.T) {
	server := mockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responses := []talkative.CompletionResponse{
			{
				Model:    talkative.DEFAULT_MODEL,
				Response: "Hello",
			},
			{
				Model:    talkative.DEFAULT_MODEL,
				Response: ", ",
			},
			{
				Model:    talkative.DEFAULT_MODEL,
				Response: "It is nice talking to you.",
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

	message := &talkative.CompletionMessage{
		Prompt: "Hi there!",
	}

	sb := strings.Builder{}

	done, err := client.Completion("", func(cr *talkative.CompletionResponse, err error) {
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			sb.WriteString(cr.Response)
		}
	}, message)

	assert.NotNil(t, done)
	assert.Nil(t, err)

	<-done

	assert.Equal(t, "Hello, It is nice talking to you.", sb.String())
}

func TestPlainCompletionValidation(t *testing.T) {
	message := &talkative.CompletionMessage{
		Prompt: "Hi there!",
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

	done, err := client.PlainCompletion("", nil, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrCallback)
	}

	// Assert no message error
	done, err = client.PlainCompletion(talkative.DEFAULT_MODEL, func(cr string, err error) {}, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrMessage)
	}

	// Begin scenario based validation

	scenario = "not-found"
	{
		done, err = client.PlainCompletion(talkative.DEFAULT_MODEL, func(cr string, err error) {}, message)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrInvoke)
	}

	scenario = "non-json"
	{
		done, err = client.PlainCompletion(talkative.DEFAULT_MODEL, func(cr string, err error) {
			assert.ErrorIs(t, err, talkative.ErrDecoding)
		}, message)

		assert.Nil(t, err)
		assert.NotNil(t, done)

		<-done // wait for completion
	}

	scenario = "bad-request"
	{
		done, err = client.PlainCompletion("", func(cr string, err error) {}, message)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrBadRequest)
	}
}

func TestPlainCompletionResponse(t *testing.T) {
	server := mockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responses := []talkative.CompletionResponse{
			{
				Model:    talkative.DEFAULT_MODEL,
				Response: "Hello",
			},
			{
				Model:    talkative.DEFAULT_MODEL,
				Response: ", ",
			},
			{
				Model:    talkative.DEFAULT_MODEL,
				Response: "It is nice talking to you.",
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

	message := &talkative.CompletionMessage{
		Prompt: "Hi there!",
	}

	sb := strings.Builder{}

	done, err := client.PlainCompletion("", func(cr string, err error) {
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			var response *talkative.CompletionResponse
			json.Unmarshal([]byte(cr), &response)

			if response != nil {
				sb.WriteString(response.Response)
			}
		}
	}, message)

	assert.NotNil(t, done)
	assert.Nil(t, err)

	<-done

	assert.Equal(t, "Hello, It is nice talking to you.", sb.String())
}
