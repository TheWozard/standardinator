package token

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewJSONTokenizer(r io.Reader) Reader {
	return &jsonTokenizer{
		decoder: json.NewDecoder(r),
	}
}

type jsonTokenizer struct {
	decoder *json.Decoder

	arrayMode bool
	index     int
}

func (t jsonTokenizer) Next() (Token, error) {
	token, err := t.decoder.Token()
	if err != nil {
		return nil, err
	}
	switch typed := token.(type) {
	case json.Delim:
		if typed == '{' || typed == '[' {
			return StartToken{
				Key:     "$",
				Repeats: typed == '[',
			}, nil
		}
		return EndToken{}, nil
	default:
		key := fmt.Sprint(typed)
		value, err := t.decoder.Token()
		if err != nil {
			return nil, err
		}
		switch valueTyped := value.(type) {
		case json.Delim:
			return StartToken{
				Key:     key,
				Repeats: valueTyped == '[',
			}, nil
		default:
			return KVToken{
				Key:   key,
				Value: value,
			}, nil
		}
	}
}
