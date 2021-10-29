package flattener_test

import (
	"TheWozard/standardinator/pkg/flattener"
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
		config flattener.XML
		output []*flattener.Output
	}{
		{
			name:   "Empty",
			input:  ``,
			config: flattener.XML{},
			output: []*flattener.Output{},
		},
		{
			name:   "NoOutput",
			input:  `<data><f1>A</f1><f2>B</f2><f3>C</f3></data>`,
			config: flattener.XML{},
			output: []*flattener.Output{},
		},
		{
			name:  "SingleOutput",
			input: `<data><f1>A</f1><f2>B</f2><f3>C</f3></data>`,
			config: flattener.XML{Config: flattener.XMLOutputConfig{
				{
					Path: flattener.NewPath("data"),
				},
			}},
			output: []*flattener.Output{
				{Name: "data", Data: map[string]interface{}{"f1": "A", "f2": "B", "f3": "C"}},
			},
		},
		{
			name:  "SingleDeepOutput",
			input: `<data><correct><target att="1">A</target></correct><incorrect><target att="2">B</target></incorrect></data>`,
			config: flattener.XML{Config: flattener.XMLOutputConfig{
				{
					Path: flattener.NewPath("correct.target"),
				},
			}},
			output: []*flattener.Output{
				{Name: "correct.target", Data: map[string]interface{}{"@att": "1", "#text": "A"}},
			},
		},
		{
			name:  "SingleDeepOutput",
			input: `<data><correct><target att="1">A</target></correct><incorrect><target att="2">B</target></incorrect></data>`,
			config: flattener.XML{Config: flattener.XMLOutputConfig{
				{
					Path: flattener.NewPath("correct.target"),
				},
			}},
			output: []*flattener.Output{
				{Name: "correct.target", Data: map[string]interface{}{"@att": "1", "#text": "A"}},
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
