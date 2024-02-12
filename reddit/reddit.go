package reddit

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	redditquery "redditRecap/protos/codegen"
)

type server struct {
	redditquery.UnimplementedRedditQueryServiceServer
}

func (s *server) GetRedditTopicSummary(ctx context.Context, in *redditquery.QueryRequest) (*redditquery.QueryResponse, error) {
	fmt.Printf("Received query: %s\n", in.GetQuery())
	// Mock processing and return a summary
	return &redditquery.QueryResponse{Summary: "This is a mock summary of the discussion."}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	redditquery.RegisterRedditQueryServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
