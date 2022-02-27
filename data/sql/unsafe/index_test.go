package unsafe

import (
	"testing"

	m "github.com/simiancreative/simiango/mocks/data/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ConnXMock *m.ConnX

type sampleItem interface{}

var mockContent = UnsafeContent{
	// demonstrating that this is not a type safe solution and has the potential
	// to cause runtime errors
	UnsafeItem{"id": 23, "name": "Floppy diskette"},
	UnsafeItem{"id": "one", "name": 42},
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
				r, _ := v.(*UnsafeContent)
				*r = mockContent

				return nil
			},
		)
}

func TestUnsafeQuery(t *testing.T) {
	u := Unsafe{
		Cx:    ConnXMock,
		Query: "SELECT id FROM devices WHERE gateway_id = ?",
	}

	content, err := u.UnsafeSelect("1234")

	assert.NoError(t, err)
	assert.Equal(t, content, mockContent)
}
