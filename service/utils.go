package service

import (
	// "errors"
	"reflect"
	"strconv"
)

type ParamItem struct {
	Key    string
	Value  string
	Values []string
}

type ParamHolder interface {
	Get(string) (ParamItem, bool)
}

func parseParam(tagName string, v interface{}, params ParamHolder) error {
	t := reflect.TypeOf(v).Elem()
	rv := reflect.ValueOf(v).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup(tagName)

		if !ok {
			continue
		}

		val, ok := params.Get(tag)

		if !ok {
			continue
		}

		fv := rv.FieldByName(field.Name)

		if ok := fv.CanSet(); !ok {
			continue
		}

		setVal(fv, val.Value)
	}

	return nil
}

func setVal(fv reflect.Value, val string) {
	switch fv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value, err := strconv.ParseInt(val, 10, 64); err == nil {
			fv.SetInt(value)
		}

	case reflect.String:
		fv.SetString(val)

	case reflect.Bool:
		fv.SetBool(val == "true")

	case reflect.Float64:
		if value, err := strconv.ParseFloat(val, 64); err == nil {
			fv.SetFloat(value)
		}
	}
}
