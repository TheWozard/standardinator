package flattener

import (
	"fmt"
)

const (
	CompressableFieldName = "#text"
	RootObjectName        = "$Root"
)

// Factory for setting the config for OutputBuilder and building new ones.
type OutputBuilderFactory struct {
	Strict bool
}

// New creates a new empty output object with a value of null.
// To building data remember to call Object or Array to pick starting node.
func (f OutputBuilderFactory) New() *OutputBuilder {
	builder := &OutputBuilder{
		stack: []*collector{},
		value: nil,

		strict: f.Strict,
	}
	return builder
}

// NewObject equivalent to calling New followed by Object to start a new Object
func (f OutputBuilderFactory) NewObject() *OutputBuilder {
	builder := f.New()
	builder.Object(RootObjectName)
	return builder
}

// NewArray equivalent to calling New followed by Array to start a new Array
func (f OutputBuilderFactory) NewArray() *OutputBuilder {
	builder := f.New()
	builder.Array(RootObjectName)
	return builder
}

// NewNamespace equivalent to calling New followed by Namespace to start a new Namespace
func (f OutputBuilderFactory) NewNamespace() *OutputBuilder {
	builder := f.New()
	builder.Namespace(RootObjectName)
	return builder
}

// OutputBuilder is a representation of an object that is in the process of being created.
// The current state of the object might be incomplete so checking isComplete is recommended
// before pulling the value.
type OutputBuilder struct {
	stack []*collector
	value interface{}

	strict bool
}

func (builder *OutputBuilder) context() *collector {
	return builder.stack[len(builder.stack)-1]
}

// Add adds the key value into the object at the current context. Errors if key already exists ins strict mode.
func (builder *OutputBuilder) Add(key string, value interface{}) error {
	if builder.IsComplete() {
		builder.value = value
		return nil
	}
	if builder.strict {
		return builder.context().Add(key, value)
	}
	builder.context().Set(key, value)
	return nil
}

// Merge in a key adding it as part of an array. Appends the previous value of the field.
func (builder *OutputBuilder) Merge(key string, value interface{}) error {
	if builder.IsComplete() {
		builder.value = value
		return nil
	}
	return builder.context().Merge(key, value)
}

// Object creates a namespace for the object and moves add operations to add to it.
func (builder *OutputBuilder) Object(namespace string) {
	data := map[string]interface{}{}
	builder.stack = append(builder.stack, &collector{
		data:            data,
		namespace:       namespace,
		maintainRawData: true,
	})
}

// Array creates a namespace for an array element and moves add operations to add to it.
// It is important that this is not the array as a whole but a single element of it.
// When this collector is merged in it attaches itself through an array, so one does not need to be
// created explicitly.
func (builder *OutputBuilder) Array(namespace string) {
	data := map[string]interface{}{}
	builder.stack = append(builder.stack, &collector{
		data:      data,
		namespace: namespace,
		repeats:   true,
	})
}

// Namespace creates a generic namespace. Similar to Object but this will attempt to compact data.
// This will occur when the only key of a namespace is the CompressableFieldName.
func (builder *OutputBuilder) Namespace(namespace string) {
	data := map[string]interface{}{}
	builder.stack = append(builder.stack, &collector{
		data:      data,
		namespace: namespace,
	})
}

// Close closes the current collection target and passses the data up.
// To see if the object is completed call IsComplete after calling
func (builder *OutputBuilder) Close() error {
	if builder.IsComplete() {
		return AlreadyCompletedObjectError
	}
	var target *collector
	builder.stack, target = builder.stack[:len(builder.stack)-1], builder.stack[len(builder.stack)-1]
	return target.AddTo(builder)
}

// IsComplete if the builder has a completed output object
func (builder *OutputBuilder) IsComplete() bool {
	return len(builder.stack) == 0
}

// Value if the value of the completed builder
func (builder *OutputBuilder) Value() interface{} {
	return builder.value
}

// collector represents a namespace of OutputBuilder and can be stacked to produce deeply nested data
type collector struct {
	// raw data of the collector, this may be modified when the field is passed up
	// notable a collector of just the field "#text" will be converted to its value instead of an object
	data      map[string]interface{}
	namespace string

	// Helps infrom how the collector should merge its data into the provided builder
	repeats bool
	// if the "#text" should be maintained
	maintainRawData bool
}

// Data returns the final output data of a collector
func (nc *collector) Data() interface{} {
	text, ok := nc.data[CompressableFieldName]
	if len(nc.data) == 1 && ok && !nc.maintainRawData {
		// If the data only contains "#text" then return it as a single value
		return text
	} else {
		return nc.data
	}
}

// AddTo has a collector add its data into an OutputBuilder. Depending of if repeats is true it
// will try to add the data into the builder as an array field.
func (nc *collector) AddTo(builder *OutputBuilder) error {
	if nc.repeats {
		return builder.Merge(nc.namespace, nc.Data())
	}
	return builder.Add(nc.namespace, nc.Data())
}

// Add is a strict version of Set that requires the key to be new and never seen before.
func (nc *collector) Add(key string, value interface{}) error {
	if old, ok := nc.data[key]; ok {
		return fmt.Errorf("conflicting key '%s' has value '%v' and '%v'", key, old, value)
	}
	nc.Set(key, value)
	return nil
}

// Set sets the key and value for the data regardles of it has previously been set.
func (nc *collector) Set(key string, value interface{}) {
	nc.data[key] = value
}

// Merge adds a value to a field that is expected to be an array.
// This means adding it as a single element array if the field does not exist or
// appending to the previous value
func (nc *collector) Merge(key string, value interface{}) error {
	var prev []interface{}
	if old, ok := nc.data[key]; ok {
		switch typed := old.(type) {
		case []interface{}:
			prev = typed
		default:
			return fmt.Errorf("unexpected value for repeating key '%s' %v", key, typed)
		}
	} else {
		prev = []interface{}{}
	}
	nc.Set(key, append(prev, value))
	return nil
}
