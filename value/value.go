package value

import "fmt"

type Value interface {
	String() string
}

type ValueDouble struct {
	Value float64
}

func (vf ValueDouble) String() string {
	return fmt.Sprintf("%.2f", vf.Value)
}

type ValueInt struct {
	Value int64
}

func (vi ValueInt) String() string {
	return fmt.Sprintf("%d", vi.Value)
}
