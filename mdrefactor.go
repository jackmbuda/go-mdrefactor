package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Configuration constants
const (
	// Replace with the actual OpenAI API endpoint for chat completions
	openaiAPIURL = "https://api.openai.com/v1/chat/completions"
	// Default model to use. You can change this to gpt-4, etc.
	defaultModel = "gpt-3.5-turbo"
	// Default system prompt for the AI
	defaultSystemPrompt = "You are a helpful assistant that refactors Markdown content. Please improve its structure, clarity, and formatting while preserving the original meaning."
)

// APIRequest represents the request payload for the OpenAI API
type APIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"` // Set to false for simple refactoring
}

// Message represents a single message in the chat completion request
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// APIResponse represents the expected response structure from the OpenAI API
type APIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *APIError `json:"error,omitempty"`
}

// Choice represents one of the completion choices from the API
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// APIError represents an error returned by the API
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

// Global HTTP client for reuse
var httpClient = &http.Client{Timeout: 60 * time.Second}

// refactorMarkdown sends the markdown content to the OpenAI API for refactoring
func refactorMarkdown(apiKey, model, systemPrompt, markdownContent string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("OpenAI API key is not set. Please set the OPENAI_API_KEY environment variable or use the -apikey flag")
	}

	// Construct the messages for the API request
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf("Refactor the following Markdown content:\n\n%s", markdownContent)},
	}

	// Create the request payload
	apiRequest := APIRequest{
		Model:    model,
		Messages: messages,
		Stream:   false, // We want the full response, not a stream
	}

	// Marshal the request payload to JSON
	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal API request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", openaiAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	fmt.Println("Sending content to API for refactoring...")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read API response body: %w", err)
	}

	// Unmarshal the API response
	var apiResponse APIResponse
	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		// Try to print the raw response body if JSON unmarshalling fails for debugging
		fmt.Fprintf(os.Stderr, "Raw API response: %s\n", string(responseBody))
		return "", fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	// Check for API errors
	if apiResponse.Error != nil {
		return "", fmt.Errorf("API error: %s (Type: %s, Code: %s)", apiResponse.Error.Message, apiResponse.Error.Type, apiResponse.Error.Code)
	}

	// Check if choices are available
	if len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("no refactored content received from API. Raw response: %s", string(responseBody))
	}

	// Extract the refactored content
	refactoredContent := apiResponse.Choices[0].Message.Content
	fmt.Println("Refactoring successful.")
	return refactoredContent, nil
}

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the input Markdown file (required)")
	outputFile := flag.String("output", "", "Path to the output Markdown file (optional, prints to stdout if not provided)")
	apiKey := flag.String("apikey", os.Getenv("OPENAI_API_KEY"), "OpenAI API key (can also be set via OPENAI_API_KEY environment variable)")
	model := flag.String("model", defaultModel, "OpenAI model to use (e.g., gpt-3.5-turbo, gpt-4)")
	systemPrompt := flag.String("prompt", defaultSystemPrompt, "System prompt to guide the AI refactoring")

	flag.Parse()

	// Validate input file
	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: Input file path is required.")
		flag.Usage()
		os.Exit(1)
	}

	// Read the input Markdown file
	markdownBytes, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file %s: %v\n", *inputFile, err)
		os.Exit(1)
	}
	markdownContent := string(markdownBytes)

	// Check if API key is provided
	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: OpenAI API key is missing. Please provide it using the -apikey flag or set the OPENAI_API_KEY environment variable.")
		os.Exit(1)
	}

	// Call the refactoring function
	refactoredContent, err := refactorMarkdown(*apiKey, *model, *systemPrompt, markdownContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error refactoring Markdown: %v\n", err)
		os.Exit(1)
	}

	// Output the refactored content
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(refactoredContent), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file %s: %v\n", *outputFile, err)
			os.Exit(1)
		}
		fmt.Printf("Refactored content successfully written to %s\n", *outputFile)
	} else {
		// Print to stdout if no output file is specified
		fmt.Println("\n--- Refactored Markdown ---")
		fmt.Println(refactoredContent)
	}
}
