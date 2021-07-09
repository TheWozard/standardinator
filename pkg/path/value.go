package path

import (
	"errors"
)

type ValueCollector interface {
	AddValue(key string, value interface{}) error
	GetValue() (interface{}, error)
}

// When a value should only resolve to a single value
func NewSingleValue() ValueCollector {
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
func NewMultiValue() ValueCollector {
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
func NewObjectValue() ValueCollector {
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
func NewNoOpValue() ValueCollector {
	return &noOpValue{}
}

type noOpValue struct{}

func (v *noOpValue) AddValue(key string, value interface{}) error {
	return nil
}

func (v *noOpValue) GetValue() (interface{}, error) {
	return nil, errors.New("attempted to retrive value of noOpValue")
}
