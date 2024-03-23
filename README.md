[![Go Reference](https://pkg.go.dev/badge/github.com/rifaideen/talkative.svg)](https://pkg.go.dev/github.com/rifaideen/talkative)
[![Go Report Card](https://goreportcard.com/badge/github.com/rifaideen/talkative)](https://goreportcard.com/report/github.com/rifaideen/talkative)

# Talkative

A Golang Ollama REST API Wrapper

The `talkative` package simplifies interaction with the Ollama REST API in Golang applications. It provides a user-friendly interface to access Ollama functionalities without dealing with the raw HTTP requests and responses.

## Features

- **Abstraction:** Hides the underlying HTTP details, allowing you to focus on the Ollama API functionality.
- **Type Safety:** Provides strongly typed methods for interacting with the API.
- **Error Handling:** Handles errors gracefully, returning informative error messages.

* **Flexible Interactions:** Offers two interaction modes for both`/chat` and`/generate` endpoints:
  * **Typed Response:**  Suitable for scenarios where you are utilising the response directly in your program to display or process further.
  * **Untyped Response:** Ideal for scenarios where you are not utilising that response by yourself other than forwarding this to your own APIs i.e: built your own APIs utilizing `talkative` package to interact with Ollama server.

* **Typed Response:** Suitable for scenarios where you are utilizing the response directly in your program to display or process further.
* **Untyped Response:** Ideal
  for forwarding the raw response from Ollama to your own REST API
  endpoints without additional processing. This is useful when you don't
  need to modify the Ollama response within your Golang application

## Installation

```sh
go get github.com/rifaideen/talkative
```

## Usage

```go
    package main

    import (
    "fmt"
  
    "github.com/rifaideen/talkative"
    )

    func main() {
        // Create a new talkative client with Ollama server url
        client, err := talkative.New("http://localhost:11434")
  
        if err != nil {
            panic("Failed to create talkative client")
        }

        // client is ready, start fuelling with your curiosity.
    }
```

## Examples

To explore practical examples of using the `talkative` package for various tasks, navigate to the `_examples` directory within the package.

## Feedback

- [Submit feedback](https://github.com/rifaideen/talkative/issues/new)

## Contributing

We welcome contributions to the `talkative` package! While we don't have a formal CONTRIBUTING.md file yet, feel free to submit pull requests with clear descriptions of your changes. We'll be happy to review them.

## Disclaimer of Non-Liability

This project is provided **"as is"** and **without any express or implied warranties**, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose. In no event shall the authors or copyright holders be liable for any claim, damages or other liability, whether in an action of contract, tort or otherwise, arising from, out of or in connection with the software or the use or other dealings in the software.

## License

The `talkative` package is licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0).

## Additional Resources

[**Ollama API documentation**](https://github.com/ollama/ollama/blob/main/docs/api.md)
