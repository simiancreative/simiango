package mssqlexample

import (
	"github.com/simiancreative/simiango/meta"
)

type Service struct {
	ID meta.RequestId
}

type Product struct {
	ID   string `json:"id" db:"ID"`
	Name string `json:"name" db:"ProductName"`
}
