package output

func NewManager() Manager {
	return &receiver{}
}

type receiver struct {
	key  string
	data valueBuilder

	parent *receiver
}

func (r *receiver) HasResult() bool {
	return false
}

func (r *receiver) GetResult() *Result {
	return nil
}

func (r *receiver) Receive(key string, value interface{}) {
	r.data.AddEntry(key, value)
}

func (r *receiver) CreateChildNode() Manager {
	return r
}

func (r *receiver) Flush() {
	r.parent.Receive(r.key, r.data.GetValue())
}
