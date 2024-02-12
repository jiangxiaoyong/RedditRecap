package definition

type Payload struct {
	Contents []Content `json:"contents,omitempty"`
}

type Content struct {
	Parts []Part `json:"parts,omitempty"`
}

type Part struct {
	Text string `json:"text,omitempty"`
}

type Comment struct {
	Author    string `json:"author,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Body      string `json:"body,omitempty"`
	ParentID  string `json:"parent_id,omitempty"`
	ReplyList []Comment
}

type SearchResult struct {
	Title                 string `json:"title,omitempty"`
	ID                    string `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	Author                string `json:"author,omitempty"`
	SelfText              string `json:"selftext,omitempty"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed,omitempty"`
	UPS                   int    `json:"ups,omitempty"`
}
