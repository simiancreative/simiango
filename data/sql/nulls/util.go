package nulls

import (
	// "fmt"
	"strconv"
	// "encoding/binary"
	// "math"
)

func Float64fromUint8s(i []byte) float64 {
	const bitSize = 64 // Don't think about it to much. It's just 64 bits.
	float, _ := strconv.ParseFloat(string(i), bitSize)

	// bits := binary.LittleEndian.Uint64(i)
	// float := math.Float64frombits(bits)
	return float
}
