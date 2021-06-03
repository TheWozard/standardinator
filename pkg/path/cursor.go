package path

// Cursor defines a current location in the resolving of a path
type Cursor interface {
	// IsMatched returns if the current cursor has reached the end of the path and is completed
	IsMatched() bool
	// Match attempts to move the cursor forward by matching an element returning nil if no match is found and a new Cursor is the event it is matched
	Match(target string) Cursor
}
