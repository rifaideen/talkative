package talkative_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"talkative"
	"testing"
	"time"

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

		// Add more scenarios
	}))

	defer server.Close()

	client, err := talkative.New(server.URL)
	{
		assert.NoError(t, err)
		assert.NotNil(t, client)
	}

	done, err := client.Completion(nil, nil)
	{
		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrCallback)
	}

	// Begin scenario based validation

	scenario = "not-found"
	{
		done, err = client.Completion(func(cr *talkative.CompletionResponse, err error) {}, nil)

		assert.Nil(t, done)
		assert.ErrorIs(t, err, talkative.ErrMessage)

		done, err = client.Completion(func(cr *talkative.CompletionResponse, err error) {}, message)
		{
			assert.Nil(t, done)
			assert.ErrorIs(t, err, talkative.ErrInvoke)
		}

		scenario = "non-json"
		{
			done, err = client.Completion(func(cr *talkative.CompletionResponse, err error) {
				assert.ErrorIs(t, err, talkative.ErrDecoding)
			}, message)

			assert.Nil(t, err)
			assert.NotNil(t, done)

			<-done // wait for completion
		}
	}
}

func TestCompletionResponse(t *testing.T) {
	server := mockServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responses := []talkative.CompletionResponse{
			{
				Model:    "llama2",
				Response: "Hello",
			},
			{
				Model:    "llama2",
				Response: ", ",
			},
			{
				Model:    "llama2",
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

	done, err := client.Completion(func(cr *talkative.CompletionResponse, err error) {
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
