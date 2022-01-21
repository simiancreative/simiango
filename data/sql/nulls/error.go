package nulls

import "fmt"

func handleError(value interface{}) error {
	return fmt.Errorf("value not valid: %T, %v", value, value)
}
