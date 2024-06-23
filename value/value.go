package value

import "fmt"

type Value interface {
	String() string
}

type ValueFloat struct {
	Value float64
}

func (vf ValueFloat) String() string {
	return fmt.Sprintf("%.2f", vf.Value)
}

type ValueInt struct {
	Value int64
}

func (vi ValueInt) String() string {
	return fmt.Sprintf("%d", vi.Value)
}

type ValueBool struct {
	Value bool
}

func (vb ValueBool) String() string {
	return fmt.Sprintf("%t", vb.Value)
}

type ValueString struct {
	Value string
}

func (vs ValueString) String() string {
	return vs.Value
}
