package extractor

import (
	"encoding/xml"
	"fmt"
	"io"
)

func NewXML(config XmlExtractorConfig, reader io.Reader) Extractor {
	decoder := xml.NewDecoder(reader)
	decoder.Strict = config.Strict
	return &xmlExtractor{
		decoder: decoder,
		config:  config,
	}
}

type XmlExtractorConfig struct {
	Token string `json:"Token"`
	// TODO: implement arrays
	Repeats        []string `json:"Repeats"`
	Strict         bool     `json:"Strict"`
	IgnoreConflict bool     `json:"IgnoreConflict"`
}

type xmlExtractor struct {
	decoder *xml.Decoder
	config  XmlExtractorConfig
}

func (e *xmlExtractor) Next() (map[string]interface{}, error) {
	for {
		token, err := e.decoder.Token()
		if err != nil {
			return nil, err
		}
		switch typed := token.(type) {
		case xml.StartElement:
			if typed.Name.Local == e.config.Token {
				return e.decodeToMap(typed)
			}
		}
	}
}

type xmlDataContext struct {
	config XmlExtractorConfig
	stack  []*xmlContext
}

type xmlContext struct {
	data      map[string]interface{}
	namespace string
}

func (c *xmlDataContext) Add(key string, value interface{}) error {
	target := c.stack[len(c.stack)-1]
	if old, ok := target.data[key]; ok {
		if !c.config.IgnoreConflict {
			return fmt.Errorf("conflicting key '%s' has value '%v' and '%v'", key, old, value)
		}
	}
	target.data[key] = value
	return nil
}

func (c *xmlDataContext) Start(start xml.StartElement) error {
	data := map[string]interface{}{}
	c.stack = append(c.stack, &xmlContext{
		data:      data,
		namespace: start.Name.Local,
	})
	for _, att := range start.Attr {
		err := c.Add(fmt.Sprintf("@%s", att.Name.Local), att.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *xmlDataContext) End(end xml.EndElement) (map[string]interface{}, error) {
	var target *xmlContext
	c.stack, target = c.stack[:len(c.stack)-1], c.stack[len(c.stack)-1]
	if len(c.stack) == 0 {
		return target.data, nil
	}
	_, ok := target.data["#text"]
	if len(target.data) == 1 && ok {
		err := c.Add(target.namespace, target.data["#text"])
		if err != nil {
			return nil, err
		}
	} else {
		err := c.Add(target.namespace, target.data)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (e *xmlExtractor) decodeToMap(start xml.StartElement) (map[string]interface{}, error) {
	contextStack := &xmlDataContext{
		stack: []*xmlContext{},
	}
	contextStack.Start(start)
	for {
		token, err := e.decoder.Token()
		if err != nil {
			return nil, err
		}
		switch typed := token.(type) {
		case xml.StartElement:
			err := contextStack.Start(typed)
			if err != nil {
				return nil, err
			}
		case xml.EndElement:
			final, err := contextStack.End(typed)
			if err != nil {
				return nil, err
			}
			if final != nil {
				return final, nil
			}
		case xml.CharData:
			err := contextStack.Add("#text", string(typed))
			if err != nil {
				return nil, err
			}
		}
	}
}
