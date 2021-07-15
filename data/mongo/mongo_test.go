package mongo

/*
import (
	"fmt"
	"os"
	"strings"

	"bou.ke/monkey"
)

var ConnXMock *m.ConnX

func init() {
	ConnXMock = &m.ConnX{}
	Cx = ConnXMock
}

func TestConnect(t *testing.T) {
    mockError := errors.New("uh oh")
    subtests := []struct {
        name        string
        u, p, a, db string
        ConnectToMongo   func(string, string) (*sql.DB, error)
        expectedErr error
    }{
        {
            name: "Test Connect",
            ConnectToMongo: func(s string, s2 string) (a *mongo.Database, b *mongo.Client) {
                if s != "u:p@a/db" {
                    return nil, errors.New("wrong connection string")
                }
                return nil, nil
            },
        },
        {
            name: "error from Connect",
            ConnectToMongo: func(s string, s2 string) (a *mongo.Database, b *mongo.Client) {
                return nil, nil
            },
            expectedErr: mockError,
        },
    }
    for _, subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
			Connect = subtest.Connect
            _, err := init()
            if !errors.Is(err, subtest.expectedErr) {
                t.Errorf("expected error (%v), got error (%v)", subtest.expectedErr, err)
            }
        })
    }
}
*/
