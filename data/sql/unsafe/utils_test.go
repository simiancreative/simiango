package unsafe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafeCoerceBIGINT(t *testing.T) {
	b := []byte(fmt.Sprintf("%v", 42))
	res, err := coerceValue("BIGINT", b)

	assert.NoError(t, err)
	assert.Equal(t, 42, res)
}

func TestUnsafeCoerceDECIMAL(t *testing.T) {
	b := []byte(fmt.Sprintf("%v", 42.54))
	res, err := coerceValue("DECIMAL", b)

	assert.NoError(t, err)
	assert.Equal(t, 42.54, res)
}

func TestUnsafeCoerceVARCHAR(t *testing.T) {
	b := []byte("HI THERE")
	res, err := coerceValue("VARCHAR", b)

	assert.NoError(t, err)
	assert.Equal(t, "HI THERE", res)
}

func TestUnsafeCoerceUNKNOWN(t *testing.T) {
	b := []byte("HO THERE")
	res, err := coerceValue("UNKNOWN", b)

	assert.NoError(t, err)
	assert.Equal(t, "HO THERE", res)
}
