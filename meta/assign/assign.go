package assign

import (
	"reflect"
)

type Assignable map[string]interface{}

func (ps Assignable) Assign(v interface{}) error {
	return parseAssignable("assign", v, ps)
}

func parseAssignable(tagName string, v interface{}, assignable Assignable) error {
	t := reflect.TypeOf(v).Elem()
	rv := reflect.ValueOf(v).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup(tagName)

		if !ok {
			continue
		}

		val, ok := assignable[tag]

		if !ok {
			continue
		}

		fv := rv.FieldByName(field.Name)

		if ok := fv.CanSet(); !ok {
			continue
		}

		setVal(fv, val, field.Name)
	}

	return nil
}

func setVal(fv reflect.Value, val interface{}, name string) {
	switch fv.Kind() {

	case reflect.Ptr:
		if val == nil {
			src := reflect.Zero(fv.Type())
			fv.Set(src)
			return
		}

		fv.Set(reflect.New(fv.Type().Elem()))
		deref := fv.Elem()

		deref.Set(reflect.ValueOf(val).Convert(deref.Type()))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value, ok := val.(int64); ok {
			fv.SetInt(value)
		}

		if value, ok := val.(int); ok {
			fv.SetInt(int64(value))
		}

		if value, ok := val.(float64); ok {
			fv.SetInt(int64(value))
		}

	case reflect.String:
		if value, ok := val.(string); ok {
			fv.SetString(value)
		}

	case reflect.Bool:
		if value, ok := val.(bool); ok {
			fv.SetBool(value)
		}

	case reflect.Float64:
		if value, ok := val.(float64); ok {
			fv.SetFloat(value)
		}
	}
}
