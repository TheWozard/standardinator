package extractor

import (
	"encoding/json"
	"io"
)

// TODO: JsonConfig
func NewJson(token string, reader io.Reader) Extractor {
	decoder := json.NewDecoder(reader)
	return &jsonExtractor{
		token:   token,
		decoder: decoder,

		matched: token == ".",
	}
}

type jsonExtractor struct {
	token   string
	decoder *json.Decoder

	// TODO: Better matching
	matched bool
}

func (e *jsonExtractor) Next() (map[string]interface{}, error) {
	for {
		for !e.matched {
			token, err := e.decoder.Token()
			if err != nil {
				return nil, err
			}
			switch typed := token.(type) {
			case string:
				if typed == e.token {
					e.decoder.Token()
					e.matched = true
				}
			}
		}
		for e.decoder.More() {
			data := map[string]interface{}{}
			e.decoder.Decode(&data)
			return data, nil
		}
		e.matched = false
	}
}
