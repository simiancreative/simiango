package crud

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/simiancreative/simiango/data/sql"
)

type Model struct {
	MultiTenant      bool
	Paginate         bool
	Table            string
	Columns          []interface{}
	Expressions    Expressions
	AugmentListQuery AugmentableQuery
	AugmentOneQuery  AugmentableQuery
	cx               sql.ConnX
	dialect          goqu.DialectWrapper
}

func (m *Model) Initialize(dialect string, cx sql.ConnX) {
	m.cx = cx
	m.dialect = goqu.Dialect(dialect)
}
