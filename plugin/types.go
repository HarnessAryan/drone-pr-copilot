package plugin

// File represents a list of lines in the pull request
type FileDiff struct {
	Name          string
	PreviousLines []Line
	NewLines      []Line
}

type Line struct {
	Number  int
	Content string
}

// Feedback is what we receive from OpenAI and comment back on the PR
type Feedback struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
	Suggestion string `json:"suggestion"`
	Severity   string `json:"severity"` // optional
}
