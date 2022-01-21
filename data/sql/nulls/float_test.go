package nulls

func init() {
	registerTest(test{
		Name: "Float64",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Float64{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},
		Param:   0.00,
		Matcher: 0.00,

		FailParam: "hi",
	})
}
