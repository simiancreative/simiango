package unsafe

import (
	"github.com/simiancreative/simiango/data/sql/nulls"
	"strconv"
)

func coerceValue(kind string, value []byte) (interface{}, error) {
	switch kind {
	case "DECIMAL":
		dec := nulls.Decimal{}
		dec.Scan(value)
		return dec.Value, nil
	case "BIGINT":
		val, _ := strconv.Atoi(string(value))
		return val, nil
	case "VARCHAR":
		return string(value), nil
	default:
		return string(value), nil
	}
}
