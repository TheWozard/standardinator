package extractor

import (
	"TheWozard/standardinator/pkg/matcher"
	"encoding/json"
	"io"
)

// JSONConfig defines configuration for a JSON based Extractor
type JSONConfig struct {
	Token string `json:"token"`
}

// GetMatcher provides the matcher for this config
// TODO: expand support for matching
func (c JSONConfig) GetMatcher() *matcher.Simple {
	return matcher.NewSimple(c.Token)
}

// NewJsonExtractor creates an Extractor that reads JSON data from r
func NewJsonExtractor(config JSONConfig, r io.Reader) Extractor {
	decoder := json.NewDecoder(r)
	return &jsonExtractor{
		config:  config,
		decoder: decoder,

		matcher: config.GetMatcher(),
	}
}

type jsonExtractor struct {
	config  JSONConfig
	decoder *json.Decoder

	matcher *matcher.Simple
}

func (e *jsonExtractor) Next() (map[string]interface{}, error) {
	for {
		for !e.matcher.Matched {
			// Searching for the next match
			token, err := e.decoder.Token()
			if err != nil {
				// End condition when the decoder returns a io.EOF
				return nil, err
			}
			switch typed := token.(type) {
			case string:
				// TODO: could match a string value and not just a key
				if e.matcher.Check(typed) {
					// Need to consume the next token to enter the object
					raw, err := e.decoder.Token()
					if err != nil {
						return nil, err
					}
					switch delim := raw.(type) {
					case json.Delim:
						if delim != '[' {
							e.matcher.Reset()
						}
					default:
						e.matcher.Reset()
					}
				}
			}
		}
		// Extract elements
		for e.decoder.More() {
			data := map[string]interface{}{}
			e.decoder.Decode(&data)
			return data, nil
		}
		// Start over
		e.matcher.Reset()
	}
}
