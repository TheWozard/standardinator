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
				},
			},
			paths: map[string]string{
				"LikesCars":      "$.Likes.Cars",
				"LikesComputers": "$.Likes.Computers",
				"Likes":          "$.Likes.*",
			},
			post: []map[string]interface{}{
				{"LikesCars": true, "LikesComputers": true, "Likes": []interface{}{true, true}},
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
