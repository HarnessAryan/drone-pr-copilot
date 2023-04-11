package plugin

import (
	"bufio"
	"context"
	"os"
	"testing"
)

func Test(t *testing.T) {
	t.Log("here")
	token := "YOUR_TOKEN_GOES_HERE"
	opt := WithToken(token)
	c := New(opt)
	diff, err := createDiff()
	if err != nil {
		t.Error(err)
	}
	fd := []*FileDiff{diff}
	resp := c.Feedback(context.Background(), fd)
	for _, f := range resp {
		t.Logf("feedback: %+v\n", f)
	}
}

func createDiff() (*FileDiff, error) {
	prevLines, err := parseFile("before.txt")
	if err != nil {
		return nil, err
	}
	afterLines, err := parseFile("after.txt")
	if err != nil {
		return nil, err
	}
	return &FileDiff{
		Name:          "a.txt",
		PreviousLines: prevLines,
		NewLines:      afterLines,
	}, nil
}

func parseFile(filename string) ([]Line, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []Line
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, Line{Number: lineNumber, Content: line})
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
