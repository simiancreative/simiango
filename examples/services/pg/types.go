package pgservice

import (
	"github.com/simiancreative/simiango/meta"
)

type Product struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Products []Product

type sampleService struct {
	id meta.RequestId
}
