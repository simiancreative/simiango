package crud_test

import (
	"testing"

	"github.com/simiancreative/simiango/data/sql/crud"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	m := crud.Model{
		Table: "products",
	}

	clos := m.SetupDeleteTest(t)
	defer clos()

	err := m.Delete(456)
	assert.NoError(t, err)
}
