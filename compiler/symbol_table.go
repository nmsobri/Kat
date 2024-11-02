package compiler

type SymbolScope string

type Symbol struct {
	Name  string
	Index int
	Scope SymbolScope
}

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type SymbolTable struct {
	Symbols             map[string]Symbol
	NumberOfIdentifiers int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Symbols:             make(map[string]Symbol),
		NumberOfIdentifiers: 0,
	}
}

func (st *SymbolTable) Define(name string) Symbol {
	index := st.NumberOfIdentifiers
	st.NumberOfIdentifiers++

	symbol := Symbol{
		Name:  name,
		Index: index,
		Scope: GlobalScope,
	}

	st.Symbols[name] = symbol
	return symbol
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.Symbols[name]
	return sym, ok
}
