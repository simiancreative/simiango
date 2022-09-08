package crud

import "github.com/doug-martin/goqu/v9"

type AugmentableFunc func(*goqu.SelectDataset, interface{}) *goqu.SelectDataset
type AugmentableQuery *func(*goqu.SelectDataset) *goqu.SelectDataset
type Augmentations map[string]AugmentableFunc
type Filters map[string]interface{}
type OrderColumn struct {
	Name string
	Direction
}
type Order []OrderColumn
type Columns []interface{}
type Size int
type Page int

type Direction int

const (
	ASC Direction = iota
	DSC
)
