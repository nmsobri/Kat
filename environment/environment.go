package environment

import (
	"fmt"
	"kat/value"
	"log"
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

func (env *Environment) Get(key string) (value.Value, bool) {
	val, ok := env.Envs[key]

	if !ok && env.Parent != nil {
		return env.Parent.Get(key)
	}

	return val, ok
}

func (env *Environment) Set(key string, value value.Value) {
	env.Envs[key] = value
}

func (env *Environment) Assign(key string, value value.Value) {
	ok := env.setWithParent(key, value)

	if !ok {
		msg := fmt.Sprintf("Variable %s is not found", key)
		log.Fatal(msg)
	}
}

func (env *Environment) setWithParent(key string, value value.Value) bool {
	if _, ok := env.Envs[key]; ok {
		env.Envs[key] = value
		return true
	}

	if env.Parent != nil {
		return env.Parent.setWithParent(key, value)
	}

	return false
}

func (env *Environment) String() string {
	return fmt.Sprintf("%v", env.Envs)
}

func (env *Environment) Type() value.Type {
	return value.TYPE_ENVIRONMENT
}
