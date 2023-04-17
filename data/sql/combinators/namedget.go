package combinators

import "github.com/simiancreative/simiango/data/sql"

type Model struct {
	Cx sql.ConnX
}

func New(cx sql.ConnX) Model {
	return Model{Cx: cx}
}

func (m Model) NamedGet(dst interface{}, query string, params interface{}) error {
	nstmt, err := m.Cx.PrepareNamed(query)
	if err != nil {
		return err
	}

	return nstmt.Get(dst, params)
}

func (m Model) NamedSelect(dst interface{}, query string, params interface{}) error {
	nstmt, err := m.Cx.PrepareNamed(query)
	if err != nil {
		return err
	}

	return nstmt.Select(dst, params)
}
