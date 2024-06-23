package environment

import "kat/value"

type Environment struct {
	Envs map[string]value.Value
}

func New() *Environment {
	return &Environment{
		Envs: make(map[string]value.Value),
	}
}
