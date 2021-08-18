package pipeline

import (
	"context"
	"encoding/json"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

// DecodeSelect ...
func DecodeSelect(raw json.RawMessage) (PipelineStep, error) {
	config := map[string]string{}
	err := json.Unmarshal(raw, &config)
	if err != nil {
		return nil, err
	}
	paths := map[string]gval.Evaluable{}
	for key, val := range config {
		paths[key], err = jsonpath.New(val)
		if err != nil {
			return nil, err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &SelectStep{
		Context: ctx,
		Cancel:  cancel,
		Paths:   paths,
	}, nil
}

type SelectStep struct {
	Context context.Context
	Cancel  context.CancelFunc
	Paths   map[string]gval.Evaluable
}

func (s SelectStep) Init(prev GetNext) GetNext {
	return Wrap(prev, func(data map[string]interface{}) (map[string]interface{}, error) {
		new := map[string]interface{}{}
		for key, eval := range s.Paths {
			data, err := eval(s.Context, data)
			if err != nil {
				return nil, err
			}
			new[key] = data
		}
		return new, nil
	})
}

func (s SelectStep) Close() error {
	s.Cancel()
	return nil
}
