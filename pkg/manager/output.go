package manager

import (
	"TheWozard/standardinator/pkg/token"
	"fmt"
)

func NewManager() Manager {
	return &outputManager{}
}

type outputManager struct {
	key  string
	data valueCollector

	matches []resolver // TODO: optimize

	collecting bool
	parent     *outputManager
}

func (m *outputManager) HasResult() bool {
	return false
}

func (m *outputManager) GetResult() *Result {
	return nil
}

func (m *outputManager) Receive(key string, value interface{}) error {
	// First add to something that this tree might be a part of
	m.data.AddValue(key, value)
	// Next try and resolve any paths we can. We are not making a new child so we should be resolving paths at this point
	for _, resolver := range m.matches {
		if matched := resolver.Cursor.Match(key); matched != nil {
			if !matched.IsMatched() {
				return fmt.Errorf("unexpected end to path %s at element %s", matched, key)
			}
			err := resolver.Collector.AddValue(key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *outputManager) CreateChildNode(t token.StartToken) (Manager, error) {

	matches := []resolver{}
	collecting := m.collecting
	for _, r := range m.matches {
		if matched := r.Cursor.Match(t.Key); matched != nil {
			if !matched.IsMatched() {
				matches = append(matches, resolver{
					Cursor:    matched,
					Collector: r.Collector,
				})
			} else {
				collecting = true
			}
		}
	}

	var data valueCollector
	if !collecting {
		data = newNoOpValue()
	} else if t.Array {
		data = newMultiValue()
	} else {
		data = newObjectValue()
	}

	return &outputManager{
		key:  t.Key,
		data: data,

		matches: matches,

		collecting: collecting,
		parent:     m,
	}, nil
}

func (m *outputManager) Flush() error {
	if m.parent == nil {
		return nil
	}
	value, err := m.data.GetValue()
	if err != nil {
		return err
	}
	m.parent.Receive(m.key, value)
	return nil
}
