package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"redditRecap/definition"
	"redditRecap/retry"
)

const (
	apiEndpoint = "https://generativelanguage.googleapis.com/v1/models/gemini-pro:generateContent"
)

var client = &http.Client{}

func Inquiry(text string) string {
	geminiEndpoint := apiEndpoint + "?key=" + os.Getenv("GEMINI_API_KEY")
	fmt.Printf("Post to LLM gemini endpoint = %v\n", geminiEndpoint)

	payload := definition.Payload{
		// {"contents":[{"parts":[{"text": "hello"}]}]}
		Contents: []definition.Content{
			{
				Parts: []definition.Part{
					{Text: text},
				},
			},
		},
	}
	reqBody, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", geminiEndpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := retry.HttpRetry(client, req)
	if err != nil {
		fmt.Println("Error inquiry llm:", err)
		return ""
	}
	defer resp.Body.Close()

	// Reading the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading llm response body:", err)
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Response Status NOT OK from llm:", resp.Status)
		return ""
	} else {
		return string(responseBody)
	}
}

type Text struct {
	Text string `json:"text,omitempty"`
}
type Content struct {
	Parts []Text `json:"parts,omitempty"`
}
type Candidate struct {
	Content Content `json:"content,omitempty"`
}
type Response struct {
	Candidates []Candidate `json:"candidates,omitempty"`
}

func ProcessResponse(responseRaw string) (string, error) {
	var response Response
	err := json.Unmarshal([]byte(responseRaw), &response)
	if err != nil {
		fmt.Println("Error unmarshaling llm response JSON:", err)
		return "", err
	}

	if len(response.Candidates) > 0 &&
		len(response.Candidates[0].Content.Parts) > 0 &&
		response.Candidates[0].Content.Parts[0].Text != "" {

		return response.Candidates[0].Content.Parts[0].Text, nil // Text exists
	} else {
		return "", errors.New("empty llm response") // Text does not exist
	}
}
