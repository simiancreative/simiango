package unsafe

import (
	"github.com/jmoiron/sqlx"
)

// Item is used as a container for sql rows with unknown values
type Item map[string]interface{}

// Content is used as a container for Unsafe Items
type Content []Item

// Column is used to inform the column order as well as column type
type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Result is the container for all sql result info
type Result struct {
	Columns []Column      `json:"columns"`
	Content []interface{} `json:"content"`
}

func (r *Result) addColumns(rows *sqlx.Rows) error {
	if len(r.Content) > 0 {
		return nil
	}

	colTypes, err := rows.ColumnTypes()

	cols, err := rows.Columns()

	if err != nil {
		return err
	}

	for i, col := range cols {
		r.Columns = append(r.Columns, Column{
			Name: col,
			Type: string(colTypes[i].DatabaseTypeName()),
		})
	}

	return nil
}

func (r *Result) addItem(rows *sqlx.Rows) error {
	item := Item{}
	values, err := rows.SliceScan()

	if err != nil {
		return err
	}

	for i, value := range values {
		col := r.Columns[i]

		switch v := value.(type) {
		case []byte:
			newVal, err := coerceValue(col.Type, v)
			if err != nil {
				return err
			}

			item[col.Name] = newVal
		default:
			item[col.Name] = v
		}

	}

	r.Content = append(r.Content, item)

	return nil
}
