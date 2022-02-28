package service

import (
	"fmt"
	"reflect"
)

type ContentResponse struct {
	Content interface{} `json:"content"`
	ContentResponseContext
	ContentResponseMeta
}

type ContentResponseContext struct {
	First      bool `json:"first"`
	Last       bool `json:"last"`
	TotalPages int  `json:"total_pages"`
}

type ContentResponseMeta struct {
	Size  int   `json:"size"`
	Total int   `json:"total"`
	Page  int   `json:"current_page"`
	Error error `json:"error,omitempty"`
}

func ToContentResponse(
	resources interface{},
	meta ContentResponseMeta,
) ContentResponse {
	rv := reflect.ValueOf(resources)
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	if rv.Kind() != reflect.Slice {
		meta.Error = fmt.Errorf(
			"cant create a content response from '%v'",
			resources,
		)
		return ContentResponse{resources, ContentResponseContext{}, meta}
	}

	if meta.Total == 0 {
		meta.Total = rv.Len()
	}

	if meta.Size == 0 {
		meta.Size = 25
	}

	if meta.Page == 0 {
		meta.Page = 1
	}

	pages := meta.Total / meta.Size

	totalPages := float64(pages)
	if pages <= 0 {
		totalPages = 1
	}

	last := false
	if rv.Len() == 0 || int(totalPages) == meta.Page {
		last = true
	}

	first := false
	if meta.Page == 1 {
		first = true
	}

	context := ContentResponseContext{
		First:      first,
		Last:       last,
		TotalPages: int(totalPages),
	}

	return ContentResponse{
		Content:                resources,
		ContentResponseContext: context,
		ContentResponseMeta:    meta,
	}
}
