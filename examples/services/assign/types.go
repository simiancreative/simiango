package assign

import "github.com/simiancreative/simiango/meta/assign"

type Event struct {
	Type string            `json:"type"`
	Body assign.Assignable `json:"body"`
}

// Types that can be null and must remain in null
// must be pointers
type AssignableStruct struct {
	ID         string  `json:"-"`
	StatusName string  `json:"status_name,omitempty" assign:"status_name"`
	Error      *string `json:"error,omitempty" assign:"error"`
	Status     *int    `json:"status,omitempty" assign:"status"`
}
