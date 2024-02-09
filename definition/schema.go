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
