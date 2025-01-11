package temporal

import (
	"fmt"

	"github.com/simiancreative/simiango/errors"
)

func createModelNotFoundError(name string) error {
	list := ""
	for name := range registered {
		list += fmt.Sprintf("✅ %v\n", name)
	}

	return errors.New(
		"Model not found (❌ %v) available models are:\n\n%v",
		name, list,
	)
}
