package assign

import "github.com/simiancreative/simiango/meta/assign"

type Event struct {
	Type string            `json:"type"`
	Body assign.Assignable `json:"body"`
}

// Types that can be null and must remain in null
// must be pointers
type AssignableStruct struct {
	ID         string
	StatusName string  `assign:"status_name"`
	Error      *string `assign:"error"`
	Status     *int    `assign:"status"`
}
