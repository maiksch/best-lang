package evaluator

import "github.com/maiksch/best-lang/parser"

type Environment struct {
	identifiers map[string]Object
}

func NewEnvrionment() *Environment {
	return &Environment{
		identifiers: make(map[string]Object),
	}
}

func (e *Environment) get(identifier *parser.Identifier) Object {
	val, ok := e.identifiers[identifier.Value]
	if !ok {
		return newError("unknown identifier %s", identifier.Value)
	}
	return val
}

func (e *Environment) set(identifier *parser.Identifier, value Object) Object {
	e.identifiers[identifier.Value] = value
	return value
}
