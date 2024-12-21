package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MichaelMure/go-term-markdown"
	"golang.org/x/term"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiRequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ApiResponseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func llm(config Config, prompt string, model string) string {
	// create request body
	requestBody := ApiRequestBody{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// convert request body to json
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	// create request
	req, err := http.NewRequest("POST", config.ApiBaseUrl, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		panic(err)
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.ApiKey))

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// check if status is 200
	if resp.StatusCode != 200 {
		panic("status is not 200")
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var apiResponseBody ApiResponseBody
	if err := json.Unmarshal(responseBody, &apiResponseBody); err != nil {
		panic(err)
	}

	return apiResponseBody.Choices[0].Message.Content
}

func main() {
	config, err := LoadConfig("config.toml")
	if err != nil {
		panic(err)
	}

	prompt := flag.String("p", "Hello there", "Use this flag to pass the prompt")
	model := flag.String("m", config.DefaultModel, "Use this flag to pass the model")

	flag.Parse()

	// check if an alias exists for the current model name else use the literal name
	if aliasValue, exists := (*config).ModelAlias[*model]; exists {
		*model = aliasValue
	}
	response := llm(*config, *prompt, *model)

	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		fmt.Println("Couldn't get terminal width, using default width 80")
		terminalWidth = 80
	} else {
		fmt.Println("terminalWidth: ", terminalWidth)
	}

	printResult := markdown.Render(response, terminalWidth, 6)
	fmt.Println(string(printResult))
}
