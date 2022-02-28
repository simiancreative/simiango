package nulls

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*Bool)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*Bool)
		return inst.Value
	}

	registerTest(test{
		Name: "Bool",
		GetInst: func() interface{} {
			return &Bool{}
		},
		Valid:   valid,
		Value:   value,
		Param:   true,
		Matcher: true,

		FailParam: "hi",
		MarshalledParam: true,
	})
}
