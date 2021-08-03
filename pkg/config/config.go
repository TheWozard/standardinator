package config

import (
	"TheWozard/standardinator/pkg/extractor"
	"encoding/json"
	"io"
	"os"
)

func NewConfigFromFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err

	}
	defer file.Close()
	return NewConfig(file)
}

func NewConfig(r io.Reader) (Config, error) {
	config := Config{}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&config)
	return config, err
}

type Config []OutputConfig

func (configs Config) GetDecoder() (extractor.Decoder, error) {
	decoders := make([]extractor.Decoder, len(configs))
	for i, config := range configs {
		decoder, err := config.GetDecoder()
		if err != nil {
			return nil, err
		}
		decoders[i] = decoder
	}
	return extractor.Multi{Decoders: decoders}, nil
}

type OutputConfig struct {
	Name            string          `json:"name"`
	Extractor       string          `json:"extractor"`
	ExtractorConfig json.RawMessage `json:"extractor_config"`
}

func (c *OutputConfig) GetDecoder() (extractor.Decoder, error) {
	return GetDecoder(c.Extractor, c.ExtractorConfig)
}
