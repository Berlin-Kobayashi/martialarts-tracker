package query

const (
	KindAnd      kind = iota
	KindOr
	KindContains
)

type Query struct {
	Q map[string]FieldQuery
}

type FieldQuery struct {
	Kind   kind
	Values []interface{}
}

type kind uint
