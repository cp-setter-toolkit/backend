package language

import "github.com/thepluck/cp-setter-toolkit/helper/errors"

// Registry is a registry of programming languages.
type Registry struct {
	Languages map[string]Language
}

func NewRegistry() *Registry {
	return &Registry{Languages: make(map[string]Language)}
}

func (r *Registry) Register(l Language) {
	r.Languages[l.Name()] = l
}

func (r *Registry) Get(name string) (Language, error) {
	lang, ok := r.Languages[name]
	if !ok {
		return nil, errors.Errorf("language %s not found", name)
	}
	return lang, nil
}

var DefaultRegistry = NewRegistry()
