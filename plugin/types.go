package plugin

// File represents a list of lines in the pull request
type File struct {
	Name          string
	PreviousLines []Line
	DiffLines     []Line
}

type Line struct {
	Number  int
	Content string
	Removed bool // whether the line was added or removed, default: false
}

// Feedback is what we receive from OpenAI and comment back on the PR
type Feedback struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
	Suggestion string `json:"suggestion"`
	Severity   string `json:"severity"` // optional
}
