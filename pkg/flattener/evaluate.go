package flattener

import (
	"context"
	"fmt"
	"io"

	"github.com/PaesslerAG/gval"
)

// IntermediateEvaluator implements a generic Evaluator for handling IntermediateContext evaluation
type IntermediateEvaluator struct {
	// Map of rules for each object Name
	Evaluations map[string]*Evaluation

	// Extractor that will provide the *IntermediateContext
	Extractor Extractor
	// Storage for unresolved contexts. This should be empty before returning io.EOF
	Holding []*IntermediateContext
}

func (ie *IntermediateEvaluator) Next() (*Output, error) {
	for {
		// Check Holding for any now completed objects
		for i, context := range ie.Holding {
			eval, ok := ie.Evaluations[context.Name]
			if !ok {
				// We use a default Evaluation that makes no changes
				eval = &Evaluation{}
			}
			if eval.IsReady(context) {
				ie.Holding = append(ie.Holding[:i], ie.Holding[i+1:]...)
				return eval.Eval(context)
			}
		}

		// Get the next value for return
		extracted, err := ie.Extractor.Next()
		if err == io.EOF && len(ie.Holding) > 0 {
			return nil, fmt.Errorf("completed output before all objects were completed, %d object left uncompleted", len(ie.Holding))
		}
		if err != nil {
			return nil, err
		}
		ie.Holding = append(ie.Holding, extracted)
	}

}

// Evaluation applies data from OutputContext.Related data
type Evaluation struct {
	// Import a specific jsonpath to a new object
	Import map[string]map[string]gval.Evaluable
}

func (e *Evaluation) IsReady(context *IntermediateContext) bool {
	for name, _ := range e.Import {
		for _, builder := range context.Related[name] {
			if !builder.IsComplete() {
				return false
			}
		}
	}
	return true
}

func (e *Evaluation) Eval(inter *IntermediateContext) (*Output, error) {
	// If we add more rules on how to join data in a stack they should be added here
	for name, paths := range e.Import {
		builders := inter.Related[name]
		if len(builders) == 0 {
			return nil, fmt.Errorf("could not locate context information for name %s for %s", name, inter.Name)
		}
		if len(builders) > 1 {
			return nil, fmt.Errorf("ambiguous context information for name %s for %s found %d instances but expected 1", name, inter.Name, len(builders))
		}
		data := builders[0].Value()
		for key, path := range paths {
			value, err := path(context.Background(), data)
			if err != nil {
				// If we fail to get the data for any reason just use nil
				value = nil
			}
			err = inter.Data.Add(key, value)
			if err != nil {
				return nil, err
			}
		}
	}
	return inter.Finalize(), nil
}
