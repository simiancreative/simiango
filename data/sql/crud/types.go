package crud

import "github.com/doug-martin/goqu/v9"

type FilterableFunc func(interface{}) (Kind, interface{})
type Filterable map[string]FilterableFunc
type AugmentListQuery *func(*goqu.SelectDataset) *goqu.SelectDataset
type Exp map[string]interface{}
type Filters map[string]interface{}
type Order map[string]interface{}
type Size int
type Page int

type Joinable struct {
	Table string
	Exp
}

type Kind int

const (
	AND Kind = iota
	OR
	JOIN
)
