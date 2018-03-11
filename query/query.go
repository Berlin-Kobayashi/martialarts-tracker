package query

const KindAnd = kind("AND")
const KindOr = kind("OR")

type Query struct {
	Q map[string]FieldQuery
}

type FieldQuery struct {
	Kind   kind
	Values []interface{}
}

type kind string
