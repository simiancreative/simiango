# Validate

Basic wrapper around https://github.com/go-validator/validator#usage with some added features to support service.ResultError. 

## Builtin validators

```
len
	For numeric numbers, len will simply make sure that the
	value is equal to the parameter given. For strings, it
	checks that the string length is exactly that number of
	characters. For slices,	arrays, and maps, validates the
	number of items. (Usage: len=10)

max
	For numeric numbers, max will simply make sure that the
	value is lesser or equal to the parameter given. For strings,
	it checks that the string length is at most that number of
	characters. For slices,	arrays, and maps, validates the
	number of items. (Usage: max=10)

min
	For numeric numbers, min will simply make sure that the value
	is greater or equal to the parameter given. For strings, it
	checks that the string length is at least that number of
	characters. For slices, arrays, and maps, validates the
	number of items. (Usage: min=10)

nonzero
	This validates that the value is not zero. The appropriate
	zero value is given by the Go spec (e.g. for int it's 0, for
	string it's "", for pointers is nil, etc.) For structs, it
	will not check to see if the struct itself has all zero
	values, instead use a pointer or put nonzero on the struct's
	keys that you care about. For pointers, the pointer's value
	is used to test for nonzero in addition to the pointer itself
	not being nil. To just check for not being nil, use `nonnil`.
	(Usage: nonzero)

regexp
	Only valid for string types, it will validate that the
	value matches the regular expression provided as parameter.
	Commas need to be escaped with 2 backslashes `\\`.
	(Usage: regexp=^a.*b$)

nonnil
	Validates that the given value is not nil. (Usage: nonnil)
```
## Custom validators

Define a custom validator by using AddValidation.

```
// Very simple validator
func notZZ(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("notZZ only validates strings")
	}
	if st.String() == "ZZ" {
		return errors.New("value cannot be ZZ")
	}
	return nil
}
```

Then call AddValidation

```
validate.AddValidation("notzz", notZZ)
```

Now you can use notzz as a validation tag. This will print "Field A error: value cannot be ZZ"

```
type T struct {
	A string  `validate:"nonzero,notzz"`
}

t := T{"ZZ"}

if errs := validator.Validate(t); errs != nil {
	fmt.Printf("Field A error: %s\n", errs["A"][0])
}
```
