package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	p := Page{
		Cx:         ConnXMock,
		Attributes: "id",
		From:       "devices",
		Where:      "gateway_id = ?",
		Order:      "created_at DESC",
		Page:       1,
		PageSize:   25,
	}

	v := Items{}

	content, err := p.Select(&v, "1234")
	contentItems, _ := content.Content.(*Items)

	assert.NoError(t, err)
	assert.Equal(t, (*contentItems)[0].ID, "456")
	assert.Equal(t, v[0].ID, "456")
}
