package token_test

import (
	"TheWozard/standardinator/pkg/token"
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
				token.StartToken{Key: "$", Repeats: false},
				token.KVToken{Key: "Data", Value: "Test"},
				token.EndToken{},
			},
		},
		{
			name:  "Correct Types",
			input: `{"Int":1,"Bool":true,"Float":0.2,"String":"Test","Nil":null}`,
			tokens: []token.Token{
				token.StartToken{Key: "$", Repeats: false},
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
				token.StartToken{Key: "$", Repeats: false},
				token.StartToken{Key: "A", Repeats: false},
				token.EndToken{},
				token.StartToken{Key: "B", Repeats: true},
				token.EndToken{},
				token.EndToken{},
			},
		},
		// {
		// 	name:   "Stacked Namespaces",
		// 	input:  `{"A":{"B":[ "Data" ]},"B": []}`, // How do we handle this
		// 	tokens: []token.Token{
		// 		// TODO: Solve
		// 	},
		// },
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, test.name), func(t *testing.T) {
			tokenizer := token.NewJSONTokenizer(bytes.NewBufferString(test.input))
			for _, expected := range test.tokens {
				token, err := tokenizer.Next()
				require.NoError(t, err)
				require.Equal(t, expected, token)
			}
			final, err := tokenizer.Next()
			require.Equal(t, err, io.EOF)
			require.Nil(t, final)
		})
	}
}
