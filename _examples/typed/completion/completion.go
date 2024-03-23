package main

import (
	"fmt"
	"talkative"
)

func main() {
	client, err := talkative.New("http://localhost:11434")

	if err != nil {
		panic(err)
	}

	// Name of the model to use
	model := talkative.DEFAULT_MODEL
	// Callback function to handle the response
	callback := func(cr *talkative.CompletionResponse, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print(cr.Response)
	}
	// The chat message to send
	message := &talkative.CompletionMessage{
		Prompt: "Why is sky blue?",
		CompletionParams: &talkative.CompletionParams{
			System: "You are Mario from Super Mario Bros.",
		},
	}

	done, err := client.Completion(model, callback, message)

	if err != nil {
		panic(err)
	}

	<-done // wait for the request to complete
	fmt.Println()
}
