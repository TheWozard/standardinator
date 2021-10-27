package data

// TODO: Consider a raw interface and functions hung on it to determin type
type Raw = map[string]interface{}

// TODO: Better naming
type Payload struct {
	Name string
	Data Raw
}
