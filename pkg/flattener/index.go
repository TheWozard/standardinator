package flattener

import (
	"errors"
	"io"
)

var (
	ErrAlreadyCompleted = errors.New("calling close on an already completed OutputBuilder")
)

// Decoder creates new Extractors for provided io.Reader
type Decoder interface {
	// New returns a new Extractor that will read from r
	New(r io.Reader) Extractor
}

// Extractor provides an Next function to iterate over the output contexts of objects
type Extractor interface {
	// Next gets the next available context or returns io.EOF when no more remain
	Next() (*IntermediateContext, error)
}

// Evaluator provides an Next function to iterate over the output objects with any joins performed
type Evaluator interface {
	// Next gets the next available output or returns io.EOF when no more remain
	Next() (*Output, error)
}
