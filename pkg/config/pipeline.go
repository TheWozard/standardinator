package config

import (
	"TheWozard/standardinator/pkg/pipeline"
	"encoding/json"
	"fmt"
)

const (
	SelectStep PipelineStepKind = "select"
)

type PipelineStepKind = string

func GetPipelineStep(kind PipelineStepKind, raw json.RawMessage) (pipeline.PipelineStep, error) {
	switch kind {
	case SelectStep:
		config := pipeline.Select{}
		err := json.Unmarshal(raw, &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	default:
		return nil, fmt.Errorf("invalid step kind '%s'", kind)
	}
}
