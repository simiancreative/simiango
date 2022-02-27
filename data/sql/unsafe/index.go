package unsafe

import (
	"github.com/simiancreative/simiango/data/sql"
)

type Unsafe struct {
	Cx    sql.ConnX
	Query string
}

type UnsafeItem map[string]interface{}
type UnsafeContent []UnsafeItem

func (u *Unsafe) UnsafeSelect(params ...interface{}) (UnsafeContent, error) {
	items := UnsafeContent{}

	if err := u.Cx.Select(&items, u.Query, params...); err != nil {
		return nil, err
	}

	return items, nil
}
