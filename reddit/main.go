package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"redditRecap/definition"
	redditquery "redditRecap/protos/codegen"
)

var client = &http.Client{}

type myServiceServer struct {
	redditquery.UnimplementedMyServiceServer
}

func (s *myServiceServer) Search(ctx context.Context, req *redditquery.SearchRequest) (*redditquery.SearchResponse, error) {
	query := req.GetQuery()
	fmt.Printf("Search Request Received: %s\n", req.Query)

	searchList, err := Search(client, query, "relevance", "5")
	if err != nil {
		return &redditquery.SearchResponse{Result: nil}, err
	}
	if searchList == nil || len(searchList) == 0 {
		return &redditquery.SearchResponse{Result: nil},
			errors.New(fmt.Sprintf("please try again, empty search result for %v", query))
	}

	topic := searchList[0]
	fmt.Println("reddit search done")
	return &redditquery.SearchResponse{Result: convertSearchResultToProto(&topic)}, nil
}

func (s *myServiceServer) GetComments(ctx context.Context, req *redditquery.CommentRequest) (*redditquery.CommentResponse, error) {
	topicId := req.GetTopicId()
	subredditName := req.GetSubredditName()
	fmt.Printf("GetComments Request Received for Article ID: %s, subreddit: %s\n", topicId, subredditName)

	comments, err := Comments(client, subredditName, topicId)
	if err != nil {
		return &redditquery.CommentResponse{Comments: nil}, err
	}
	if comments == nil {
		return &redditquery.CommentResponse{Comments: nil},
			errors.New(fmt.Sprintf("please try again, empty search result"))
	}

	// Convert Comment structs to protobuf Comment messages
	redditComments := convertCommentsToProto(comments)

	fmt.Println("reddit get comments done")
	return &redditquery.CommentResponse{Comments: redditComments}, nil
}

func convertCommentsToProto(comments []definition.Comment) []*redditquery.Comment {
	var pbComments []*redditquery.Comment
	for _, c := range comments {
		pbComment := &redditquery.Comment{
			Author:    c.Author,
			Id:        c.ID,
			Name:      c.Name,
			Body:      c.Body,
			ParentId:  c.ParentID,
			ReplyList: convertCommentsToProto(c.ReplyList), // Recursively convert replies
		}
		pbComments = append(pbComments, pbComment)
	}
	return pbComments
}

func convertSearchResultToProto(sr *definition.SearchResult) *redditquery.SearchResult {
	if sr == nil {
		return nil
	}
	return &redditquery.SearchResult{
		Title:                 sr.Title,
		Id:                    sr.ID,
		Name:                  sr.Name,
		Author:                sr.Author,
		SelfText:              sr.SelfText,
		SubredditNamePrefixed: sr.SubredditNamePrefixed,
		Ups:                   int32(sr.UPS), // Note: Convert int to int32 as required by protobuf
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	redditquery.RegisterMyServiceServer(s, &myServiceServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
