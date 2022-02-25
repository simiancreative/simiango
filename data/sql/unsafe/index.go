package unsafe

import (
	"github.com/simiancreative/simiango/data/sql"
)

func Query(Cx sql.ConnX, query string, params ...interface{}) (interface{}, error) {
	rows, err := Cx.Query(query, params...)

	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))
	dest := make([]interface{}, len(cols))

	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}
	}
	return result, nil
}
