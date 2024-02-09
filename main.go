package main

import (
	"fmt"
	"net/http"
	"redditRecap/llm"
	"redditRecap/reddit"
)

func main() {

	http.HandleFunc("/reddit", func(w http.ResponseWriter, r *http.Request) {
		// Parse the query parameters from the request URL
		queryParams := r.URL.Query()

		// Extract and print specific query parameters
		query := queryParams.Get("query")
		fmt.Println("user query:", query)

		searchList, err := reddit.Search(query, query, "relevance", "5")
		if err != nil || searchList == nil || len(searchList) == 0 {
			fmt.Fprintf(w, "please try again %+v", err.Error())
			return
		}
		topic := searchList[0]

		comments, err := reddit.Comments(topic.SubredditNamePrefixed, topic.ID)
		if err != nil || comments == nil {
			fmt.Fprintf(w, "please try again %+v", err.Error())
			return
		}

		prompt := llm.Prompt(topic, comments)
		res := llm.Inquiry(prompt)
		text := llm.ProcessResponse(res)
		fmt.Println("\n\n", text)

		fmt.Fprintf(w, text)
	})

	// Start the HTTP server and listen on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
