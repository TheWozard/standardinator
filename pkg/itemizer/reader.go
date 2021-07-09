package itemizer

type Reader interface {
	Next() (interface{}, error)
}
