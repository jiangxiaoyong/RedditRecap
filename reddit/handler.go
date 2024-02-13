package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"redditRecap/definition"
	"redditRecap/retry"
	"sort"
)

func Search(client *http.Client, query, sortType, limit string) ([]definition.SearchResult, error) {
	baseURL := fmt.Sprintf("https://www.reddit.com/search.json")
	params := url.Values{}
	params.Add("q", query)
	params.Add("sort", sortType)
	params.Add("limit", limit)

	// Construct the URL with query parameters
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Println("search endpoint:", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	resp, err := retry.HttpRetry(client, req)
	if err != nil {
		fmt.Println("Error making search request:", err)
		return nil, errors.New("error making search request")
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Reddit search response %v", resp.Status)
		return nil, errors.New(msg)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding search JSON:", err)
		return nil, errors.New("error decoding search JSON")
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid response format, no data in search response")
		return nil, errors.New("invalid response format, no data in search response")
	}

	children, ok := data["children"].([]interface{})
	if !ok {
		fmt.Println("Invalid response format, no children in search response")
		return nil, errors.New("invalid response format, no children in search response")
	}

	// Decode the JSON response
	var searchResults []definition.SearchResult
	for _, child := range children {
		post, ok := child.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid response format, no post")
			continue
		}

		postData, ok := post["data"].(map[string]interface{})
		if !ok {
			fmt.Println("Invalid response format, no post data")
			continue
		}

		_, ok = postData["title"].(string)
		if !ok {
			fmt.Println("Invalid response format, no title")
			continue
		}

		var searchResult definition.SearchResult
		bytes, _ := json.Marshal(postData)
		err := json.Unmarshal(bytes, &searchResult)
		if err != nil {
			fmt.Println("Failed to decode searchResult")
			return nil, errors.New("failed to decode searchResult")
		}
		searchResults = append(searchResults, searchResult)
	}

	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].UPS > searchResults[j].UPS
	})

	return searchResults, nil
}

func Comments(client *http.Client, subreddit, articleID string) ([]definition.Comment, error) {
	endpoint := fmt.Sprintf("https://www.reddit.com/%s/comments/%s.json", subreddit, articleID)
	fmt.Println("comments endpoint:", endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)
	resp, err := retry.HttpRetry(client, req)
	if err != nil {
		fmt.Println("Error making comments request:", err)
		return nil, errors.New("error making comments request")
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Reddit comments response %v", resp.Status)
		return nil, errors.New(msg)
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var comments []interface{}
	err = json.NewDecoder(resp.Body).Decode(&comments)
	if err != nil {
		fmt.Println("Error decoding comments JSON:", err)
		return nil, errors.New("error decoding comments JSON")
	}

	var commentList []definition.Comment
	for _, comment := range comments[1].(map[string]interface{})["data"].(map[string]interface{})["children"].([]interface{}) {
		commentData := comment.(map[string]interface{})["data"].(map[string]interface{})
		bytes, _ := json.Marshal(commentData)
		var comment definition.Comment
		err := json.Unmarshal(bytes, &comment)
		if err != nil {
			fmt.Println("Failed to decode comment")
			return nil, errors.New("failed to decode comment")
		}

		var replyList []definition.Comment
		if replies, ok := commentData["replies"]; ok && replies != "" {
			replies := commentData["replies"].(map[string]interface{})
			if _, ok := replies["data"]; ok {
				children := replies["data"].(map[string]interface{})["children"].([]interface{})
				for _, child := range children {
					replyData := child.(map[string]interface{})["data"]
					bytes, _ := json.Marshal(replyData)
					var reply definition.Comment
					err := json.Unmarshal(bytes, &reply)
					if err != nil {
						fmt.Println("Failed to decode reply")
						return nil, errors.New("failed to decode reply")
					}
					replyList = append(replyList, reply)
				}
			}
		}
		comment.ReplyList = replyList
		commentList = append(commentList, comment)
	}
	return commentList, nil
}
