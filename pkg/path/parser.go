package path

// Parsed defines a path that is ready to be used to match during traversal.
type Parsed interface {
	// IsRoot returns if the path is relative to the root of the object.
	IsRoot() bool
	// IsRelativeTo returns if the path is relative to the passed element.
	IsRelativeTo(string) bool
	// MatchRoot attempts to start a cursor by matching the root. Returns nil on failed matches.
	MatchRoot() Cursor
	// Match attempts to start a cursor by matching the element. Returns nil on failed matches.
	Match(string) Cursor
}
