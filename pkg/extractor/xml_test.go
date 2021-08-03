package extractor_test

import (
	"TheWozard/standardinator/pkg/extractor"
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXMLExtractor(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config extractor.XML
		output []map[string]interface{}
	}{
		{
			name:   "Empty",
			input:  ``,
			config: extractor.XML{},
			output: []map[string]interface{}{},
		},
		{
			name:   "NoMatch",
			input:  `<data></data>`,
			config: extractor.XML{},
			output: []map[string]interface{}{},
		},
		{
			name:   "EmptyObject",
			input:  `<data></data>`,
			config: extractor.XML{Token: "data"},
			output: []map[string]interface{}{{}},
		},
		{
			name:   "ObjectWithData",
			input:  `<data att="value">5<A att="foo">1</A><B>2</B><C></C></data>`,
			config: extractor.XML{Token: "data"},
			output: []map[string]interface{}{
				{
					"@att":  "value",
					"#text": "5",
					"A": map[string]interface{}{
						"@att":  "foo",
						"#text": "1",
					},
					"B": "2",
					"C": map[string]interface{}{},
				},
			},
		},
		{
			name:   "MultiTargets",
			input:  `<entities><data>A</data><data>B</data><data>C</data></entities>`,
			config: extractor.XML{Token: "data"},
			output: []map[string]interface{}{
				{
					"#text": "A",
				},
				{
					"#text": "B",
				},
				{
					"#text": "C",
				},
			},
		},
		{
			name:   "SplitTargets",
			input:  `<entities><data>A</data></entities><entities><data>B</data></entities>`,
			config: extractor.XML{Token: "data"},
			output: []map[string]interface{}{
				{
					"#text": "A",
				},
				{
					"#text": "B",
				},
			},
		},
		{
			name:   "RepeatElements",
			input:  `<entities><data>A</data><data>B</data></entities>`,
			config: extractor.XML{Token: "entities", Repeats: []string{"data"}},
			output: []map[string]interface{}{
				{
					"data": []interface{}{
						"A",
						"B",
					},
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]%s", i, test.name), func(t *testing.T) {
			extractor := test.config.New(bytes.NewBufferString(test.input))
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
