package nulls

import (
	"time"
)

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*Time)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*Time)
		return inst.Time
	}

	t := time.Time{}

	registerTest(test{
		Name:            "Time",
		Valid:           valid,
		Value:           value,
		Param:           t,
		Matcher:         t,
		FailParam:       "hi",
		MarshalledParam: t,

		GetInst: func() interface{} {
			return &Time{}
		},
	})
}
