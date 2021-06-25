package config

import (
	"TheWozard/standardinator/pkg/token"
	"TheWozard/standardinator/pkg/token/json"
	"fmt"
	"io"
)

func NewTokenizer(t token.ReaderType, r io.Reader) (token.Reader, error) {
	switch t {
	case token.JSONReader:
		return json.NewTokenizer(r), nil
	default:
		return nil, fmt.Errorf("unknown tokenizer type '%s'", t)
	}
}
