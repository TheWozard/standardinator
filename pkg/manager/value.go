package manager

import (
	"TheWozard/standardinator/pkg/path"
	"errors"
)

// resolver is the link between a cursor and the value it resolves to. In order to take advantage of the immutability
// of the cursors while having a static value that once resolves persists up the tree
type resolver struct {
	Cursor    path.Cursor
	Collector valueCollector
}

type valueCollector interface {
	AddValue(key string, value interface{}) error
	GetValue() (interface{}, error)
}

// When a value should only resolve to a single value
func newSingleValue() valueCollector {
	return &singleValue{}
}

type singleValue struct {
	value interface{}
}

func (v *singleValue) AddValue(key string, value interface{}) error {
	if v.value != nil {
		return errors.New("can not accept multiple values")
	}
	v.value = value
	return nil
}

func (v *singleValue) GetValue() (interface{}, error) {
	if v.value == nil {
		return nil, errors.New("no value found")
	}
	return v.value, nil
}

// When a value should resolve to a list of values
func newMultiValue() valueCollector {
	return &multiValue{
		value: []interface{}{},
	}
}

type multiValue struct {
	value []interface{}
}

func (v *multiValue) AddValue(key string, value interface{}) error {
	v.value = append(v.value, value)
	return nil
}

func (v *multiValue) GetValue() (interface{}, error) {
	return v.value, nil
}

// When a value should only resolve to a object
func newObjectValue() valueCollector {
	return &objectValue{
		value: map[string]interface{}{},
	}
}

type objectValue struct {
	value map[string]interface{}
}

func (v *objectValue) AddValue(key string, value interface{}) error {
	v.value[key] = value
	return nil
}

func (v *objectValue) GetValue() (interface{}, error) {
	return v.value, nil
}

// When a value should never resolve
func newNoOpValue() valueCollector {
	return &noOpValue{}
}

type noOpValue struct{}

func (v *noOpValue) AddValue(key string, value interface{}) error {
	return nil
}

func (v *noOpValue) GetValue() (interface{}, error) {
	return nil, errors.New("attempted to retrive value of noOpValue")
}
