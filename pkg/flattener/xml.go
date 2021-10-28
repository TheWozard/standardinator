package flattener

import (
	"encoding/xml"
	"fmt"
	"io"
)

// XML creates an extractor for xml data
type XML struct {
	Elements XMLElementConfig
	// Map of Elements to output
	Output XMLOutputConfig

	// DisableStrict if the default golang XML decoder strict option
	DisableStrict bool
}

type XMLOutputConfig map[string]XMLOutput

// XMLOutput defines how an output is formed the XML
type XMLOutput struct {
	IncludeInParent bool
	Import          []string
	SubConfig       XMLOutputConfig
}

type XMLElementConfig struct {
	// Elements that are known to repeat inside of the element
	Arrays []string
	// Elements that are know to be objects and to not compress out #text if it exists
	Objects []string

	SubConfigs map[string]XMLElementConfig
}

// New creates an Extractor that reads XML data from r
func (c XML) New(reader io.Reader) Extractor {
	decoder := xml.NewDecoder(reader)
	decoder.Strict = !c.DisableStrict
	return &xmlExtractor{
		next:     decoder.Token,
		elements: c.Elements,
		output:   c.Output,
		factory: OutputBuilderFactory{
			Strict: !c.DisableStrict,
		},
	}
}

// xmlExtractor reads tokens from a Token function of xml.Decoder
// Uses the tokens and a stack of contexts to build output objects
type xmlExtractor struct {
	next     func() (xml.Token, error)
	elements XMLElementConfig
	output   XMLOutputConfig
	factory  OutputBuilderFactory

	// current execution stack
	context []xmlExtractorContext

	// Used for entering namespace specific configs
	subInstance *xmlExtractor
}

type xmlExtractorContext struct {
	builder *OutputBuilder
	config  XMLOutput
	name    string
}

func (e *xmlExtractor) Next() (*OutputContext, error) {
	// Evaluate any running subInstances
	if e.subInstance != nil {
		data, err := e.subInstance.Next()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if data != nil {
			return data, nil
		}
		// Release subInstance
		e.subInstance = nil
	}

	// Token Processing
	for {
		token, err := e.next()
		if err != nil {
			return nil, err
		}
		switch typed := token.(type) {
		case xml.StartElement:
			// Do we need to open a new context for this data.
			if config, ok := e.output[typed.Name.Local]; ok {
				if config.SubConfig != nil {
					e.subInstance = &xmlExtractor{
						next: func() (xml.Token, error) {
							// TODO: Wrapping this to get it back out
							return e.next()
						},
						elements: ,
					}
				}
				builder := e.factory.New()
				e.context = append(e.context, xmlExtractorContext{
					name:    typed.Name.Local,
					config:  config,
					builder: builder,
				})
			}
			if len(e.context) > 0 {
				current := e.currentContext()
				name := typed.Name.Local
				if e.isArray(name) {
					current.builder.Array(name)
				} else if e.isObject(name) {
					current.builder.Object(name)
				} else {
					current.builder.Namespace(name)
				}
				for _, att := range typed.Attr {
					e.currentContext().builder.Add(fmt.Sprintf("@%s", att.Name.Local), att.Value)
				}
			}
		case xml.CharData:
			if len(e.context) > 0 {
				// CompressableFieldName will be converted to a pure value in the event there is no other parts in the context
				// and we opened a namespace and not an object
				err := e.currentContext().builder.Add(CompressableFieldName, string(typed))
				if err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			if len(e.context) > 0 {
				current := e.currentContext()
				err = current.builder.Close()
				if err != nil {
					return nil, err
				}
				if current.builder.IsComplete() {
					return e.Close()
				}
			}
		}
	}
}

// currentContext returns the currently building context
func (e *xmlExtractor) currentContext() xmlExtractorContext {
	return e.context[len(e.context)-1]
}

// isArray returns if the passed field is expected to be an array
func (e *xmlExtractor) isArray(field string) bool {
	for _, key := range e.elements.Arrays {
		if field == key {
			return true
		}
	}
	return false
}

// isObject returns if the passed field is expect to be an object
func (e *xmlExtractor) isObject(field string) bool {
	for _, key := range e.elements.Objects {
		if field == key {
			return true
		}
	}
	return false
}

// Close closes the current context and outputs with the completed data.
func (e *xmlExtractor) Close() (*OutputContext, error) {
	var current xmlExtractorContext
	e.context, current = e.context[:len(e.context)-1], e.context[len(e.context)-1]
	if current.config.IncludeInParent {
		var err error
		if e.isArray(current.name) {
			err = e.currentContext().builder.Merge(current.name, current.builder.Value())
		} else {
			err = e.currentContext().builder.Add(current.name, current.builder.Value())
		}
		if err != nil {
			return nil, err
		}
	}
	Related := map[string][]*OutputBuilder{}
	for _, parent := range e.context {
		Related[parent.name] = append(Related[parent.name], parent.builder)
	}
	return &OutputContext{
		Name:    current.name,
		Data:    current.builder,
		Related: Related,
	}, nil
}
