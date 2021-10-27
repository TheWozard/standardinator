package pipeline_test

import (
	"TheWozard/standardinator/pkg/pipeline"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {

	tests := []struct {
		name  string
		pre   []map[string]interface{}
		paths map[string]string
		post  []map[string]interface{}
	}{
		{
			name:  "Empty",
			pre:   []map[string]interface{}{},
			paths: map[string]string{},
			post:  []map[string]interface{}{},
		},
		{
			name: "SingleSimple",
			pre: []map[string]interface{}{
				{
					"FirstName": "John",
					"LastName":  "Doe",
					"Age":       34,
					"Height":    70,
					"Likes": map[string]interface{}{
						"Cars":      true,
						"Computers": true,
					},
					"prefered": 1,
					"Files": []interface{}{
						map[string]interface{}{
							"Name": "A",
							"Url":  "<URL>",
						},
						map[string]interface{}{
							"Name": "B",
							"Url":  "<URL>",
						},
						map[string]interface{}{
							"Name": "C",
							"Url":  "<URL>",
						},
						map[string]interface{}{
							"Name": "D",
							"Url":  "<URL>",
						},
					},
					"Data1": "DataA",
					"Data2": "DataB",
					"Data3": "DataC",
				},
			},
			paths: map[string]string{
				"Data": `{#0: $[?(@ =~ "Data.*")]}`,
			},
			post: []map[string]interface{}{
				{},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]%s", i, test.name), func(t *testing.T) {
			raw, err := json.Marshal(test.paths)
			require.NoError(t, err)
			step, err := pipeline.DecodeSelect(raw)
			require.NoError(t, err)
			next := step.Init(
				pipeline.ProvideStep{
					Data: test.pre,
				}.Init(pipeline.NextEOF),
			)

			VerifyNext(t, next, test.post)
		})
	}
}
