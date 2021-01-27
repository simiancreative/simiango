package pgservice

import (
	p "github.com/simiancreative/simiango/data/pg"
)

func (s sampleService) Result() (interface{}, error) {
	products := Products{}

	err := p.Cx.Select(&products, "SELECT * FROM products")

	return products, err
}
