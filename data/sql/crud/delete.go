package crud

import "github.com/doug-martin/goqu/v9"

func (m *Model) Delete(id interface{}) error {
	ds := m.dialect.Delete(m.Table).Where(goqu.Ex{"id": id})

	query, qParams, _ := ds.ToSQL()
	_, err := m.cx.Exec(query, qParams...)
	if err != nil {
		return err
	}

	return nil
}
