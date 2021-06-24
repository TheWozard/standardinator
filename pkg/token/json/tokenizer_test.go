package json_test

import (
	"TheWozard/standardinator/pkg/token"
	"TheWozard/standardinator/pkg/token/json"
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []token.Token
	}{
		{
			name:  "Simple",
			input: `{"Data":"Test"}`,
			tokens: []token.Token{
				token.StartToken{Key: "$", Array: false},
				token.KVToken{Key: "Data", Value: "Test"},
				token.EndToken{},
			},
		},
		{
			name:  "Correct Types",
			input: `{"Int":1,"Bool":true,"Float":0.2,"String":"Test","Nil":null}`,
			tokens: []token.Token{
				token.StartToken{Key: "$", Array: false},
				token.KVToken{Key: "Int", Value: 1.0}, // All number come across as floats. Thanks Javascript.
				token.KVToken{Key: "Bool", Value: true},
				token.KVToken{Key: "Float", Value: 0.2},
				token.KVToken{Key: "String", Value: "Test"},
				token.KVToken{Key: "Nil", Value: nil},
				token.EndToken{},
			},
		},
		{
			name:  "Namespaces",
			input: `{"A":{},"B":[]}`,
			tokens: []token.Token{
				token.StartToken{Key: "$", Array: false},
				token.StartToken{Key: "A", Array: false},
				token.EndToken{},
				token.StartToken{Key: "B", Array: true},
				token.EndToken{},
				token.EndToken{},
			},
		},
		{
			name:  "Stacked Namespaces",
			input: `{"A":{"B":[ "Data" ],"C": []}}`,
			tokens: []token.Token{
				token.StartToken{Key: "$"},
				token.StartToken{Key: "A"},
				token.StartToken{Key: "B", Array: true},
				token.KVToken{Key: "0", Value: "Data"},
				token.EndToken{},
				token.StartToken{Key: "C", Array: true},
				token.EndToken{},
				token.EndToken{},
				token.EndToken{},
			},
		},
		{
			name:  "Lists of simple elements",
			input: `["A","B","C","D","E","F"]`,
			tokens: []token.Token{
				token.StartToken{Key: "$", Array: true},
				token.KVToken{Key: "0", Value: "A"},
				token.KVToken{Key: "1", Value: "B"},
				token.KVToken{Key: "2", Value: "C"},
				token.KVToken{Key: "3", Value: "D"},
				token.KVToken{Key: "4", Value: "E"},
				token.KVToken{Key: "5", Value: "F"},
				token.EndToken{},
			},
		},
		{
			name:  "Lists of objects",
			input: `[{},{},{}]`,
			tokens: []token.Token{
				token.StartToken{Key: "$", Array: true},
				token.StartToken{Key: "0"},
				token.EndToken{},
				token.StartToken{Key: "1"},
				token.EndToken{},
				token.StartToken{Key: "2"},
				token.EndToken{},
				token.EndToken{},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, test.name), func(t *testing.T) {
			tokenizer := json.NewTokenizer(bytes.NewBufferString(test.input))
			for i, expected := range test.tokens {
				token, err := tokenizer.Next()
				require.NoError(t, err, fmt.Sprintf("Unexpected error on index %d", i))
				require.Equal(t, expected, token, fmt.Sprintf("Unexpected token on index %d", i))
			}
			final, err := tokenizer.Next()
			require.Equal(t, err, io.EOF)
			require.Nil(t, final)
		})
	}
}
