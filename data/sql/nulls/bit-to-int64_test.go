package nulls

func init() {
	registerTest(test{
		Name: "BitToInt64",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &BitToInt64{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},
		Param:   true,
		Matcher: int64(1),

		FailParam: "hi",
	})
}
