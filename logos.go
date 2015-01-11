package logos

/*
Logos

This program analyzes articles, identifying important elements in highly ranking articles.
*/

import (
  "bufio"
  "io"
  "strings"
  "unicode"
  "strconv"
)

type Publication struct {
  Score float64
  Author string
  Text PublicationBody
}

type PublicationBody interface {
  NextLine() (string)
  HasNextLine() (bool)
  NextWord() (string)
  HasNextWord() (bool)
  ResetSeeker()
}

type StringPublicationBody struct {
  Lines []string
  Words [][]string
  CurrentLine int
  CurrentWord int
}

func (s StringPublicationBody) HasNextLine() (bool) {
  return len(s.Lines) > s.CurrentLine
}

func (s StringPublicationBody) NextLine() (string) {
  res := s.Lines[s.CurrentLine]
  s.CurrentLine++
  return res
}

func (s StringPublicationBody) HasNextWord() (bool) {
  if (len(s.Lines) <= s.CurrentLine) {
    return false
  }

  return (len(s.Words[s.CurrentLine]) > s.CurrentWord) || (len(s.Lines) > s.CurrentLine + 1)
}

func (s StringPublicationBody) NextWord() (string) {
  word := s.Words[s.CurrentLine][s.CurrentWord]
  if (len(s.Words[s.CurrentLine]) < s.CurrentWord + 1) {
    s.CurrentLine++
    s.CurrentWord = 0
  }
  return word
}

func (s StringPublicationBody) ResetSeeker() {
  s.CurrentLine = 0
  s.CurrentWord = 0
}

func CreateStringPubBody(input io.Reader) (StringPublicationBody, error) {
  scanner := bufio.NewScanner(input)
  lines := make([]string, 0)
  words := make([][]string, 0)

  for scanner.Scan() {
    line := scanner.Text()
    lineWords := splitWords(line)

    if (len(lineWords) > 0) {
      lines = append(lines, line)
      words = append(words, lineWords)
    }
  }

  body := StringPublicationBody{
    Lines: lines,
    Words: words,
    CurrentLine: 0,
    CurrentWord: 0,
  }
  return body, nil
}

/*
---------------------------------------------------------------
-------------- Measures of Publication Quality ----------------
---------------------------------------------------------------
*/

func (p PublicationBody) WordCount() (int) {
  count := 0

  for (p.HasNextWord()) {
    l := p.NextWord()
    count++
  }

  return count
}

func (p PublicationBody) AverageWordsPerLine() (float64) {
  sum := 0
  count := 0

  for (p.HasNextLine()) {
    l := p.NextLine()
    words := splitWords(l)
    sum += len(words)
    count++
  }

  return float64(sum) / float64(count)
}

func (p PublicationBody) AverageWordLength() (float64) {
  sum := 0
  count := 0

  for (p.HasNextWord()) {
    w := p.NextWord()
    sum += len(w)
    count++
  }

  return float64(sum) / float64(count)
}

func (p PublicationBody) WordsLongerThan(x int) (int) {
  count := 0

  for (p.HasNextWord()) {
    w := p.NextWord()
    if (len(w) > x) {
      count++
    }
  }

  return count
}

func (p PublicationBody) WordsIn(list WordList) (int) {
  count := 0

  for (p.HasNextWord()) {
    w := p.NextWord()
    if (list.Contains(w)) {
      count++
    }
  }

  return count
}

func (p PublicationBody) ConstructMarkovMatrix(ngramSize int) (MarkovMatrix) {
  prevNGram := NGram{ngramSize, []string{}}
  matrix := CreateMarkovMatrix()

  for (p.HasNextWord()) {
    w := p.NextWord()
    newWords := append(prevNGram.Words, w)

    if (len(newWords) == ngramSize + 1) {
      ngram := NGram{ngramSize, newWords[1:]}
      curCount := matrix.GetProbability(prevNGram, ngram)

      matrix.SetProbability(prevNGram, ngram, curCount + 1.0)
    }
  }

  res := CreateMarkovMatrix()
  for i := range matrix.Matrix {
    total := 0.0

    for j := range matrix.Matrix[i] {
      total += matrix.Matrix[i][j]
    }

    for j := range matrix.Matrix[i] {
      ig := HashToNGram(i)
      jg := HashToNGram(j)
      res.SetProbability(ig, jg, matrix.GetProbability(ig, jg) / total)
    }
  }

  return res
}

type NGram struct {
  Size int
  Words []string
}

func (n NGram) Hash() (string) {
  s := []string{strconv.Itoa(n.Size)}
  s = append(s, n.Words...)
  return strings.Join(s, ".")
}

func HashToNGram(hash string) (NGram) {
  s := strings.Split(hash, ".")
  size, _ := strconv.Atoi(s[0])
  return NGram{size, s[1:]}
}

type MarkovMatrix struct {
  Matrix map[string]map[string]float64
}

func CreateMarkovMatrix() (MarkovMatrix) {
  m := make(map[string]map[string]float64)
  return MarkovMatrix{m}
}

func (m MarkovMatrix) SetProbability(i, j NGram, prob float64) {
  entry := m.Matrix[i.Hash()]
  if (entry == nil) {
    entry = make(map[string]float64)
  }
  entry[j.Hash()] = prob
  m.Matrix[i.Hash()] = entry
}

func (m MarkovMatrix) GetProbability(i, j NGram) (float64) {
  entry := m.Matrix[i.Hash()]
  if (entry == nil) {
    return 0.0
  } else {
    return entry[j.Hash()]
  }
}

type WordList struct {
  Words map[string]bool
}

func ConstructWordList(words []string) (WordList) {
  list := make(map[string]bool)

  for _, w := range words {
    list[w] = true
  }

  return WordList{list}
}

func (w WordList) Contains(word string) (bool) {
  return w.Words[word]
}

/*
splitWords returns the separate words that make up a particular string, making
sure to remove punctuation and spaces.
*/
func splitWords(line string) ([]string) {
  f := func(c rune) bool {
    return unicode.IsPunct(c) || unicode.IsSpace(c)
  }
  return strings.FieldsFunc(line, f)
}
