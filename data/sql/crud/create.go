package crud

func (m *Model) Create(params interface{}, destination interface{}) error {
	ds := m.dialect.Insert(m.Table).Rows(params)

	query, qParams, _ := ds.ToSQL()
	rows, err := m.cx.Exec(query, qParams...)
	if err != nil {
		return err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return err
	}

	if err := m.One(destination, id); err != nil {
		return err
	}

	return nil
}
