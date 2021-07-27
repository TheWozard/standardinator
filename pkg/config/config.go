package config

import (
	"TheWozard/standardinator/pkg/extractor"
	"encoding/json"
	"io"
)

type ExtractionConfig = []OutputConfig

type OutputConfig struct {
	Name            string          `json:"name"`
	Extractor       string          `json:"extractor"`
	ExtractorConfig json.RawMessage `json:"extractor_config"`
}

func (c *OutputConfig) GetExtractor(reader io.Reader) (extractor.Extractor, error) {
	return GetExtractor(c.Extractor, c.ExtractorConfig, reader)
}
