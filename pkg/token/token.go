package token

type Token interface{}

// StartToken defines the start of some namespace of data
type StartToken struct {
	Key     string
	Repeats bool
}

// EndToken defines the end of some namespace of data
type EndToken struct {
}

// KVToken represents a data value and its assigned name
type KVToken struct {
	Key   string
	Value interface{}
}
