package extractor

import (
	"encoding/xml"
	"fmt"
	"io"
)

// XmlExtractorConfig  defines configuration for a XML based Extractor
type XmlExtractorConfig struct {
	// Token the target token to output element of
	Token string `json:"token"`
	// Elements that are known to repeat inside of the element
	Repeats []string `json:"repeats"`
	// DisableStrict if the default golang XML decoder strict option
	DisableStrict bool `json:"disable_strict"`
}

// NewXMLExtractor creates an Extractor that reads XML data from r
func NewXMLExtractor(config XmlExtractorConfig, reader io.Reader) Extractor {
	decoder := xml.NewDecoder(reader)
	decoder.Strict = !config.DisableStrict
	return &xmlExtractor{
		decoder: decoder,
		config:  config,
	}
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

// xmlDataContext defines the overall element extraction context
type xmlDataContext struct {
	stack        []*xmlContext
	repeatLookup map[string]struct{}
}

type xmlContext struct {
	data      map[string]interface{}
	namespace string
}

// Add adds the past key and value to the current stack context
func (c *xmlDataContext) Add(key string, value interface{}) error {
	target := c.stack[len(c.stack)-1]
	if _, ok := c.repeatLookup[key]; ok {
		var prev []interface{}
		if old, ok := target.data[key]; ok {
			switch typed := old.(type) {
			case []interface{}:
				prev = typed
			default:
				return fmt.Errorf("unexpected value for repeating key '%s' %v", key, typed)
			}
		} else {
			prev = []interface{}{}
		}
		target.data[key] = append(prev, value)
		return nil
	}
	if old, ok := target.data[key]; ok {
		return fmt.Errorf("conflicting key '%s' has value '%v' and '%v'", key, old, value)
	}
	target.data[key] = value
	return nil
}

// Start creates a new parent context to the current stack based on the passed start element
func (c *xmlDataContext) Start(start xml.StartElement) error {
	data := map[string]interface{}{}
	c.stack = append(c.stack, &xmlContext{
		data:      data,
		namespace: start.Name.Local,
	})
	for _, att := range start.Attr {
		// Add all the elements attributes
		err := c.Add(fmt.Sprintf("@%s", att.Name.Local), att.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// End ends the current stack context and provides it data up the stack. If the stack would be empty the data is returned.
func (c *xmlDataContext) End(end xml.EndElement) (map[string]interface{}, error) {
	var target *xmlContext
	c.stack, target = c.stack[:len(c.stack)-1], c.stack[len(c.stack)-1]
	if len(c.stack) == 0 {
		// If the stack would be empty, return the data
		return target.data, nil
	}
	_, ok := target.data["#text"]
	if len(target.data) == 1 && ok {
		// If the data only contains #text then return it as a single value
		err := c.Add(target.namespace, target.data["#text"])
		if err != nil {
			return nil, err
		}
	} else {
		// Else return the completed data context as is
		err := c.Add(target.namespace, target.data)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// decodeToMap continues to pull tokens and decode them to a map
func (e *xmlExtractor) decodeToMap(start xml.StartElement) (map[string]interface{}, error) {
	lookup := map[string]struct{}{}
	for _, key := range e.config.Repeats {
		lookup[key] = struct{}{}
	}
	contextStack := &xmlDataContext{
		stack:        []*xmlContext{},
		repeatLookup: lookup,
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
			// in the event we get a final back we have received the completed map
			if err != nil {
				return nil, err
			}
			if final != nil {
				return final, nil
			}
		case xml.CharData:
			// #text will be converted to a pure value in the event there is no other parts in the context
			err := contextStack.Add("#text", string(typed))
			if err != nil {
				return nil, err
			}
		}
	}
}
