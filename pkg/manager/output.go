package manager

// func NewManager() Manager {
// 	return createOutputManager("$", []resolver{}, false, false)
// }

// func createOutputManager(key string, matches []resolver, collect bool, array bool) Manager {

// 	collecting := collect
// 	newMatches := []resolver{}
// 	for _, r := range matches {
// 		if matched := r.Cursor.Match(key); matched != nil {
// 			if !matched.IsMatched() {
// 				newMatches = append(newMatches, resolver{
// 					Cursor:    matched,
// 					Collector: r.Collector,
// 				})
// 			} else {
// 				collecting = true
// 			}
// 		}
// 	}

// 	var data valueCollector
// 	if !collecting {
// 		data = newNoOpValue()
// 	} else if array {
// 		data = newMultiValue()
// 	} else {
// 		data = newObjectValue()
// 	}

// 	return &outputManager{
// 		key:  key,
// 		data: data,

// 		matches: newMatches,
// 	}
// }

// type outputManager struct {
// 	key  string
// 	data valueCollector

// 	matches []resolver // TODO: optimize

// 	collecting bool
// 	parent     *outputManager
// }

// func (m *outputManager) HasResult() bool {
// 	return false
// }

// func (m *outputManager) GetResult() *Result {
// 	return nil
// }

// func (m *outputManager) Receive(key string, value interface{}) error {
// 	// First add to something that this tree might be a part of
// 	m.data.AddValue(key, value)
// 	// Next try and resolve any paths we can. We are not making a new child so we should be resolving paths at this point
// 	for _, resolver := range m.matches {
// 		if matched := resolver.Cursor.Match(key); matched != nil {
// 			if !matched.IsMatched() {
// 				return fmt.Errorf("unexpected end to path %s at element %s", matched, key)
// 			}
// 			err := resolver.Collector.AddValue(key, value)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func (m *outputManager) CreateChildNode(t token.StartToken) (Manager, error) {
// 	return createOutputManager(t.Key, m.matches, m.collecting, t.Array), nil
// }

// func (m *outputManager) Flush() error {
// 	if m.parent == nil {
// 		return nil
// 	}
// 	value, err := m.data.GetValue()
// 	if err != nil {
// 		return err
// 	}
// 	m.parent.Receive(m.key, value)
// 	return nil
// }
