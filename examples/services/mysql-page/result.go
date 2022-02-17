package mysqlpage

import (
	m "github.com/simiancreative/simiango/data/mysql"
	"github.com/simiancreative/simiango/data/sql/pagination"
)

func (s Service) Result() (interface{}, error) {
	rows := []Product{}

	p := pagination.Page{
		Cx:         m.Cx,
		Attributes: "id, name",
		From:       "products",
		Where:      `name = ? or name = ?`,
		Order:      "name DESC",
		Page:       s.Page,
		PageSize:   s.Size,
	}

	return p.Select(&rows, s.Pattern, s.Pattern2)
}
