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

func (p ParamItem) SetValue(value string) {
	p.Value = value
	p.Values = []string{value}
}

func (p *ParamItem) SetValues(values []string) {
	if len(values) > 0 {
		p.Value = values[0]
	}
	p.Values = values
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

func setVal(fv reflect.Value, val interface{}) {
	switch fv.Kind() {

	case reflect.Ptr:
		handlePtr(fv, val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value, err := strconv.ParseInt(val.(string), 10, 64); err == nil {
			fv.SetInt(value)
		}

	case reflect.String:
		fv.SetString(val.(string))

	case reflect.Bool:
		fv.SetBool(val == "true")

	case reflect.Float64:
		if value, err := strconv.ParseFloat(val.(string), 64); err == nil {
			fv.SetFloat(value)
		}
	}
}

func handlePtr(fv reflect.Value, val interface{}) {
	if val == nil {
		src := reflect.Zero(fv.Type())
		fv.Set(src)
		return
	}

	fv.Set(reflect.New(fv.Type().Elem()))
	deref := fv.Elem()

	setVal(deref, val)
}
