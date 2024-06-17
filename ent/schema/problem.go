package schema

import "entgo.io/ent"

// Problem holds the schema definition for the Problem entity.
type Problem struct {
	ent.Schema
}

// Fields of the Problem.
func (Problem) Fields() []ent.Field {
	return nil
}

// Edges of the Problem.
func (Problem) Edges() []ent.Edge {
	return nil
}
