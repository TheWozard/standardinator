package token

type Reader interface {
	Next() (Token, error)
}
