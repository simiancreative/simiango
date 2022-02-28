package nulls

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*Decimal)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*Decimal)
		return inst.Value
	}

	registerTest(test{
		Name: "Decimal",
		GetInst: func() interface{} {
			return &Decimal{}
		},

		Valid:     valid,
		Value:     value,
		Param:     0.00,
		Matcher:   0.00,
		FailParam: "hi",
	})

	registerTest(test{
		Name: "Decimal - unint8",
		GetInst: func() interface{} {
			return &Decimal{}
		},
		Valid:     valid,
		Value:     value,
		Param:     make([]uint8, 8),
		Matcher:   0.00,
		FailParam: "hi",
		MarshalledParam: 0.00,
	})
}
