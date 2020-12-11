package service

type ContentResponse struct {
	Content          []interface{} `json:"content"`
	First            bool          `json:"first"`
	Last             bool          `json:"last"`
	Number           int           `json:"number"`
	NumberOfElements int           `json:"number_of_elements"`
	Size             int           `json:"size"`
	Total            int           `json:"total"`
	TotalPages       int           `json:"total_pages"`
}

func ToContentResponse(resources []interface{}) ContentResponse {
	return ContentResponse{Content: resources}
}
