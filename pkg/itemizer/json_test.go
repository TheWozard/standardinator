package itemizer_test

import (
	"TheWozard/standardinator/pkg/itemizer"
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Next(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		input string
	}{
		{
			name:  "example",
			path:  "$",
			input: "[{},{},{},{}]",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, test.name), func(t *testing.T) {
			reader, err := itemizer.CreateJSONReader(ioutil.NopCloser(bytes.NewBufferString(test.input)), test.path)
			require.NoError(t, err)
			for {
				_, err := reader.Next()
				require.NoError(t, err)
			}
		})
	}
}
