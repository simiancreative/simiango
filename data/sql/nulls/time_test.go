package nulls

import "time"

func init() {
	t := time.Time{}

	registerTest(test{
		Name: "Time",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Time{}
			err := inst.Scan(v)

			return inst.Valid, inst.Time, err
		},
		Param:   t,
		Matcher: t,

		FailParam: "hi",
	})
}
