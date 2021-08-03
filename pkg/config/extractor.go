package config

import (
	"TheWozard/standardinator/pkg/extractor"
	"encoding/json"
	"fmt"
)

const (
	JSONExtractor ExtractorKind = "json"
	XMLExtractor  ExtractorKind = "xml"
)

type ExtractorKind = string

func GetDecoder(kind ExtractorKind, raw json.RawMessage) (extractor.Decoder, error) {
	switch kind {
	case JSONExtractor:
		config := extractor.JSON{}
		err := json.Unmarshal(raw, &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	case XMLExtractor:
		config := extractor.XML{}
		err := json.Unmarshal(raw, &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	default:
		return nil, fmt.Errorf("invalid extractor kind '%s'", kind)
	}
}
