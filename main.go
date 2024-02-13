package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"redditRecap/definition"
	"redditRecap/llm"
	redditquery "redditRecap/protos/codegen"
)

func convertProtoCommentsToGo(pbComments []*redditquery.Comment) []definition.Comment {
	var comments []definition.Comment
	for _, pbComment := range pbComments {
		comment := definition.Comment{
			Author:    pbComment.GetAuthor(),
			ID:        pbComment.GetId(),
			Name:      pbComment.GetName(),
			Body:      pbComment.GetBody(),
			ParentID:  pbComment.GetParentId(),
			ReplyList: convertProtoCommentsToGo(pbComment.ReplyList), // Recursively convert replies
		}
		comments = append(comments, comment)
	}
	return comments
}

// convertProtoSearchResultToGo converts a protobuf SearchResult message to its Go struct equivalent.
func convertProtoSearchResultToGo(protoSR *redditquery.SearchResult) *definition.SearchResult {
	if protoSR == nil {
		return nil
	}
	return &definition.SearchResult{
		Title:                 protoSR.GetTitle(),
		ID:                    protoSR.GetId(),
		Name:                  protoSR.GetName(),
		Author:                protoSR.GetAuthor(),
		SelfText:              protoSR.GetSelfText(),
		SubredditNamePrefixed: protoSR.GetSubredditNamePrefixed(),
		UPS:                   int(protoSR.GetUps()), // Note: Convert int32 to int as required by the Go struct
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	// Establish a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	client := redditquery.NewMyServiceClient(conn)

	// Perform a search request
	searchResProto, err := client.Search(context.Background(), &redditquery.SearchRequest{Query: query})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	searchResGo := convertProtoSearchResultToGo(searchResProto.Result)
	fmt.Printf("Search Result: %v\n", searchResGo)

	// Perform a request to get commentsGO
	commentsProto, err := client.GetComments(
		context.Background(),
		&redditquery.CommentRequest{SubredditName: searchResGo.SubredditNamePrefixed, TopicId: searchResGo.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	commentsGO := convertProtoCommentsToGo(commentsProto.GetComments())
	fmt.Printf("comments: %v\n", commentsGO)

	prompt := llm.Prompt(*searchResGo, commentsGO)
	res := llm.Inquiry(prompt)
	text, err := llm.ProcessResponse(res)
	fmt.Println("LLM response:\n\n", text)

	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("please try again %v", err.Error()))
	} else {
		fmt.Fprintf(w, text)
	}
}

func main() {
	http.HandleFunc("/reddit", handler)
	// Start the HTTP server and listen on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
