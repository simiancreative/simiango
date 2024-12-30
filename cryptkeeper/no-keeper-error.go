package cryptkeeper

import (
	_ "embed"

	"fmt"
	"strings"
)

//go:embed no-keeper-error.txt
var notFoundTxt string

func NotFoundError(kind int) error {
	return NotFound{kind: kind}
}

type NotFound struct {
	kind int
}

func (e NotFound) Error() string {
	return fmt.Sprintf(notFoundTxt, e.kind, availableKeepers())
}

func availableKeepers() string {
	keys := []string{""}

	for k := range registered {
		keys = append(keys, fmt.Sprintf("✅ %v", names[k]))
	}

	return strings.Join(keys, "\n  ")
}
