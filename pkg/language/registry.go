package language

import "errors"

var ErrorLanguageNotFound = errors.New("language not found")

type NotFoundError struct {
	Id string
}

func (err NotFoundError) Error() string {
	return "language not found: " + err.Id
}

// MapRegistry is a implementation of Registry using a map.
type MapRegistry struct {
	Languages map[string]Language
}

func NewMapRegistry() *MapRegistry {
	return &MapRegistry{Languages: make(map[string]Language)}
}

func (r *MapRegistry) Register(lang Language) {
	r.Languages[lang.Id()] = lang
}

func (r *MapRegistry) Get(id string) (Language, error) {
	lang, ok := r.Languages[id]
	if !ok {
		return nil, NotFoundError{Id: id}
	}
	return lang, nil
}

var DefaultRegistry Registry

func init() {
	DefaultRegistry = NewMapRegistry()
}
