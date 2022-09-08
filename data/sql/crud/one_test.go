package crud_test

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/simiancreative/simiango/data/sql/crud"
	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	m := crud.Model{
		Table: "products",
	}

	item := &Item{}

	clos := m.SetupOneTest(t)
	defer clos()

	err := m.One(item, 456)
	assert.NoError(t, err)
}

func TestAugmentOne(t *testing.T) {
	augmentOneFunc := func(ds *goqu.SelectDataset) *goqu.SelectDataset {
		return ds
	}

	m := crud.Model{
		Table:           "products",
		AugmentOneQuery: &augmentOneFunc,
	}

	item := &Item{}

	clos := m.SetupOneTest(t)
	defer clos()

	err := m.One(item, 456)
	assert.NoError(t, err)
}
