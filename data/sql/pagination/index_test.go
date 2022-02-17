package pagination

import (
	"testing"

	m "github.com/simiancreative/simiango/mocks/data/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ConnXMock *m.ConnX
var total int

type sampleItem interface{}

func init() {
	ConnXMock = &m.ConnX{}

	ConnXMock.
		On(
			"Get",
			&total,
			`
	SELECT COUNT(*)
	FROM devices
	WHERE gateway_id = ?
	`,
			"1234",
		).
		Return(
			func(v interface{}, s string, s2 ...interface{}) error {
				r, _ := v.(*int)
				*r = 10000
				return nil
			},
		)

	ConnXMock.
		On(
			"Select",
			mock.Anything,
			`
	SELECT id
	FROM devices
	WHERE gateway_id = ?
	ORDER BY created_at DESC
	LIMIT 25
	OFFSET 0
	`,
			"1234",
		).
		Return(
			func(v interface{}, s string, s2 ...interface{}) error {
				r, _ := v.(*Items)
				*r = Items{{ID: "456"}, {ID: "789"}}
				return nil
			},
		)
}

type Item struct {
	ID string `db:"id"`
}

type Items []Item

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
