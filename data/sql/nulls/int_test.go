package nulls

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*Int64)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*Int64)
		return inst.Value
	}

	registerTest(test{
		Name: "Int64",
		GetInst: func() interface{} {
			return &Int64{}
		},
		Valid:   valid,
		Value:   value,
		Param:   int64(1),
		Matcher: int64(1),

		FailParam:       "hi",
		MarshalledParam: 1,
	})
}
