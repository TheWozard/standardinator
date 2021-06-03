package jsonpath

import (
	"TheWozard/standardinator/pkg/path"
	"errors"
	"strings"
)

type parsed struct {
	path []string
	// TODO: If the path object is planned to be immutable, it would be faster to pre-compute functions like IsRoot and attach it to the struct
}

// Parse parses a string into a modified form of JSONPath
// TODO: define additions to the syntax and currently supported options
func Parse(path string) (path.Parsed, error) {
	if path == "" {
		return nil, errors.New("cannot parse empty path")
	}
	splitPath := strings.Split(path, ".")
	if !(splitPath[0] == "$" || strings.HasPrefix(splitPath[0], "@")) {
		return nil, errors.New("path must begin with either a $ or @ character")
	}
	return &parsed{
		path: splitPath,
	}, nil
}

func (p parsed) IsRoot() bool {
	return p.path[0] == "$"
}

func (p parsed) IsRelativeTo(element string) bool {
	return strings.HasPrefix(p.path[0], "@") && p.path[0][1:] == element
}

func (p parsed) Match(target string) path.Cursor {
	if p.IsRelativeTo(target) {
		return cursor{
			cursor: 1,
			parent: p,
		}
	}
	return nil
}

func (p parsed) MatchRoot() path.Cursor {
	if p.IsRoot() {
		return cursor{
			cursor: 1,
			parent: p,
		}
	}
	return nil
}
