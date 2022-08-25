package nulls

func init() {
	valid := func(v interface{}) bool {
		inst := v.(*String)
		return inst.Valid
	}

	value := func(v interface{}) interface{} {
		inst := v.(*String)
		return inst.Value
	}

	registerTest(test{
		Name: "String",
		GetInst: func() interface{} {
			return &String{}
		},
		Valid:   valid,
		Value:   value,
		Param:   "hi",
		Matcher: "hi",

		FailParam:       0,
		MarshalledParam: "hi",
	})

	registerTest(test{
		Name: "String",
		GetInst: func() interface{} {
			return &String{}
		},
		Valid:   valid,
		Value:   value,
		Param:   []byte("hi"),
		Matcher: "hi",

		FailParam:       0,
		MarshalledParam: "hi",
	})
}
