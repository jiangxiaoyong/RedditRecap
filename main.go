package main

import (
	"fmt"
	"net/http"
	"redditRecap/llm"
	"redditRecap/reddit"
)

func main() {
	client := &http.Client{}

	http.HandleFunc("/reddit", func(w http.ResponseWriter, r *http.Request) {
		// Parse the query parameters from the request URL
		queryParams := r.URL.Query()

		// Extract and print specific query parameters
		query := queryParams.Get("query")
		fmt.Println("user query:", query)

		searchList, err := reddit.Search(client, query, query, "relevance", "5")
		if err != nil {
			fmt.Fprintf(w, "please try again %+v", err.Error())
			return
		}
		if searchList == nil || len(searchList) == 0 {
			fmt.Fprintf(w, fmt.Sprintf("please try again, empty search result for %v", query))
			return
		}
		topic := searchList[0]

		comments, err := reddit.Comments(client, topic.SubredditNamePrefixed, topic.ID)
		if err != nil {
			fmt.Fprintf(w, "please try again %+v", err.Error())
			return
		}
		if comments == nil {
			fmt.Fprintf(w, "please try again, empty comments")
			return
		}

		prompt := llm.Prompt(topic, comments)
		res := llm.Inquiry(client, prompt)
		text, err := llm.ProcessResponse(res)
		fmt.Println("\n\n", text)

		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("please try again %v", err.Error()))
		} else {
			fmt.Fprintf(w, text)
		}
	})

	// Start the HTTP server and listen on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
