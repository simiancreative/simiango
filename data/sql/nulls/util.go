package nulls

import (
	"encoding/binary"
	"math"
)

func Float64fromUint8s(i []byte) float64 {
	bits := binary.LittleEndian.Uint64(i)
	float := math.Float64frombits(bits)
	return float
}

