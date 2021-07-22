package matcher

func NewSimple(target string) *Simple {
	return &Simple{
		Matched: target == ".",
		target:  target,
	}
}

// Simple provides the current implementation for matching a location
type Simple struct {
	Matched bool
	target  string
}

func (s *Simple) Check(value string) bool {
	s.Matched = (value == s.target)
	return s.Matched
}

func (s *Simple) Reset() {
	s.Matched = false
}
