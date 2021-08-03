package extractor_test

import (
	"TheWozard/standardinator/pkg/extractor"
	"bytes"
	"io"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiExtractor(t *testing.T) {
	t.Skip()
	data := `{}`
	file := bytes.NewBufferString(data)
	output := []map[string]interface{}{{"A": 1.0}, {"B": 2.0}, {"C": 3.0}, {"D": 4.0}}

	decoder := extractor.Multi{
		Decoders: []extractor.Decoder{
			extractor.JSON{},
			extractor.JSON{},
		},
	}

	extractor := decoder.New(file)
	precount := runtime.NumGoroutine()
	for _, expected := range output {
		actual, err := extractor.Next()
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
	actual, err := extractor.Next()
	require.Nil(t, actual, "Unexpected Extra Results")
	require.Equal(t, io.EOF, err)
	require.Equal(t, precount, runtime.NumGoroutine())
}
