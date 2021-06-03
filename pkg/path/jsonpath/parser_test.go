package jsonpath_test

import (
	"errors"
	"fmt"
	"testing"

	"TheWozard/standardinator/pkg/path"
	"TheWozard/standardinator/pkg/path/jsonpath"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		err      error
		validate func(*testing.T, path.Parsed)
	}{
		{
			name:  "golden path root",
			input: "$.Destination",
			err:   nil,
			validate: func(t *testing.T, parsed path.Parsed) {
				require.True(t, parsed.IsRoot())
			},
		},
		{
			name:  "golden path relative",
			input: "@Output.Destination",
			err:   nil,
			validate: func(t *testing.T, parsed path.Parsed) {
				require.True(t, parsed.IsRelativeTo("Output"))
			},
		},
		{
			name:     "empty path",
			input:    "",
			err:      errors.New("cannot parse empty path"),
			validate: func(t *testing.T, parsed path.Parsed) {},
		},
		{
			name:     "no start symbol",
			input:    "Destination.Bad",
			err:      errors.New("path must begin with either a $ or @ character"),
			validate: func(t *testing.T, parsed path.Parsed) {},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, test.name), func(t *testing.T) {
			parsed, err := jsonpath.Parse(test.input)
			require.Equal(t, test.err, err)
			test.validate(t, parsed)
		})
	}
}

func TestIsRoot(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output bool
	}{
		{
			name:   "golden path",
			input:  "$.Destination",
			output: true,
		},
		{
			name:   "relative start",
			input:  "@.Destination",
			output: false,
		},
		{
			name:   "relative named start",
			input:  "@Output.Destination",
			output: false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, test.name), func(t *testing.T) {
			parsed, err := jsonpath.Parse(test.input)
			require.NoError(t, err)
			require.Equal(t, test.output, parsed.IsRoot())
		})
	}
}

func TestIsRelativeTo(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		target string
		output bool
	}{
		{
			name:   "golden path",
			input:  "@Output.Destination",
			target: "Output",
			output: true,
		},
		{
			name:   "bad location",
			input:  "@Output.Destination",
			target: "BadLocation",
			output: false,
		},
		{
			name:   "root path",
			input:  "$.Destination",
			target: "Output",
			output: false,
		},
		{
			name:   "partial location",
			input:  "@Output.Destination",
			target: "Out",
			output: false,
		},
		{
			name:   "empty location",
			input:  "@Output.Destination",
			target: "",
			output: false,
		},
		{
			name:   "longer location",
			input:  "@Output.Destination",
			target: "OutputLocation",
			output: false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, test.name), func(t *testing.T) {
			parsed, err := jsonpath.Parse(test.input)
			require.NoError(t, err)
			require.Equal(t, test.output, parsed.IsRelativeTo(test.target))
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		path     []string
		succeeds bool
	}{
		{
			name:     "golden path",
			input:    "@Output.Destination.Element.Type",
			path:     []string{"Output", "Destination", "Element", "Type"},
			succeeds: true,
		},
		{
			name:     "incomplete path traversal",
			input:    "@Output.Destination",
			path:     []string{"Output"},
			succeeds: false,
		},
		{
			name:     "invalid root",
			input:    "@Output.Destination",
			path:     []string{"Bad"},
			succeeds: false,
		},
		{
			name:     "invalid path",
			input:    "@Output.Destination",
			path:     []string{"Output", "Bad"},
			succeeds: false,
		},
		{
			name:     "overrun path",
			input:    "@Output.Destination",
			path:     []string{"Output", "Destination", "Bad"},
			succeeds: false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, test.name), func(t *testing.T) {
			parsed, err := jsonpath.Parse(test.input)
			require.NoError(t, err)
			match := parsed.Match(test.path[0])
			for i, path := range test.path[1:] {
				if match == nil {
					require.False(t, test.succeeds, fmt.Sprintf("path unexpectedly failed matching at index %d for key %s", i, test.path[i]))
					return
				}
				match = match.Match(path)
			}
			if match == nil {
				require.False(t, test.succeeds, "match was not returned at end of path that was intended to succeed")
				return
			}
			if test.succeeds {
				require.True(t, match.IsMatched(), "path was unexpectedly incomplete at the end of a test that was supposed to succeed")
			} else {
				require.False(t, match.IsMatched(), "path was unexpectedly completed at the end of a test that was supposed to fail")
			}
		})
	}
}

func TestRoot(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		path     []string
		succeeds bool
	}{
		{
			name:     "golden path",
			input:    "@Output.Destination.Element.Type",
			path:     []string{"Output", "Destination", "Element", "Type"},
			succeeds: true,
		},
		{
			name:     "incomplete path traversal",
			input:    "@Output.Destination",
			path:     []string{"Output"},
			succeeds: false,
		},
		{
			name:     "invalid root",
			input:    "@Output.Destination",
			path:     []string{"Bad"},
			succeeds: false,
		},
		{
			name:     "invalid path",
			input:    "@Output.Destination",
			path:     []string{"Output", "Bad"},
			succeeds: false,
		},
		{
			name:     "overrun path",
			input:    "@Output.Destination",
			path:     []string{"Output", "Destination", "Bad"},
			succeeds: false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, test.name), func(t *testing.T) {
			parsed, err := jsonpath.Parse(test.input)
			require.NoError(t, err)
			match := parsed.Match(test.path[0])
			for i, path := range test.path[1:] {
				if match == nil {
					require.False(t, test.succeeds, fmt.Sprintf("path unexpectedly failed matching at index %d for key %s", i, test.path[i]))
					return
				}
				match = match.Match(path)
			}
			if match == nil {
				require.False(t, test.succeeds, "match was not returned at end of path that was intended to succeed")
				return
			}
			if test.succeeds {
				require.True(t, match.IsMatched(), "path was unexpectedly incomplete at the end of a test that was supposed to succeed")
			} else {
				require.False(t, match.IsMatched(), "path was unexpectedly completed at the end of a test that was supposed to fail")
			}
		})
	}
}
