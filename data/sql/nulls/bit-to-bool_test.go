package nulls

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*BitBool)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*BitBool)
		return inst.Value
	}

	inst := func() interface{} {
		return &BitBool{}
	}

	registerTest(test{
		Name:    "BitBool",
		GetInst: inst,
		Valid:   valid,
		Value:   value,
		Param:   []byte{1},
		Matcher: true,

		FailParam:       "hi",
		MarshalledParam: true,
	})

	registerTest(test{
		Name:    "BitBool",
		GetInst: inst,
		Valid:   valid,
		Value:   value,
		Param:   true,
		Matcher: true,

		FailParam:       "hi",
		MarshalledParam: true,
	})
}
