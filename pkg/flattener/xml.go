package flattener

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/PaesslerAG/gval"
)

// XML creates an extractor for xml data
type XML struct {
	Config XMLOutputConfig

	// DisableStrict if the default golang XML decoder strict option
	DisableStrict bool
	// Disables erroring on finding embedded HTML data
	HTMLEscape bool
}

type XMLOutputConfig []*XMLOutput

// XMLOutput defines how an output is formed the XML
type XMLOutput struct {
	Path    *Path
	Name    string
	Repeat  []*Path
	Objects []*Path
	// gval.Evaluable is an evaluate-able jsonpath created by jsonpath.New(path)
	Import map[string]map[string]gval.Evaluable
}

func (o *XMLOutput) GetName() string {
	if o.Name != "" {
		return o.Name
	}
	return o.Path.toString()
}

// New creates an Evaluator that reads XML data from r
func (c XML) New(reader io.Reader) Evaluator {
	decoder := xml.NewDecoder(reader)
	decoder.Strict = !c.DisableStrict
	//TODO: HTML Escape
	evaluations := map[string]*Evaluation{}
	for _, output := range c.Config {
		evaluations[output.GetName()] = &Evaluation{
			Import: output.Import,
		}
	}
	return &IntermediateEvaluator{
		Extractor: &xmlExtractor{
			next:   decoder.Token,
			config: c.Config,
			factory: OutputBuilderFactory{
				Strict: !c.DisableStrict,
			},
			path: NewPath("$"),
		},
		Evaluations: evaluations,
	}
}

// xmlExtractor reads tokens from a Token function of xml.Decoder
// Uses the tokens and a stack of contexts to build output objects
type xmlExtractor struct {
	next    func() (xml.Token, error)
	config  XMLOutputConfig
	factory OutputBuilderFactory

	// current execution stack
	context []xmlExtractorContext
	path    *Path
}

type xmlExtractorContext struct {
	builder *OutputBuilder
	config  *XMLOutput
	name    string
}

func (e *xmlExtractor) Next() (*IntermediateContext, error) {
	// Token Processing
	for {
		token, err := e.next()
		if err != nil {
			return nil, err
		}
		switch typed := token.(type) {
		case xml.StartElement:
			e.path.Enter(typed.Name.Local)
			// Do we need to open a new context for this data.
			if config := e.getConfig(); config != nil {
				builder := e.factory.New()
				e.context = append(e.context, xmlExtractorContext{
					name:    config.GetName(),
					config:  config,
					builder: builder,
				})
			}
			if len(e.context) > 0 {
				current := e.currentContext()
				name := typed.Name.Local
				if e.isArray() {
					current.builder.Array(name)
				} else if e.isObject() {
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
			e.path.Exit()
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
func (e *xmlExtractor) isArray() bool {
	for _, target := range e.currentContext().config.Repeat {
		if e.path.Includes(target) {
			return true
		}
	}
	return false
}

// isObject returns if the passed field is expect to be an object
func (e *xmlExtractor) isObject() bool {
	for _, target := range e.currentContext().config.Objects {
		if e.path.Includes(target) {
			return true
		}
	}
	return false
}

// getConfig attempts to return the first XMLOutput that matches the current path returns nil otherwise
func (e *xmlExtractor) getConfig() *XMLOutput {
	for _, config := range e.config {
		if e.path.Includes(config.Path) {
			return config
		}
	}
	return nil
}

// Close closes the current context and outputs with the completed data.
func (e *xmlExtractor) Close() (*IntermediateContext, error) {
	var current xmlExtractorContext
	e.context, current = e.context[:len(e.context)-1], e.context[len(e.context)-1]
	Related := map[string][]*OutputBuilder{}
	// TODO: Children
	for _, parent := range e.context {
		Related[parent.name] = append(Related[parent.name], parent.builder)
	}
	return &IntermediateContext{
		Name:    current.name,
		Data:    current.builder,
		Related: Related,
	}, nil
}
