package core

import (
	"TheWozard/standardinator/pkg/output"
	"TheWozard/standardinator/pkg/token"
	"fmt"
	"io"
)

// NewIterator creates new Iterator with the provided Reader. This will read the reader and convert the tokens to configured output
func NewIterator(reader token.Reader) Iterator {
	return &synchronous{
		reader:  reader,
		manager: output.NewManager(),
	}
}

type synchronous struct {
	reader  token.Reader
	manager output.Manager

	closed bool
	child  *synchronous

	// TODO: config
}

func (s *synchronous) HasNext() bool {
	return s.manager.HasResult() || !s.closed
}

func (s *synchronous) Next() (output.Result, error) {
	// First we resolve all issues with children
	if s.child != nil {
		result, err := s.child.Next()
		if err != nil {
			if err == io.EOF {
				s.child = nil
			} else {
				return nil, err
			}
		} else {
			return result, nil
		}
	}
	// Next we resolve all of our own backlog
	if s.manager.HasResult() {
		return s.manager.GetResult(), nil
	}
	// Then we check to see if we have completed pulling new tokens
	if s.closed {
		// We have already closed this node in the tree and there is no backlogged results to be returned, time to actually close out this node.
		return nil, io.EOF
	}
	// Then back to pulling new tokens
	for {
		next, err := s.reader.Next()
		if err != nil {
			return nil, err
		}

		switch typed := next.(type) {
		case token.StartToken:
			s.spawnNewChild()
			// We call out own next as the child might not actually have anything in which we would have to continue our own next search
			return s.Next()
		case token.EndToken:
			// We done!
			return s.close()
		case token.KVToken:
			// TODO:
		default:
			return nil, fmt.Errorf("unexpected token %v", typed)
		}
	}
}

func (s *synchronous) spawnNewChild() {
	s.child = &synchronous{
		reader:  s.reader,
		manager: s.manager.CreateChildNode(),
	}
}

// close sets the iterator into a closed state that will attempt to return any remaining unfinished output children
func (s *synchronous) close() (output.Result, error) {
	s.closed = true
	s.manager.Flush()
	return s.manager.GetResult(), nil
}
