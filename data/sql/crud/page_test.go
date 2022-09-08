package crud_test

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/simiancreative/simiango/data/sql/crud"
	"github.com/stretchr/testify/assert"
)

type sampleItem interface{}

type Item struct {
	ID string `db:"id"`
}

type Items []Item

func setupPage(
	m *crud.Model,
	filters crud.Filters,
	order crud.Order,
	page int,
	size int,
) *crud.Model {

	if m == nil {
		m = &crud.Model{Table: "products"}
	}

	m.SetupPageTest(
		&Items{},
		filters,
		order,
		page,
		size,
		func(v interface{}, s string, s2 ...interface{}) error {
			r, _ := v.(*Items)
			*r = Items{{ID: "456"}, {ID: "789"}}
			return nil
		},
	)

	return m
}

func TestPage(t *testing.T) {
	filters := crud.Filters{}
	order := crud.Order{}
	page := 1
	size := 25

	m := setupPage(nil, filters, order, page, size)

	items := Items{}
	result, err := m.Page(&items, filters, order, page, size)
	contentItems, _ := result.Content.(*Items)

	assert.NoError(t, err)
	assert.Equal(t, (*contentItems)[0].ID, "456")
	assert.Equal(t, items[0].ID, "456")
}

func TestColumns(t *testing.T) {
	filters := crud.Filters{}
	order := crud.Order{}
	page := 1
	size := 25

	m := &crud.Model{Table: "products", Columns: crud.Columns{
		goqu.T("products").Col("id"),
	}}
	m = setupPage(m, filters, order, page, size)

	items := Items{}
	_, err := m.Page(&items, filters, order, page, size)

	assert.NoError(t, err)
}

func TestOrder(t *testing.T) {
	filters := crud.Filters{}
	order := crud.Order{
		{Name: "name", Direction: crud.ASC},
		{Name: "id", Direction: crud.DSC},
	}
	page := 1
	size := 25

	m := setupPage(nil, filters, order, page, size)

	items := Items{}
	_, err := m.Page(&items, filters, order, page, size)

	assert.NoError(t, err)
}

func TestFilters(t *testing.T) {
	filters := crud.Filters{"name": "hi", "id": "nope"}
	order := crud.Order{}
	page := 1
	size := 25

	m := &crud.Model{Table: "products", Expressions: crud.Expressions{
		"name": func(ds *goqu.SelectDataset, value interface{}) *goqu.SelectDataset {
			return ds.Where(goqu.Ex{"products.id": value})
		},
	}}
	m = setupPage(m, filters, order, page, size)

	items := Items{}
	_, err := m.Page(&items, filters, order, page, size)

	assert.NoError(t, err)
}

func TestAugmentListQuery(t *testing.T) {
	filters := crud.Filters{}
	order := crud.Order{}
	page := 1
	size := 25

	augmentListQuery := func(ds *goqu.SelectDataset) *goqu.SelectDataset {
		return ds
	}

	m := &crud.Model{
		Table:            "products",
		AugmentListQuery: &augmentListQuery,
	}
	m = setupPage(m, filters, order, page, size)

	items := Items{}
	_, err := m.Page(&items, filters, order, page, size)

	assert.NoError(t, err)
}
