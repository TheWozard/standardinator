package pipeline

import "github.com/PaesslerAG/jsonpath"

type Select struct {
	Paths map[string]string `json:"paths"`
}

func (s Select) Init(prev GetNext) GetNext {
	return Wrap(prev, func(data map[string]interface{}) (map[string]interface{}, error) {
		new := map[string]interface{}{}
		for key, path := range s.Paths {
			data, err := jsonpath.Get(path, data)
			if err != nil {
				return nil, err
			}
			new[key] = data
		}
		return new, nil
	})
}

func (s Select) Close() error {
	return nil
}
