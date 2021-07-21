package extractor

// Extractor provides an Next function to iterate over a set of elements
type Extractor interface {
	// Next gets the next available element or returns io.EOF when no more remain
	Next() (map[string]interface{}, error)
}

// TODO: Common Decoder