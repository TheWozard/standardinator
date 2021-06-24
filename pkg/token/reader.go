package token

// ReaderType string representation of different reader types
type ReaderType string

const (
	JSONReader ReaderType = "json"
	// Stretch goals
	// XMLReader  ReaderType = "xml"
	// YamlReader ReaderType = "yaml"
	// CSVReader  ReaderType = "csv"
)

// Reader produces a stream of tokens. Returns io.EOF when the token stream has ended
type Reader interface {
	Next() (Token, error)
}
