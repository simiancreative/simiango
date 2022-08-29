package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	p := Page{
		Attributes: "id",
		From:       "devices",
		Where:      "gateway_id = ?",
		Order:      "created_at DESC",
		Page:       1,
		PageSize:   25,
	}

	p.SetupTest(
		&Items{},
		func(v interface{}, s string, s2 ...interface{}) error {
			r, _ := v.(*Items)
			*r = Items{{ID: "456"}, {ID: "789"}}
			return nil
		},
		"1234",
	)

	v := Items{}

	content, err := p.Select(&v, "1234")
	contentItems, _ := content.Content.(*Items)

	assert.NoError(t, err)
	assert.Equal(t, (*contentItems)[0].ID, "456")
	assert.Equal(t, v[0].ID, "456")
}
