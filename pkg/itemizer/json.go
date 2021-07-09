package itemizer

import (
	"TheWozard/standardinator/pkg/path"
	"TheWozard/standardinator/pkg/path/jsonpath"
	"encoding/json"
	"fmt"
	"io"
)

func CreateJSONReader(input io.ReadCloser, rawPath string) (Reader, error) {
	parsed, err := jsonpath.Parse(rawPath)
	if err != nil {
		return nil, err
	}
	if !parsed.IsRoot() {
		return nil, fmt.Errorf("expected path to match root")
	}
	decoder := json.NewDecoder(input)
	decoder.Token()
	return &jsonReader{
		reader:  input,
		decoder: decoder,
		stack:   &cursorStack{stack: []path.Cursor{parsed.MatchRoot()}},
	}, nil
}

type jsonReader struct {
	reader  io.ReadCloser
	decoder *json.Decoder
	stack   *cursorStack
}

func (j *jsonReader) Next() (interface{}, error) {
	if !j.stack.IsMatched() {
		token, err := j.decoder.Token()
		if err != nil {
			return nil, err
		}
		switch typed := token.(type) {
		case json.Delim:
			if j.stack.IsMatched() {
				break
			}
			if typed == ']' || typed == '}' {
				j.stack.Pop()
			}
		case string:
			j.stack.Match(typed)
			j.decoder.Token()
		}
	}
	for j.decoder.More() {
		d := data{}
		j.decoder.Decode(&d)
		return d, nil
	}
	j.stack.Pop()
	j.decoder.Token()
	return j.Next()
}

type data struct{}

func (d *data) UnmarshalJSON(b []byte) error {
	fmt.Println(string(b))
	return nil
}

type cursorStack struct {
	stack []path.Cursor
}

func (path *cursorStack) IsMatched() bool {
	if len(path.stack) == 0 {
		return false
	}
	target := path.stack[len(path.stack)-1]
	if target == nil {
		return false
	}
	return target.IsMatched()
}

func (path *cursorStack) Match(target string) {
	path.stack = append(path.stack, path.stack[len(path.stack)-1].Match(target))
}

func (path *cursorStack) Pop() {
	if len(path.stack) > 0 {
		path.stack = path.stack[:len(path.stack)-1]
	}
}
