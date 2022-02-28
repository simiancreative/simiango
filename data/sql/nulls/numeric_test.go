package nulls

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*Numeric)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*Numeric)
		return inst.Value
	}

	registerTest(test{
		Name: "Numeric - Float",
		GetInst: func() interface{} {
			return &Numeric{}
		},
		Valid:           valid,
		Value:           value,
		Param:           0.00,
		Matcher:         0.00,
		FailParam:       "hi",
		MarshalledParam: 0.00,
	})

	registerTest(test{
		Name: "Numeric - []unit8",
		GetInst: func() interface{} {
			return &Numeric{}
		},
		Valid:           valid,
		Value:           value,
		Param:           make([]uint8, 8),
		Matcher:         0.00,
		FailParam:       "hi",
		MarshalledParam: 0.00,
	})

	registerTest(test{
		Name: "Numeric - []byte 6",
		GetInst: func() interface{} {
			return &Numeric{}
		},
		Valid:           valid,
		Value:           value,
		Param:           make([]byte, 6),
		Matcher:         float64(0),
		FailParam:       "hi",
		MarshalledParam: 0.00,
	})
}
