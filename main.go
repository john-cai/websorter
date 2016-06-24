package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
)

var validChar = regexp.MustCompile(`^[A-Za-z]+$`)

type payload struct {
	Words   []string `json:"words,omitempty"`
	Result  []string `json:"result"`
	Reverse bool     `json:"reverse"`
}

func sortArray(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// check HTTP verb is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// check content type
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// try to unmarshal the payload
	var p payload

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check if all the words match [A-Za-z]. Ideally I would have liked to use
	// a reader(commented out below) so we wouldn't have to go through these characters twice since we're
	// already reading them once. However it seems that the standard library's json decoder
	// will swallow errors from the reader

	if !validWords(p.Words) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("words must match [A-Za-z]"))
		return
	}

	reverse := p.Reverse

	sortableWords := newSortableWords(p.Words, reverse)
	sort.Sort(sortableWords)

	p.Words = nil
	p.Result = sortableWords.words

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)

	return
}

func validWords(s []string) bool {
	for _, w := range s {
		if !validChar.MatchString(w) {
			return false
		}
	}
	return true
}

type sortableWords struct {
	words   []string
	reverse bool
}

func (s *sortableWords) Len() int {
	return len(s.words)
}

func (s *sortableWords) Swap(i, j int) {
	s.words[i], s.words[j] = s.words[j], s.words[i]
}

func (s *sortableWords) Less(i, j int) bool {
	less := s.words[i] < s.words[j]
	if s.reverse {
		return !less
	}
	return less
}

func newSortableWords(words []string, reverse bool) *sortableWords {
	return &sortableWords{reverse: reverse, words: words}
}

/*
type invalidCharacterError struct {
	Err string
}

func (i *invalidCharacterError) Error() string {
	return i.Err
}

func newInvalidCharacterError(s string) *invalidCharacterError {
	return &invalidCharacterError{
		Err: fmt.Sprintf("Invalid character found in payload %s", s),
	}
}

type validationReader struct {
	r io.ReadCloser
}

func (v *validationReader) Read(p []byte) (int, error) {
	n, err := v.r.Read(p)

	if err != nil && err != io.EOF {
		return n, err
	}

	for _, b := range p[:n] {
		if !validChar.MatchString(string(b)) {
			return n, newInvalidCharacterError(string(b))
		}
	}

	return n, err
}
func newValidationReader(r io.ReadCloser) *validationReader {
	return &validationReader{
		r: r,
	}
}
*/

func main() {
	http.HandleFunc("/sort", sortArray)
	http.ListenAndServe(":8080", nil)
}
