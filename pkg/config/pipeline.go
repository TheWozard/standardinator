package config

import (
	"TheWozard/standardinator/pkg/pipeline"
	"encoding/json"
	"fmt"
)

// PipelineStepKind provides a type for configuring a step
type PipelineStepKind = string

const (
	// SelectStep selects an element out of the passed data
	SelectStep PipelineStepKind = "select"
)

var (
	// StepIndex provides the index for decoding incoming
	StepIndex map[PipelineStepKind]pipeline.PipelineStepDecoder = map[PipelineStepKind]pipeline.PipelineStepDecoder{
		SelectStep: pipeline.DecodeSelect,
	}
)

// Converts a json.RawMessage into a configured pipeline.PipelineStep
func GetPipelineStep(kind PipelineStepKind, raw json.RawMessage) (pipeline.PipelineStep, error) {
	if decode, ok := StepIndex[kind]; ok {
		return decode(raw)
	}
	return nil, fmt.Errorf("invalid step kind '%s'", kind)
}
