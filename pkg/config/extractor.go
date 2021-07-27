package config

import (
	"TheWozard/standardinator/pkg/extractor"
	"encoding/json"
	"fmt"
	"io"
)

const (
	JSONExtractor ExtractorKind = "json"
	XMLExtractor  ExtractorKind = "xml"
)

type ExtractorKind = string

func GetExtractor(kind ExtractorKind, raw json.RawMessage, reader io.Reader) (extractor.Extractor, error) {
	switch kind {
	case JSONExtractor:
		config := extractor.JSONConfig{}
		err := json.Unmarshal(raw, config)
		if err != nil {
			return nil, err
		}
		return extractor.NewJsonExtractor(config, reader), nil
	case XMLExtractor:
		config := extractor.XmlConfig{}
		err := json.Unmarshal(raw, config)
		if err != nil {
			return nil, err
		}
		return extractor.NewXMLExtractor(config, reader), nil
	default:
		return nil, fmt.Errorf("invalid extractor kind '%s'", kind)
	}
}
