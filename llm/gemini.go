package llm

import (
	"bytes"
	"encoding/json"
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

func Inquiry(client *http.Client, text string) string {
	geminiEndpoint := apiEndpoint + "?key=" + os.Getenv("BARD_API_KEY")
	fmt.Printf("gemini endpoint = %v\n", geminiEndpoint)

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
	//fmt.Printf("json = %v\n", string(reqBody))

	req, err := http.NewRequest("POST", geminiEndpoint, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := retry.HttpRetry(client, req)
	if err != nil {
		fmt.Println("Error making HTTP POST request:", err)
		return ""
	}
	defer resp.Body.Close()

	// Reading the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Response Status NOT OK:", resp.Status)
		return ""
	} else {
		//fmt.Println("Response Body:", string(responseBody))
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

func ProcessResponse(responseRaw string) string {
	var response Response
	err := json.Unmarshal([]byte(responseRaw), &response)
	if err != nil {
		fmt.Println("Error unmarshaling response JSON:", err)
		return ""
	}

	return response.Candidates[0].Content.Parts[0].Text
}
