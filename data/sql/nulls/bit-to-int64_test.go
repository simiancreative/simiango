package nulls

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*BitToInt64)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*BitToInt64)
		return inst.Value
	}

	registerTest(test{
		Name: "BitToInt64",
		GetInst: func() interface{} {
			return &BitToInt64{}
		},
		Valid:   valid,
		Value:   value,
		Param:   true,
		Matcher: int64(1),

		FailParam:       "hi",
		MarshalledParam: true,
	})
}
