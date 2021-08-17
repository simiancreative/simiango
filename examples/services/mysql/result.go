package mysqlexample

import (
	m "github.com/simiancreative/simiango/data/mysql"
)

func (s Service) Result() (interface{}, error) {
	rows := []Product{}

	err := m.Cx.Select(&rows, "select * from products")

	return rows, err
}
