package plugin

// FileDiff represents a list of lines in the pull request
type FileDiff struct {
	Name          string `json:"name"`
	PreviousLines []Line `json:"previous_lines"`
	NewLines      []Line `json:"new_lines"`
}

type Line struct {
	Number  int    `json:"number"`
	Content string `json:"content"`
}

// Feedback is what we receive from OpenAI and comment back on the PR
type Feedback struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
	Suggestion string `json:"suggestion"`
	Severity   string `json:"severity"` // optional
}
