package crud_test

import (
	"testing"

	"github.com/simiancreative/simiango/data/sql/crud"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	m := crud.Model{
		Table: "products",
	}

	params := Item{ID: "456"}
	item := &Item{}

	clos := m.SetupUpdateTest(t, params, item)
	defer clos()

	err := m.Update(456, params, item)
	assert.NoError(t, err)
}
