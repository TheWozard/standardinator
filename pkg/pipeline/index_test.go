package pipeline_test

import (
	"TheWozard/standardinator/pkg/pipeline"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// VerifyNext ensures the next function provides the expect list of output data and ends with an io.EOF
func VerifyNext(t *testing.T, next pipeline.GetNext, expected []map[string]interface{}) {
	index := 0
	for {
		data, err := next()
		if err != nil {
			require.Equal(t, io.EOF, err)
			require.Equal(t, len(expected), index)
			break
		}
		require.Equal(t, expected[index], data)
		index += 1
	}
}

func TestPipeline(t *testing.T) {

	tests := []struct {
		name  string
		steps []pipeline.PipelineStep
		data  []map[string]interface{}
	}{
		{
			name:  "Empty",
			steps: []pipeline.PipelineStep{},
			data:  []map[string]interface{}{},
		},
		{
			name: "SingleButEmpty",
			steps: []pipeline.PipelineStep{
				pipeline.ProvideStep{
					Data: []map[string]interface{}{},
				},
			},
			data: []map[string]interface{}{},
		},
		{
			name: "Single",
			steps: []pipeline.PipelineStep{
				pipeline.ProvideStep{
					Data: []map[string]interface{}{
						{"Data": "A"},
					},
				},
			},
			data: []map[string]interface{}{
				{"Data": "A"},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]%s", i, test.name), func(t *testing.T) {
			pipe := pipeline.NewPipeline(test.steps, pipeline.NextEOF)

			VerifyNext(t, pipe.Next, test.data)
			require.NoError(t, pipe.Close())
		})
	}
}
