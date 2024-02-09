package llm

import (
	"fmt"
	"redditRecap/reddit"
	"strings"
)

func Prompt(topic reddit.SearchResult, comments []reddit.Comment) string {
	context := `
I will copy past a topic discussion among a group of people.
The discussion is structured similarly to a Reddit thread, 
featuring nested comments and connections between comments through 
two fields: comment_id and reply_to. 
The comment_id field represents the unique identifier for each comment, 
while the reply_to field indicates which comment a user is responding to.
Try to give general conclusion for the whole discussion and include the Topic name in the conclusion section.
Try to summarize the opinion of each people and list each user's opinion with their user_name.


Please summarize the discussion in the following format:

Topic: [the topic]
Conclusion: [the summarized discussion]

user_name opinion: [the user's opinion]

==============Please learn following example============
-------------Example input-----------------
Topic: go language
user: John
comment_id: 1
comment_body: I enjoy coding in go language

user: Mike
reply_to: 1
comment_id: 2
comment_body: are you sure, I don’t like it

------------Example output----------------
Topic: go language
Conclusion: The opinions on the go programming language are mixed.

John opinion: John enjoys coding in go and finds it to be a productive and enjoyable experience.

Mike opinion: Mike does not like go and finds it less productive or enjoyable to use.
==============End of example============
...
`

	topicSection := fmt.Sprintf("Topic: %v\nuser_name: %v\ncomment_id: %v\ncomment_body: %v\n\n",
		topic.Title, topic.Author, topic.Name, topic.SelfText)

	var posts strings.Builder
	for _, comment := range comments {
		msg := buildMessage(comment)
		posts.WriteString(msg)

		for _, reply := range comment.ReplyList {
			msg = buildMessage(reply)
			posts.WriteString(msg)
		}
	}

	res := context + topicSection + posts.String()
	return res
}

func buildMessage(comment reddit.Comment) string {
	if comment.Author != "" &&
		comment.Body != "" &&
		!strings.Contains(strings.ToLower(comment.Author), "delete") &&
		!strings.Contains(strings.ToLower(comment.Body), "delete") {
		return fmt.Sprintf("user_name: %v\nreply_to: %v\ncomment_id: %v\ncomment_body: %v\n\n",
			comment.Author, comment.ParentID, comment.Name, comment.Body)
	} else {
		return ""
	}
}
