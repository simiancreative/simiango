package combinators

import "github.com/simiancreative/simiango/data/sql"

type Model struct {
	Cx sql.ConnX
}

func New(cx sql.ConnX) Model {
	return Model{Cx: cx}
}
