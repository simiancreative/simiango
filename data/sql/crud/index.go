package crud

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/simiancreative/simiango/data/sql"
)

type Model struct {
	MultiTenant      bool
	Paginate         bool
	Table            string
	Filterable       Filterable
	AugmentListQuery AugmentListQuery
	AugmentOneQuery  AugmentListQuery
	cx               sql.ConnX
	dialect          goqu.DialectWrapper
}

func (m *Model) Initialize(cx sql.ConnX, dialect string) {
	m.cx = cx
	m.dialect = goqu.Dialect(dialect)
}
