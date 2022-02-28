package unsafe

import (
	"testing"

	m "github.com/simiancreative/simiango/mocks/data/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ConnXMock *m.ConnX

var mockItem = Item{"id": 23, "name": "Floppy diskette"}

var mockContent = Content{
	// demonstrating that this is not a type safe solution and has the potential
	// to cause runtime errors
	Item{"id": 23, "name": "Floppy diskette"},
	Item{"id": "one", "name": 42},
}

func init() {
	ConnXMock = &m.ConnX{}

	ConnXMock.
		On(
			"Select",
			mock.Anything,
			`SELECT id FROM devices WHERE gateway_id = ?`,
			"1234",
		).
		Return(
			func(v interface{}, s string, s2 ...interface{}) error {
				r, _ := v.(*Content)
				*r = mockContent

				return nil
			},
		)

	ConnXMock.
		On(
			"Get",
			mock.Anything,
			`SELECT id FROM devices WHERE gateway_id = ?`,
			"1234",
		).
		Return(
			func(v interface{}, s string, s2 ...interface{}) error {
				r, _ := v.(*Item)
				*r = mockItem

				return nil
			},
		)
}

func TestUnsafeSelect(t *testing.T) {
	u := Unsafe{Cx: ConnXMock}

	content, err := u.UnsafeSelect(
		"SELECT id FROM devices WHERE gateway_id = ?",
		"1234",
	)

	assert.NoError(t, err)
	assert.Equal(t, content, mockContent)
}

func TestUnsafeGet(t *testing.T) {
	u := Unsafe{Cx: ConnXMock}

	content, err := u.UnsafeGet(
		"SELECT id FROM devices WHERE gateway_id = ?",
		"1234",
	)

	assert.NoError(t, err)
	assert.Equal(t, content, mockItem)
}
