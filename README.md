# RedditRecap

With nested comments and long threads, Reddit discussions can be hard to follow. This project simplifies Reddit discussions by condensing topics and organizing individual opinions, making it easier to understand the conversation.



The project consists of two services that communicate using gRPC.

* Http server
* Reddit query service

### To build the project
* Go to Google AI Studio to get API key, and store the api_key as env varible `GEMINI_API_KEY`
* Run the Protobuf code gen
`protoc --go_out=. --go-grpc_out=. redditquery.proto`

### To run the project

* Go to folder reddit, and `go run .`
* Go to root foler of redditRecap, and `go run .`
* curl `http://localhost:8080/reddit?query=golang`, repladce the query string to the topic that you're interested in.


## Demo
| Topic discussion | Reddit Recap |
| -------- | -------- |
| <img width="200" alt="Screenshot 2024-02-17 at 11 46 52 AM" src="https://github.com/jiangxiaoyong/RedditRecap/assets/5414525/2dc9a540-70a0-45d3-abd3-2ffe7a861d9c">|<img width="200" alt="Screenshot 2024-02-17 at 11 47 17 AM" src="https://github.com/jiangxiaoyong/RedditRecap/assets/5414525/a04cb0ae-8c9d-4f35-b06a-22cfeb86a1c3">|

