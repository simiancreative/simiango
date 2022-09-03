package crud

import (
	"github.com/doug-martin/goqu/v9"
)

func (m *Model) Update(id interface{}, params interface{}, destination interface{}) error {
	ds := m.dialect.Update(m.Table).Set(params).Where(goqu.Ex{"id": id})

	query, qParams, _ := ds.ToSQL()
	if _, err := m.cx.Exec(query, qParams...); err != nil {
		return err
	}

	if err := m.One(destination, id); err != nil {
		return err
	}

	return nil
}
