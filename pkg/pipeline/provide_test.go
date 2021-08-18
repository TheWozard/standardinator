package pipeline_test

import (
	"TheWozard/standardinator/pkg/pipeline"
	"fmt"
	"testing"
)

func TestProvide(t *testing.T) {

	tests := []struct {
		name string
		data []map[string]interface{}
	}{
		{
			name: "Empty",
			data: []map[string]interface{}{},
		},
		{
			name: "Single",
			data: []map[string]interface{}{
				{"Data": "A"},
			},
		},
		{
			name: "Multiple",
			data: []map[string]interface{}{
				{"Data": "A"}, {"Data": "B"}, {"Data": "C"},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]%s", i, test.name), func(t *testing.T) {
			next := pipeline.ProvideStep{
				Data: test.data,
			}.Init(pipeline.NextEOF)

			VerifyNext(t, next, test.data)
		})
	}
}
