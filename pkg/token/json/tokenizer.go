package json

import (
	"TheWozard/standardinator/pkg/token"
	"encoding/json"
	"fmt"
	"io"
)

func NewTokenizer(r io.Reader) token.Reader {
	return &tokenizer{
		decoder: json.NewDecoder(r),
		stack:   []keyProvider{newStaticKeyProvider("$")},
	}
}

type keyProvider func() (string, token.Token, error)

type tokenizer struct {
	decoder *json.Decoder

	stack []keyProvider
}

func (t *tokenizer) Next() (token.Token, error) {
	key, keyToken, err := t.stack[len(t.stack)-1]()
	if err != nil {
		return nil, err
	}
	if keyToken != nil {
		return keyToken, nil
	}
	next, err := t.decoder.Token()
	if err != nil {
		return nil, err
	}
	switch typed := next.(type) {
	case json.Delim:
		if typed == '{' {
			t.stack = append(t.stack, t.tokenKeyProvider)
			return token.StartToken{
				Key:   key,
				Array: false,
			}, nil
		}
		if typed == '[' {
			t.stack = append(t.stack, newIndexProvider())
			return token.StartToken{
				Key:   key,
				Array: true,
			}, nil
		}
		t.stack = t.stack[:len(t.stack)-1]
		return token.EndToken{}, nil
	default:
		return token.KVToken{
			Key:   key,
			Value: typed,
		}, nil
	}
}

func (t *tokenizer) tokenKeyProvider() (string, token.Token, error) {
	next, err := t.decoder.Token()
	if err != nil {
		return "", nil, err
	}
	switch typed := next.(type) {
	case json.Delim:
		t.stack = t.stack[:len(t.stack)-1]
		return "", token.EndToken{}, nil
	default:
		return fmt.Sprint(typed), nil, nil
	}
}

func newStaticKeyProvider(key string) keyProvider {
	return func() (string, token.Token, error) { return key, nil, nil }
}

func newIndexProvider() keyProvider {
	index := 0
	return func() (string, token.Token, error) {
		defer func() { index = index + 1 }()
		return fmt.Sprintf("%d", index), nil, nil
	}
}
