package extractor

import (
	"io"
	"reflect"
)

func Multi(extractors []Extractor) Extractor {
	return &multi{
		extractors: extractors,
	}
}

type data map[string]interface{}
type payload struct {
	data
	error
}

type multi struct {
	extractors []Extractor
	get        func() (map[string]interface{}, error)
}

func (m *multi) Next() (map[string]interface{}, error) {
	if m.get != nil {
		return m.get()
	}
	if len(m.extractors) == 0 {
		return nil, io.EOF
	}
	chans := make([]chan payload, len(m.extractors))
	for i, extractor := range m.extractors {
		channel := make(chan payload)
		chans[i] = channel
		go func() {
			for {
				rtn, err := extractor.Next()
				channel <- payload{rtn, err}
				if err != nil {
					close(channel)
					break
				}
			}
		}()
	}
	cases := make([]reflect.SelectCase, len(chans))
	for i, ch := range chans {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	m.get = func() (map[string]interface{}, error) {
		for {
			if len(cases) == 0 {
				return nil, io.EOF
			}
			chosen, value, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases[:chosen], cases[chosen+1:]...)
				continue
			}
			p := value.Interface().(payload)
			return p.data, p.error
		}
	}
	return m.get()
}
