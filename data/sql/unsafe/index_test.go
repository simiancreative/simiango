package unsafe

import (
	"testing"

	m "github.com/simiancreative/simiango/mocks/data/mysql"
	"github.com/stretchr/testify/assert"
)

var ConnXMock *m.ConnX

type sampleItem interface{}

func init() {
	ConnXMock = &m.ConnX{}

	ConnXMock.
		On(
			"Query",
			`SELECT id FROM devices WHERE gateway_id = ?`,
			"1234",
		).
		Return(
			func(v interface{}, s string, s2 ...interface{}) error {
				r, _ := v.(*Items)
				*r = Items{{ID: "456"}}
				return nil
			},
		)
}

type Item struct {
	ID string `db:"id"`
}

type Items []Item

type Request struct {
	Input1 int64
}

func TestUnsafeQuery(t *testing.T) {
	v := Items{}

	content, err := Query(ConnXMock, "SELECT id FROM devices WHERE gateway_id = ?", "1234")
	contentItems, _ := content.(*Items)

	assert.NoError(t, err)
	assert.Equal(t, (*contentItems)[0].ID, "456")
	assert.Equal(t, v[0].ID, "456")

}
