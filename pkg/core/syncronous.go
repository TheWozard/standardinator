package core

import (
	"TheWozard/standardinator/pkg/output"
	"TheWozard/standardinator/pkg/token"
)

type synchronous struct {
	reader  token.Reader
	manager output.Manager

	closed bool
}

func (s *synchronous) HasNext() bool {
	return s.manager.HasResult() || !s.closed
}

func (s *synchronous) Next() (*output.Result, error) {
	if s.manager.HasResult() {
		return s.manager.GetResult(), nil
	}
	if s.closed {
		// We have already closed this node in the tree and there is no backlogged results to be returned, time to actually close out this node.
		return nil, nil
	}
	token, err := s.reader.Next()
	if err != nil {
		return nil, err
	}

	return s.close()
}

// close sets the iterator into a closed state that will attempt to return any remaining unfinished output children
func (s *synchronous) close() (*output.Result, error) {
	s.closed = true
	s.manager.Flush()
	return s.manager.GetResult(), nil
}
