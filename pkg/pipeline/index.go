package pipeline

import (
	"github.com/hashicorp/go-multierror"
)

// Pipeline an overall view of a complete pipeline and the primary location for interacting with it.
type Pipeline interface {
	// Next gets the next element in the pipeline. Returns io.EOF at the end.
	Next() (map[string]interface{}, error)
	// Closes the pipeline and all steps. Should always be called once done.
	Close() error
}

// NewPipeline creates a new Pipeline based on the passed PipelineStep pulling from the passed provider
func NewPipeline(steps []PipelineStep, provider GetNext) Pipeline {
	next := provider
	for _, step := range steps {
		next = step.Init(next)
	}

	return &pipeline{
		next:  next,
		steps: steps,
	}
}

// pipeline implementation of the Pipeline interface
type pipeline struct {
	next  GetNext
	steps []PipelineStep
}

// Next gets the next element in the pipeline. Returns io.EOF at the end.
func (p *pipeline) Next() (map[string]interface{}, error) {
	return p.next()
}

// Closes the pipeline and all steps. Should always be called once done.
func (p *pipeline) Close() error {
	var result error

	for _, step := range p.steps {
		if err := step.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}
