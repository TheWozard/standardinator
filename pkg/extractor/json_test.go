package extractor_test

import (
	"TheWozard/standardinator/pkg/extractor"
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONExtractor(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		token  string
		output []map[string]interface{}
	}{
		{
			name:   "Empty",
			input:  ``,
			token:  "Anything",
			output: []map[string]interface{}{},
		},
		{
			name:   "EmptyObject",
			input:  `{}`,
			token:  "Anything",
			output: []map[string]interface{}{},
		},
		{
			name:  "Root",
			input: `{}`,
			token: ".",
			output: []map[string]interface{}{
				{},
			},
		},
		{
			name:  "RootWithData",
			input: `{"Stringy":"data", "Numbers": 1, "Boolean": true, "Object": {"More": "data"}, "List": ["value"]}`,
			token: ".",
			output: []map[string]interface{}{
				{
					"Stringy": "data",
					"Numbers": 1.0,
					"Boolean": true,
					"Object": map[string]interface{}{
						"More": "data",
					},
					"List": []interface{}{"value"},
				},
			},
		},
		{
			name:   "TargetToken",
			input:  `{"entities":[{"A":1},{"B":2}]}`,
			token:  "entities",
			output: []map[string]interface{}{{"A": 1.0}, {"B": 2.0}},
		},
		{
			name:   "DeepToken",
			input:  `{"parent":[{"entities":[{"A":1},{"B":2}]},{"entities":[{"C":3},{"D":4}]}]}`,
			token:  "entities",
			output: []map[string]interface{}{{"A": 1.0}, {"B": 2.0}, {"C": 3.0}, {"D": 4.0}},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]%s", i, test.name), func(t *testing.T) {
			extractor := extractor.NewJson(test.token, bytes.NewBufferString(test.input))
			for _, expected := range test.output {
				actual, err := extractor.Next()
				require.NoError(t, err)
				require.Equal(t, expected, actual)
			}
			actual, err := extractor.Next()
			require.Nil(t, actual, "Unexpected Extra Results")
			require.Equal(t, io.EOF, err)
		})
	}
}
