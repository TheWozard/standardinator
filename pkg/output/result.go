package output

// Result a finally output object that has completed standardization
type Result struct {
	Name string
	Data interface{}
}

type valueBuilder interface {
	AddEntry(key string, value interface{})
	GetValue() interface{}
}

type objectValueBuilder struct {
	data map[string]interface{}
}

func (b *objectValueBuilder) AddEntry(key string, value interface{}) {
	b.data[key] = value
}

func (b *objectValueBuilder) GetValue() interface{} {
	return b.data
}

type arrayValueBuilder struct {
	data []interface{}
}

func (b *arrayValueBuilder) AddEntry(key string, value interface{}) {
	// TODO: kinda lazy but works for now.
	b.data = append(b.data, value)
}

func (b *arrayValueBuilder) GetValue() interface{} {
	return b.data
}

type noOpValueBuilder struct {
	data []interface{}
}

func (b *noOpValueBuilder) AddEntry(key string, value interface{}) {

}

func (b *noOpValueBuilder) GetValue() interface{} {
	return nil
}
