package nulls

func init() {
	registerTest(test{
		Name: "Bool",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Bool{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},
		Param:   true,
		Matcher: true,

		FailParam: "hi",
	})
}
