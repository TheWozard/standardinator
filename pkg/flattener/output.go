package flattener

// IntermediateContext is a completed Data object that has not been finalized into a data.
type IntermediateContext struct {
	// Mapping of Names to all builders in the current stack. This is the union of all output parents and children.
	// This does not contain a list of all siblings to the current data.
	// An ObjectBuilder may not be completed at the time this Context is first evaluated. It is important to check the
	// Status of the ObjectBuilder.
	Related map[string][]*OutputBuilder
	// The actual data of this output object.
	Data *OutputBuilder
	// Name of the object being output. This is not unique.
	Name string
}

// Finalize extracts the data from the OutputBuilder to build a final Output object.
func (ic *IntermediateContext) Finalize() *Output {
	return &Output{
		Name: ic.Name,
		Data: ic.Data.Value(),
	}
}

// Output is the final output object of a flattener
type Output struct {
	// Name of the output object. This is not unique.
	Name string
	// Raw data of the output object
	Data interface{}
}
