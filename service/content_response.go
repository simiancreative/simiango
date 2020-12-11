package service

import (
	"fmt"
	"math"
)

type ContentResponse struct {
	Content    []interface{} `json:"content"`
	First      bool          `json:"first"`
	Last       bool          `json:"last"`
	TotalPages int           `json:"total_pages"`
	ContentResponseMeta
}

type ContentResponseMeta struct {
	Size  int `json:"size"`
	Total int `json:"total"`
	Page  int `json:"current_page"`
}

func ToContentResponse(
	resources []interface{},
	meta ContentResponseMeta,
) ContentResponse {
	if meta.Total == 0 {
		meta.Total = len(resources)
	}

	if meta.Size == 0 {
		meta.Size = 25
	}

	if meta.Page == 0 {
		meta.Page = 1
	}

	pages := meta.Total / meta.Size
	fmt.Println(pages)

	var total_pages float64
	if pages <= 0 {
		total_pages = math.Ceil(total_pages)
	}
	if pages < meta.Total {
		total_pages = 1
	}

	last := false
	if len(resources) == 0 || int(total_pages) == meta.Page {
		last = true
	}

	first := false
	if meta.Page == 1 {
		first = true
	}

	return ContentResponse{
		Content:             resources,
		First:               first,
		Last:                last,
		TotalPages:          int(total_pages),
		ContentResponseMeta: meta,
	}
}
