package core

import (
	"TheWozard/standardinator/pkg/output"
)

type Iterator interface {
	//HasEnded returns if the reader has reached the end of the input stream. Does not garrantee that there is another result to be returned, only that the end has not been reached.
	HasNext() bool
	//Next returns the next completed object from the input stream or io.EOF in the event the iterator has reached the end.
	Next() (output.Result, error)
}
