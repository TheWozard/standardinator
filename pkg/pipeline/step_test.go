package pipeline_test

import (
	"TheWozard/standardinator/pkg/pipeline"
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {

	simpleProvider := func(data map[string]interface{}) (map[string]interface{}, error) {
		data["provider"] = "processed"
		return data, nil
	}

	tests := []struct {
		name     string
		pre      []map[string]interface{}
		post     []map[string]interface{}
		provider pipeline.ProvideNext
	}{
		{
			name:     "Empty",
			pre:      []map[string]interface{}{},
			post:     []map[string]interface{}{},
			provider: simpleProvider,
		},
		{
			name: "Single",
			pre: []map[string]interface{}{
				{},
			},
			post: []map[string]interface{}{
				{"provider": "processed"},
			},
			provider: simpleProvider,
		},
		{
			name: "Multi",
			pre: []map[string]interface{}{
				{}, {}, {},
			},
			post: []map[string]interface{}{
				{"provider": "processed"},
				{"provider": "processed"},
				{"provider": "processed"},
			},
			provider: simpleProvider,
		},
		{
			name: "MultiData",
			pre: []map[string]interface{}{
				{"Data": "A"}, {"Data": "B"}, {"Data": "C"},
			},
			post: []map[string]interface{}{
				{"provider": "processed", "Data": "A"},
				{"provider": "processed", "Data": "B"},
				{"provider": "processed", "Data": "C"},
			},
			provider: simpleProvider,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]%s", i, test.name), func(t *testing.T) {
			pre := pipeline.ProvideStep{
				Data: test.pre,
			}.Init(pipeline.NextEOF)
			next := pipeline.Wrap(pre, test.provider)

			VerifyNext(t, next, test.post)
		})
	}
}

func TestBufferedWrap(t *testing.T) {

	cloneProvider := func(count int) pipeline.ProvideBufferedNext {
		return func(data map[string]interface{}) ([]map[string]interface{}, error) {
			clones := make([]map[string]interface{}, count)
			for i, _ := range clones {
				clones[i] = map[string]interface{}{}
				for key, value := range data {
					clones[i][key] = value
				}
				clones[i]["provider"] = fmt.Sprintf("clone-%d", i+1)
			}
			return clones, nil
		}
	}

	tests := []struct {
		name     string
		pre      []map[string]interface{}
		post     []map[string]interface{}
		provider pipeline.ProvideBufferedNext
	}{
		{
			name:     "Empty",
			pre:      []map[string]interface{}{},
			post:     []map[string]interface{}{},
			provider: cloneProvider(3),
		},
		{
			name: "Single",
			pre: []map[string]interface{}{
				{},
			},
			post: []map[string]interface{}{
				{"provider": "clone-1"},
				{"provider": "clone-2"},
				{"provider": "clone-3"},
			},
			provider: cloneProvider(3),
		},
		{
			name: "Multi",
			pre: []map[string]interface{}{
				{}, {}, {},
			},
			post: []map[string]interface{}{
				{"provider": "clone-1"},
				{"provider": "clone-2"},
				{"provider": "clone-1"},
				{"provider": "clone-2"},
				{"provider": "clone-1"},
				{"provider": "clone-2"},
			},
			provider: cloneProvider(2),
		},
		{
			name: "MultiData",
			pre: []map[string]interface{}{
				{"Data": "A"}, {"Data": "B"}, {"Data": "C"},
			},
			post: []map[string]interface{}{
				{"provider": "clone-1", "Data": "A"},
				{"provider": "clone-2", "Data": "A"},
				{"provider": "clone-1", "Data": "B"},
				{"provider": "clone-2", "Data": "B"},
				{"provider": "clone-1", "Data": "C"},
				{"provider": "clone-2", "Data": "C"},
			},
			provider: cloneProvider(2),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]%s", i, test.name), func(t *testing.T) {
			pre := pipeline.ProvideStep{
				Data: test.pre,
			}.Init(pipeline.NextEOF)
			next := pipeline.WrapBuffered(pre, test.provider)

			VerifyNext(t, next, test.post)
		})
	}
}
