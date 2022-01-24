package nulls

func init() {
	registerTest(test{
		Name: "Numeric - Float",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Numeric{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},

		Param:     0.00,
		Matcher:   0.00,
		FailParam: "hi",
	})

	registerTest(test{
		Name: "Numeric - []unit8",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Numeric{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},

		Param:     make([]uint8, 8),
		Matcher:   0.00,
		FailParam: "hi",
	})

	registerTest(test{
		Name: "Numeric - []byte 6",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &Numeric{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},

		Param:     make([]byte, 6),
		Matcher:   float64(0),
		FailParam: "hi",
	})
}
