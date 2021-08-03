package extractor

import "io"

// Extractor provides an Next function to iterate over a set of elements
type Extractor interface {
	// Next gets the next available element or returns io.EOF when no more remain
	// TODO: more generic return to support non map based data
	Next() (map[string]interface{}, error)
}

// Decoder creates new Extractors for provided io.Reader
type Decoder interface {
	// New returns a new Extractor that will read from r
	New(r io.Reader) Extractor
}
