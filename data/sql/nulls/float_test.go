package nulls

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*Float64)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*Float64)
		return inst.Value
	}
	registerTest(test{
		Name: "Float64",
		GetInst: func() interface{} {
			return &Float64{}
		},
		Valid:   valid,
		Value:   value,
		Param:   0.00,
		Matcher: 0.00,

		FailParam: "hi",
		MarshalledParam: 0.00,
	})
}
