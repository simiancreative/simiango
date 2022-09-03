package crud

import "github.com/doug-martin/goqu/v9"

type FilterableFunc func(interface{}) (WhereKind, interface{})
type Filterable map[string]FilterableFunc
type AugmentListQuery *func(*goqu.SelectDataset) *goqu.SelectDataset
type Filters map[string]interface{}
type Orders map[string]interface{}
type Size int
type Page int

type WhereKind int

const (
	AND WhereKind = iota
	OR
)
