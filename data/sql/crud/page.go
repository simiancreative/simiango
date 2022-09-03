package crud

import "github.com/simiancreative/simiango/service"

func (m *Model) Page(items interface{}, filters Filters, orders Orders, page int, size int) (*service.ContentResponse, error) {
	ds := m.query(filters, orders)
	count, content := m.pageQueries(ds, page, size)

	var total int

	countQuery, countParams, _ := count.ToSQL()
	if err := m.cx.Get(&total, countQuery, countParams...); err != nil {
		return nil, err
	}

	contentQuery, contentParams, _ := content.ToSQL()
	if err := m.cx.Select(items, contentQuery, contentParams...); err != nil {
		return nil, err
	}

	meta := service.ContentResponseMeta{Size: size, Page: page, Total: total}
	resp := service.ToContentResponse(items, meta)

	return &resp, nil
}
