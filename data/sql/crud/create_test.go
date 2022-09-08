package crud_test

import (
	"testing"

	"github.com/simiancreative/simiango/data/sql/crud"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	m := crud.Model{
		Table: "products",
	}

	params := Item{ID: "456"}
	item := &Item{}

	clos := m.SetupCreateTest(t, params)
	defer clos()

	err := m.Create(params, item)
	assert.NoError(t, err)
}
