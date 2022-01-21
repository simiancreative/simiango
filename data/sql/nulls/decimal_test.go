package nulls

func init() {
	registerTest(test{
		Name: "Decimal",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Decimal{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},

		Param:     0.00,
		Matcher:   0.00,
		FailParam: "hi",
	})

	registerTest(test{
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Decimal{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},

		Param:     make([]uint8, 8),
		Matcher:   0.00,
		FailParam: "hi",
	})
}
