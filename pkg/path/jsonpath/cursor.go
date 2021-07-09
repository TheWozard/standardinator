package jsonpath

import "TheWozard/standardinator/pkg/path"

// cursor is a path.Cursor implementation based on the index in the current jsonpath.parsed.path
type cursor struct {
	cursor int
	parent parsed
	value  path.ValueCollector
}

func (c cursor) AddValue(key string, value interface{}) error {
	return c.value.AddValue(key, value)
}

func (c cursor) GetValue() (interface{}, error) {
	return c.value.GetValue()
}

func (c cursor) IsMatched() bool {
	return c.cursor >= len(c.parent.path)
}

func (c cursor) Match(target string) path.Cursor {
	if c.IsMatched() {
		return nil // If we already matched we cant match anything new
	}
	if c.parent.path[c.cursor] == "" {
		// We are doing a deep search here so we are going to check based on the next element and when it matches move the cursor forward
		// otherwise we are going to keep returning out current cursor as its still kinda a match
		result := cursor{
			cursor: c.cursor + 1,
			parent: c.parent,
			value:  c.value,
		}.Match(target)
		if result != nil {
			return c
		} else {
			return result
		}
	}
	if c.parent.path[c.cursor] == target {
		// Standard matching situation, moving the cursor forward one
		return cursor{
			cursor: c.cursor + 1,
			parent: c.parent,
			value:  c.value,
		}
	}
	return nil
}
