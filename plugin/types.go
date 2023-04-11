package plugin

// File represents a list of lines in the pull request
type File struct {
	Name  string
	Lines []Line
}

type Line struct {
	Number  int
	Content string
}

// Feedback is what we receive from OpenAI and comment back on the PR
type Feedback struct {
	Filename   string
	LineNumber int
	Suggestion string
	Severity   string // optional
}
