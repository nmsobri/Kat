package environment

import (
	"kat/value"
)

type Environment struct {
	Envs   map[string]value.Value
	Parent *Environment
}

func New() *Environment {
	return &Environment{
		Envs:   make(map[string]value.Value),
		Parent: nil,
	}
}

func NewWithParent(parent *Environment) *Environment {
	return &Environment{
		Envs:   make(map[string]value.Value),
		Parent: parent,
	}
}

func (env *Environment) Set(key string, value value.Value) {
	env.Envs[key] = value
}

func (env *Environment) Get(key string) (value.Value, bool) {
	val, ok := env.Envs[key]

	if !ok {
		if env.Parent != nil {
			return env.Parent.Get(key)
		}
	}

	return val, ok
}
