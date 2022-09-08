package crud

import (
	"strings"

	"github.com/simiancreative/simiango/service"
)

func (m *Model) PageFromReq(dest interface{}, req service.Req) (interface{}, error) {
	order := OrderFromReq(req)
	filters := FiltersFromReq(req)

	return m.Page(
		dest,
		filters,
		order,
		req.Params.GetWithFallback("page", "1").AsInt(),
		req.Params.GetWithFallback("size", "5").AsInt(),
	)
}

func FiltersFromReq(req service.Req) Filters {
	filters := Filters{}

	for key, value := range req.Params.ValuesMap() {
		if key == "order" {
			continue
		}

		filters[key] = value
	}

	return filters
}

func OrderFromReq(req service.Req) Order {
	order := Order{}

	item, exists := req.Params.Get("order")

	if !exists {
		return order
	}

	for _, value := range item.Values {
		parts := strings.Split(value, ",")

		if len(parts) == 0 {
			continue
		}

		if len(parts) == 1 {
			order[value] = "asc"
		}

		if len(parts) > 1 {
			order[parts[0]] = parts[1]
		}
	}

	return order
}
