package crud

import "github.com/doug-martin/goqu/v9"

type AugmentableFunc func(*goqu.SelectDataset, interface{}) *goqu.SelectDataset
type AugmentableQuery *func(*goqu.SelectDataset) *goqu.SelectDataset
type Augmentations map[string]AugmentableFunc
type Filters map[string]interface{}
type Order map[string]interface{}
type Columns []interface{}
type Size int
type Page int
