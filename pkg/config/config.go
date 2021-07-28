package config

import (
	"TheWozard/standardinator/pkg/extractor"
	"encoding/json"
	"io"
	"os"
)

func NewExtractionConfigFromFile(path string) (ExtractionConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err

	}
	defer file.Close()
	return NewExtractionConfig(file)
}

func NewExtractionConfig(r io.Reader) (ExtractionConfig, error) {
	config := ExtractionConfig{}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&config)
	return config, err
}

type ExtractionConfig []OutputConfig

func (configs ExtractionConfig) GetExtractor(reader io.Reader) (extractor.Extractor, error) {
	extractors := []extractor.Extractor{}
	test := reader
	for i, config := range configs {
		r := test
		if i != len(configs)-1 {
			var pw io.Writer
			r, pw = io.Pipe()
			test = io.TeeReader(test, pw)
		}
		extractor, err := config.GetExtractor(r)
		if err != nil {
			return nil, err
		}
		extractors = append(extractors, extractor)
	}
	return extractor.Multi(extractors), nil
}

type OutputConfig struct {
	Name            string          `json:"name"`
	Extractor       string          `json:"extractor"`
	ExtractorConfig json.RawMessage `json:"extractor_config"`
}

func (c *OutputConfig) GetExtractor(reader io.Reader) (extractor.Extractor, error) {
	return GetExtractor(c.Extractor, c.ExtractorConfig, reader)
}
