/*
Logos

This program analyzes articles, identifying important elements in highly ranking articles.
*/

import (
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
  io.Reader
  io.Seeker
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

func (p PublicationBody) AverageWordsPerLine() (int) {
  sum := 0
  count := 0

  for (p.HasNextLine()) {
    l := p.NextLine()
    words := splitWords(l)
    sum += len(words)
    count++
  }

  return float64(sum) / count
}

func (p PublicationBody) AverageWordLength() (int) {
  sum := 0
  count := 0

  for (p.HasNextWord()) {
    w := p.NextWord()
    sum += len(w)
    count++
  }

  return float64(sum) / count
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
  prevWords := make([]string, ngramSize)

  for (p.HasNextWord()) {
    w := p.NextWord()
    prevWords = append(prevWords, w)
  }
}

type NGram struct {
  Size int
  Words []string
}

func (n NGram) Hash() (string) {
  s := []string{strconv.Itoa(n.Size)}
  s := append(s, n.Words...)
  return strings.Join(s, ".")
}

type MarkovMatrix struct {
  Matrix map[string]map[string]float64
}

func (m MarkovMatrix) SetProbability(i, j NGram, prob float64) {
  entry := m.Matrx[i.Hash()]
  if (entry == nil) {
    entry = make(map[string]float64)
  }
  entry[j.Hash()] = prob
  m.Matrix[i.Hash()] = entry
}

type WordList struct {
  Words map[string]bool
}

func ConstructWordList(words []string) (WordList) {
  list := make(map[string]bool)

  for _, w := range words {
    list[w] = true
  }

  return list
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
  return strings.FieldsFunc(sentence, f)
}
