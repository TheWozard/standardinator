package extractor

import (
	"TheWozard/standardinator/pkg/data"
	"io"
)

// Extractor provides an Next function to iterate over a set of elements
type Extractor interface {
	// Next gets the next available payload or returns io.EOF when no more remain
	Next() (*data.Payload, error)
}

// Decoder creates new Extractors for provided io.Reader
type Decoder interface {
	// New returns a new Extractor that will read from r
	New(r io.Reader) Extractor
}
