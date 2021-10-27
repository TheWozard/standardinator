package extractor

import (
	"TheWozard/standardinator/pkg/data"
	"context"
	"io"
	"reflect"
)

// Multi allows multiple decoders to be run in parallel through a single Extractor
type Multi struct {
	Decoders []Decoder
}

// Creates an Extractor that splits all incoming data on r to all decoders
func (c Multi) New(r io.Reader) Extractor {
	extractors := make([]Extractor, len(c.Decoders))
	next := r
	for i, decoder := range c.Decoders {
		current := next
		if i < len(c.Decoders)-1 {
			var pw io.Writer
			// We use a pipe so the read/write opperations block and there isn't a buildup
			// of unread elements in memory
			next, pw = io.Pipe()
			current = io.TeeReader(current, pw)
		}
		extractors[i] = decoder.New(current)
	}
	return &multi{
		extractors: extractors,
	}
}

type response struct {
	payload *data.Payload
	err     error
}

// Coordinates multiple tread extraction for running extractors in parallel while
// keeping the overhead as small as possible.
// This process is done in multiple threads to prevent the possibility of loading a large amount of
// the document into memory. This would happen if we called each Extractor in sequence as the Next
// function always returns either and extracted object or an error and reads until it finds one
// All the read contents would have to sit in memory for the remaining extractors to process after the
// first one has extracted its object, which could be at the end of a large file.
type multi struct {
	extractors []Extractor
	get        func() (*data.Payload, error)
}

func (m *multi) Next() (*data.Payload, error) {
	// Once we get our first call we can short circuit
	if m.get != nil {
		return m.get()
	}
	if len(m.extractors) == 0 {
		return nil, io.EOF
	}
	// Each extractor is given its own channel so we can keep track of
	// how many have not yet completed
	chans := make([]chan response, len(m.extractors))
	ctx, cancel := context.WithCancel(context.Background())
	for i, extractor := range m.extractors {
		channel := make(chan response)
		chans[i] = channel
		e := extractor
		// Starting the thread for each extractor
		go func() {
			defer func() {
				close(channel)
			}()
			for {
				payload, err := e.Next()
				select {
				case channel <- response{payload, err}:
					if err != nil {
						return
					}
				case <-ctx.Done():
					// We need a way to ensure all threads close in the event
					// of an error
					return
				}
			}
		}()
	}
	cases := make([]reflect.SelectCase, len(chans))
	for i, ch := range chans {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	m.get = func() (*data.Payload, error) {
		for {
			// Once all of the channels are closed
			if len(cases) == 0 {
				cancel()
				return nil, io.EOF
			}
			chosen, value, ok := reflect.Select(cases)
			// if we get a channel closed event we take it out of the list.
			// The thread is the only one who should be closing the channel.
			if !ok {
				cases = append(cases[:chosen], cases[chosen+1:]...)
				continue
			}
			p := value.Interface().(response)
			if p.err != nil && p.err != io.EOF {
				// if we get a terminating error we should cancel all our goroutines
				cancel()
			}
			return p.payload, p.err
		}
	}
	return m.get()
}
