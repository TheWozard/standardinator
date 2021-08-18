package pipeline

import (
	"encoding/json"
	"io"
)

// GetNext core iteration function of the pipeline process
type GetNext func() (map[string]interface{}, error)

// NextEOF returns io.EOF and works as an GetNext. Most useful for unit testing
func NextEOF() (map[string]interface{}, error) {
	return nil, io.EOF
}

// ProvideNext provides transformation data for the passed data
type ProvideNext func(data map[string]interface{}) (map[string]interface{}, error)

// ProvideBufferedNext provides multiple transformation data structures for the passed data
type ProvideBufferedNext func(data map[string]interface{}) ([]map[string]interface{}, error)

// PipelineStep a single modification operation in a passed pipeline
type PipelineStep interface {
	Init(GetNext) GetNext
	Close() error
}

// PipelineStepDecoder converts a json.RawMessage into a valid PipelineStep
type PipelineStepDecoder func(raw json.RawMessage) (PipelineStep, error)

// Wrap wraps a GetNext function to provide automatic error handling from previous steps
func Wrap(prev GetNext, provide ProvideNext) GetNext {
	return func() (map[string]interface{}, error) {
		data, err := prev()
		if err != nil {
			return nil, err
		}
		return provide(data)
	}
}

// WrapBuffered wraps a GetNext function allowing for the provider function to return a list of elements to each be returned in order
func WrapBuffered(prev GetNext, provide ProvideBufferedNext) GetNext {
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
			buffer, err = provide(data)
			if err != nil {
				return nil, err
			}
		}
	}
}
