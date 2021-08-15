package pipeline

import (
	"TheWozard/standardinator/pkg/extractor"

	"github.com/hashicorp/go-multierror"
)

type GetNext func() (map[string]interface{}, error)

// Pipeline
type Pipeline struct {
	name    string
	next    GetNext
	steps   []PipelineStep
	extract extractor.Extractor
}

// PipelineElement
type PipelineElement struct {
	Name string
	Data map[string]interface{}
}

// Next
func (p *Pipeline) Next() (*PipelineElement, error) {
	data, err := p.next()
	if err != nil {
		return nil, err
	}
	return &PipelineElement{
		Name: p.name,
		Data: data,
	}, nil
}

func (p *Pipeline) Close() error {
	var result error

	for _, step := range p.steps {
		if err := step.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

type PipelineStep interface {
	Init(GetNext) GetNext
	Close() error
}

func NewPipeline(name string, steps []PipelineStep, extract extractor.Extractor) *Pipeline {
	next := extract.Next
	for _, step := range steps {
		next = step.Init(next)
	}

	return &Pipeline{
		name:    name,
		next:    next,
		steps:   steps,
		extract: extract,
	}
}

func Wrap(prev GetNext, new func(data map[string]interface{}) (map[string]interface{}, error)) GetNext {
	return func() (map[string]interface{}, error) {
		data, err := prev()
		if err != nil {
			return nil, err
		}
		return new(data)
	}
}

func WrapBuffered(prev GetNext, new func(data map[string]interface{}) ([]map[string]interface{}, error)) GetNext {
	buffer := []map[string]interface{}{}
	return func() (map[string]interface{}, error) {
		for {
			if len(buffer) > 0 {
				var rtn map[string]interface{}
				rtn, buffer = buffer[0], buffer[1:]
				return rtn, nil
			}
			data, err := prev()
			if err != nil {
				return nil, err
			}
			buffer, err = new(data)
			if err != nil {
				return nil, err
			}
		}
	}
}
