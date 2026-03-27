package errors

import (
	stderrors "errors"
	"net/http"
)

// Translator maps sentinel errors to HTTP status codes.
// Entries are evaluated in registration order using errors.Is,
// so wrapped errors are correctly resolved.
type Translator struct {
	entries []entry
}

type entry struct {
	sentinel error
	status   int
}

// NewTranslator returns an empty Translator with no registered mappings.
func NewTranslator() *Translator {
	return &Translator{}
}

// Register adds a mapping from a sentinel error to an HTTP status code.
// Returns the Translator itself to allow method chaining.
func (t *Translator) Register(sentinel error, status int) *Translator {
	t.entries = append(t.entries, entry{sentinel, status})
	return t
}

// Translate returns the HTTP status code for the given error.
// Uses errors.Is internally, so wrapped errors are supported.
// Returns 500 if no mapping is found or if err is nil.
func (t *Translator) Translate(err error) int {
	for _, e := range t.entries {
		if stderrors.Is(err, e.sentinel) {
			return e.status
		}
	}
	return http.StatusInternalServerError
}
