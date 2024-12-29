package errors

import (
	"github.com/simiancreative/simiango/service"
)

type Reason map[string]interface{}

type Item struct {
	Status      int
	Key         string
	Description string
}

type Reasons map[int]Item

func (r Reasons) Result(id int, reasons ...Reason) *service.ResultError {
	item := r[id]

	maps := []map[string]interface{}{}
	for _, reason := range reasons {
		maps = append(maps, reason)
	}

	return &service.ResultError{
		Status:     item.Status,
		ErrMessage: item.Description,
		Message:    item.Key,
		Reasons:    maps,
	}
}
