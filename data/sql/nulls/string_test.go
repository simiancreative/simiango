package nulls

func init() {
	registerTest(test{
		Name: "String",
		Func: func(v interface{}) (bool, interface{}, error) {
			inst := &String{}
			err := inst.Scan(v)

			return inst.Valid, inst.Value, err
		},
		Param:   "hi",
		Matcher: "hi",

		FailParam: 0,
	})
}
