package token

// Token either a StartToken, EndToken, or KVToken that represents the
type Token interface{}

// StartToken defines the start of some namespace of data
type StartToken struct {
	Key   string
	Array bool
}

// EndToken defines the end of some namespace of data
type EndToken struct {
}

// KVToken represents a data value and its assigned name
type KVToken struct {
	Key   string
	Value interface{}
}
