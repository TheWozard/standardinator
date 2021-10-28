package flattener

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
)

type OutputContext struct {
	// Mapping of Names to all builders in the current stack. This is the union of all output parents and children.
	// This does not contain a list of all siblings to the current data
	// An ObjectBuilder may not be completed at the time this Context is first evaluated. It is important to check the
	// Status of the ObjectBuilder
	Related map[string][]*OutputBuilder
	// The actual data of this output object
	Data *OutputBuilder

	Name string
}

type OutputEvaluator struct {
	Import []ImportDetails
}

type ImportDetails struct {
	Name    string
	Mapping map[string]string
}

func (e *OutputEvaluator) IsReady(context *OutputContext) bool {
	for _, details := range e.Import {
		for _, builder := range context.Related[details.Name] {
			if !builder.IsComplete() {
				return false
			}
		}
	}
	return true
}

func (e *OutputEvaluator) Eval(context *OutputContext) interface{} {
	for _, details := range e.Import {
		builders := context.Related[details.Name]
		if len(builders) == 0 {
			return fmt.Errorf("Could not locate context information for name %s for %s", details.Name, context.Name)
		}
		if len(builders) > 1 {
			return fmt.Errorf("Ambiguous context information for name %s for %s found %d instances but expected 1", details.Name, context.Name, len(builders))
		}
		builder := builders[0]
		for key, path := range details.Mapping {
			value, err := jsonpath.Get(path, builder.Value())
			if err != nil {
				value = nil
			}
			context.Data.Add(key, value)
		}
	}
	return context.Data.Value()
}
