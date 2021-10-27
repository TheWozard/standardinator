package extractor_test

import (
	"TheWozard/standardinator/pkg/data"
	"TheWozard/standardinator/pkg/extractor"
	"bytes"
	"fmt"
	"io"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiExtractor(t *testing.T) {

	decoder := extractor.Multi{
		Decoders: []extractor.Decoder{
			extractor.JSON{Token: "Parent"},
			extractor.JSON{Token: "Child"},
		},
	}

	tests := []struct {
		runs   int
		name   string
		input  string
		config extractor.Decoder
		output map[string][]data.Payload
	}{
		// {
		// 	runs:   1,
		// 	name:   "Empty",
		// 	input:  `{}`,
		// 	config: decoder,
		// 	output: []data.Payload{},
		// },
		{
			runs:   20,
			name:   "Simple",
			input:  `{"Parent":[{"Child":[{},{},{}]}]}`,
			config: decoder,
			output: map[string][]data.Payload{
				"Child": {
					{Name: "Child", Data: map[string]interface{}{}},
					{Name: "Child", Data: map[string]interface{}{}},
					{Name: "Child", Data: map[string]interface{}{}},
				},
				"Parent": {
					{Name: "Parent", Data: map[string]interface{}{
						"Child": []interface{}{map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}},
					}},
				},
			},
		},
	}

	for i, test := range tests {
		for run := 0; run < test.runs; run++ {
			t.Run(fmt.Sprintf("[%d][%d]%s", i, run, test.name), func(t *testing.T) {
				file := bytes.NewBufferString(test.input)
				precount := runtime.NumGoroutine()
				extractor := decoder.New(file)
				fmt.Println("Start")
				fmt.Println(runtime.NumGoroutine())
				totalOutput := 0
				for _, outputs := range test.output {
					totalOutput += len(outputs)
				}
				for i := 0; i < totalOutput; i++ {
					actual, err := extractor.Next()
					fmt.Println(runtime.NumGoroutine())
					require.NoError(t, err)
					expected, remaining := test.output[actual.Name][0], test.output[actual.Name][1:]
					require.Equal(t, expected, *actual)
					test.output[actual.Name] = remaining
				}
				actual, err := extractor.Next()
				fmt.Println("Final")
				fmt.Println(runtime.NumGoroutine())
				require.Nil(t, actual, "Unexpected Extra Results")
				require.Equal(t, io.EOF, err)
				require.Equal(t, precount, runtime.NumGoroutine(), "Unexpected NumGoroutine")
			})
		}
	}

}
