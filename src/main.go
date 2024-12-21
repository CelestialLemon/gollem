package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		Model: "qwen/qwen-2.5-coder-32b-instruct",
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
	}
	response := llm(*config, *prompt, *model)

	response := llm(*config, prompt)
	fmt.Println(response)
}
