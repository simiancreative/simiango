package nulls

func init() {
	registerTest(test{
		Name: "Int64",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Int64{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},
		Param:   int64(1),
		Matcher: int64(1),

		FailParam: "hi",
	})
}
