package output

// Manager is responsible for keeping track of the partial and completed output objects
type Manager interface {
	// HasResult returns if there is a result available to be output
	HasResult() bool
	// GetResult provides the next result to be output. Will only output one result even in the event of multiple backlogged
	GetResult() *Result

	// Receive accepts a key value pair and passes it to all needed locations. If results are completed they are available as part of the GetResult() call
	Receive(key string, value interface{})

	// CreateChildNode creates a new child node with its own object tracking for locally available objects
	CreateChildNode() Manager

	// Flush forces all open objects to be closed and output, even in a partially complete state.
	Flush()
}
