package sorter

import (
	"bytes"
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
)

func TestValidWords(t *testing.T) {
	tests := []struct {
		Words    []string
		Expected bool
	}{
		{
			Words:    []string{"b", "c", "a"},
			Expected: true,
		},
		{
			Words:    []string{""},
			Expected: false,
		},
		{
			Words:    []string{},
			Expected: true,
		},
		{
			Words:    []string{"aa9", "b"},
			Expected: false,
		},
		{
			Words:    []string{"aaaaa*aaa"},
			Expected: false,
		},
	}

	for _, test := range tests {
		valid := validWords(test.Words)
		if valid != test.Expected {
			t.Errorf("With %v, expected %t, got %t", test.Expected, valid)
		}
	}
}

func TestSortableWords(t *testing.T) {
	tests := []struct {
		Words    []string
		Reverse  bool
		Expected []string
	}{
		{
			Words:    []string{"b", "c", "a"},
			Reverse:  false,
			Expected: []string{"a", "b", "c"},
		},
		{
			Words:    []string{"a"},
			Reverse:  false,
			Expected: []string{"a"},
		},

		{
			Words:    []string{"a", "A"},
			Reverse:  false,
			Expected: []string{"A", "a"},
		},
		{

			Words:    []string{"cats", "dogs", "pigs", "man"},
			Reverse:  true,
			Expected: []string{"pigs", "man", "dogs", "cats"},
		},

		{
			Words:    []string{"a", "aa", "aaa"},
			Reverse:  true,
			Expected: []string{"aaa", "aa", "a"},
		},
	}
	for _, test := range tests {
		sw := newSortableWords(test.Words, test.Reverse)
		sort.Sort(sw)
		result := sw.words
		if !reflect.DeepEqual(result, test.Expected) {
			t.Errorf("expected %v, got %v", test.Expected, result)
		}
	}
}

func TestSortArrayBadMethod(t *testing.T) {
	tests := []struct {
		Verb     string
		Expected int
	}{
		{
			Verb:     "POST",
			Expected: http.StatusOK,
		},
		{
			Verb:     "GET",
			Expected: http.StatusMethodNotAllowed,
		},
		{
			Verb:     "PUT",
			Expected: http.StatusMethodNotAllowed,
		},
	}

	body := bytes.NewBuffer(nil)
	body.WriteString(`{"words":["a","b","c"], "reverse": true}`)
	r, _ := http.NewRequest("POST", "", body)
	r.Header.Set("Content-Type", "application/json")

	for _, test := range tests {
		w := httptest.NewRecorder()
		r.Method = test.Verb
		SortArray(w, r)
		if w.Code != test.Expected {
			t.Errorf("expected %d, got %d", test.Expected, w.Code)
		}
	}
}

func TestSortArrayBadContentType(t *testing.T) {
	tests := []struct {
		ContentType string
		Expected    int
	}{
		{
			ContentType: "application/json",
			Expected:    http.StatusOK,
		},
		{
			ContentType: "Application/JSON",
			Expected:    http.StatusUnsupportedMediaType,
		},
		{
			ContentType: "text",
			Expected:    http.StatusUnsupportedMediaType,
		},
	}

	body := bytes.NewBuffer(nil)
	body.WriteString(`{"words":["a","b","c"], "reverse": true}`)
	r, _ := http.NewRequest("POST", "", body)

	for _, test := range tests {
		w := httptest.NewRecorder()
		r.Header.Set("Content-Type", test.ContentType)
		SortArray(w, r)
		if w.Code != test.Expected {
			t.Errorf("expected %d, got %d", test.Expected, w.Code)
		}
	}
}
func TestSortArray(t *testing.T) {
	tests := []struct {
		Words          []string
		ExpectedCode   int
		ExpectedResult string
		Reverse        bool
	}{
		{
			Words:          []string{"cats", "dogs", "turtles"},
			Reverse:        false,
			ExpectedCode:   http.StatusOK,
			ExpectedResult: `{"result":["cats","dogs","turtles"],"reverse":false}`,
		},
		{
			Words:          []string{"9", "dogs", "turtles"},
			Reverse:        false,
			ExpectedCode:   http.StatusBadRequest,
			ExpectedResult: `words must match [A-Za-z]`,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SortArray))

	for _, test := range tests {
		p, _ := json.Marshal(&payload{
			Words:   test.Words,
			Result:  nil,
			Reverse: test.Reverse,
		})
		b := bytes.NewBuffer(nil)
		b.Write(p)
		resp, err := http.Post(ts.URL, "application/json", b)
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
		actualResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
		if resp.StatusCode != test.ExpectedCode {
			t.Errorf("expected response code %d, got %d", test.ExpectedCode, resp.StatusCode)
		}

		if strings.TrimSpace(string(actualResponse)) != test.ExpectedResult {
			t.Errorf("expected response %s, got %s", test.ExpectedResult, string(actualResponse))
		}

		if resp.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected content type of response to be application/json, got %s", resp.Header.Get("Content-Type"))
		}
	}
}
