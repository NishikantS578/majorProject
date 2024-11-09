package vm

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store           map[string]Symbol
	num_definitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol), num_definitions: 0}
}

func (s *SymbolTable) Define(name string) Symbol {
	s.store[name] = Symbol{Name: name, Index: s.num_definitions, Scope: GlobalScope}
	s.num_definitions++
	return s.store[name]
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return obj, ok
}
