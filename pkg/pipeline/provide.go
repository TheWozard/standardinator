package pipeline

type ProvideStep struct {
	Data []map[string]interface{}
}

func (s ProvideStep) Init(prev GetNext) GetNext {
	index := 0
	return func() (map[string]interface{}, error) {
		if index >= len(s.Data) {
			return prev()
		}
		data := s.Data[index]
		index += 1
		return data, nil
	}
}

func (s ProvideStep) Close() error {
	return nil
}
