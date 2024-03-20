package talkative

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Streaming the response from the server asynchronously.
//
// This function takes an io.ReadCloser object (`body`) representing the response body
// and a callback function (`cb`) for handling individual responses and errors.
// It iterates through the response, decoding each message and invoking the callback for processing.
//
// In case of errors during decoding or processing, the callback is invoked with the error
// and processing stops. The function closes the response body before exiting.
func StreamResponse[T any](body io.ReadCloser, cb func(*T, error)) {
	defer body.Close()

	for {
		var response T

		err := json.NewDecoder(body).Decode(&response)

		if err == io.EOF {
			return
		}

		if err != nil {
			cb(nil, fmt.Errorf("%w: %v", ErrDecoding, err))

			return
		}

		cb(&response, nil)
	}
}

// Streaming the plain response from the server asynchronously.
//
// This function takes an io.ReadCloser object (`body`) representing the response body
// and a callback function (`cb`) for handling individual responses and errors.
// It iterates through the response and invoking the callback with plain string for processing.
//
// In case of errors during decoding or processing, the callback is invoked with the error
// and processing stops. The function closes the response body before exiting.
func StreamPlainResponse(body io.ReadCloser, cb func(string, error)) {
	defer body.Close()
	buff := bufio.NewReader(body)

	for {
		data, err := buff.ReadString('\n')

		if err == io.EOF {
			return
		}

		if err != nil {
			cb("", err)
			return
		}

		cb(data, nil)
	}
}
