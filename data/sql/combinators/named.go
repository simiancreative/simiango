package combinators

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
