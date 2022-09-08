package crud_test

import (
	"testing"

	"github.com/simiancreative/simiango/data/sql/crud"
	"github.com/simiancreative/simiango/service"
	"github.com/stretchr/testify/assert"
)

func TestPageFromReq(t *testing.T) {
	filters := crud.Filters{"name": "Car"}
	order := crud.Order{
		{Name: "name", Direction: crud.DSC},
		{Name: "id", Direction: crud.ASC},
		{Name: "turtle", Direction: crud.ASC},
	}
	page := 1
	size := 25

	m := setupPage(nil, filters, order, page, size)

	req := service.Req{}
	req.Params = service.RawParams{
		{Key: "name", Values: []string{"Car"}},
		{Key: "order", Values: []string{"name,dsc", "id", "turtle,asc"}},
	}
	items := &Items{}
	_, err := m.PageFromReq(items, req)

	assert.NoError(t, err)
}

func TestPageFromReqNoOrder(t *testing.T) {
	filters := crud.Filters{"name": "Car"}
	order := crud.Order{}
	page := 1
	size := 25

	m := setupPage(nil, filters, order, page, size)

	req := service.Req{}
	req.Params = service.RawParams{
		{Key: "name", Values: []string{"Car"}},
	}
	items := &Items{}
	_, err := m.PageFromReq(items, req)

	assert.NoError(t, err)
}
