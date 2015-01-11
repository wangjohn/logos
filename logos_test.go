package logos

import (
  "testing"
  "strings"
)

func TestStringPublicationBody(t *testing.T) {
  fixtures := []struct {
    Body string
    NumLines int
    NumWords int
  }{
    {"This is the body.\nOf the paragraph", 2, 7},
    {"A\nB\nC\nD\nE\nF", 6, 6},
    {"hell0 df .$34\n\n\n", 1, 3},
  }

  for _, fixture := range fixtures {
    input := strings.NewReader(fixture.Body)
    body, err := CreateStringPubBody(input)
    if err != nil {
      t.Errorf("Did not expect error: %v", err)
    }

    lines, words := countLinesWords(body)
    if (lines != fixture.NumLines) {
      t.Errorf("Expected %v lines, but got %v", fixture.NumLines, lines)
    }

    if (words != fixture.NumWords) {
      t.Errorf("Expected %v words, but got %v", fixture.NumWords, words)
    }
  }
}

func countLinesWords(body StringPublicationBody) (int, int) {
  body.ResetSeeker()
  lines := 0
  for (body.HasNextLine()) {
    body.NextLine()
    lines++
  }

  body.ResetSeeker()
  words := 0
  for (body.HasNextWord()) {
    body.NextWord()
    words++
  }

  return lines, words
}
