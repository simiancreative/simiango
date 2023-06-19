package pagination

import (
	"github.com/simiancreative/simiango/service"
)

func (p *Page) NamedSelect(items interface{}, arg interface{}) (*service.ContentResponse, error) {
	countQuery := p.buildCountQuery()
	query := p.buildQuery()

	var total int

	nstmt, err := p.Cx.PrepareNamed(countQuery)
	if nstmt != nil {
		defer nstmt.Close()
	}

	if err != nil {
		return nil, err
	}

	if err := nstmt.Get(&total, arg); err != nil {
		return nil, err
	}

	nstmt, err = p.Cx.PrepareNamed(query)
	if nstmt != nil {
		defer nstmt.Close()
	}

	if err != nil {
		return nil, err
	}

	if err := nstmt.Select(items, arg); err != nil {
		return nil, err
	}

	meta := service.ContentResponseMeta{Size: p.PageSize, Page: p.Page, Total: total}
	resp := service.ToContentResponse(items, meta)

	return &resp, nil
}
