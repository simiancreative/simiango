package pagination

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/simiancreative/simiango/data/sql"
	"github.com/simiancreative/simiango/service"
)

type Page struct {
	Cx         sql.ConnX
	Attributes string
	From       string
	Where      string
	Order      string
	Page       int
	PageSize   int
}

func (p *Page) Select(items interface{}, params ...interface{}) (*service.ContentResponse, error) {
	countQuery := p.buildCountQuery()
	query := p.buildQuery()

	var total int

	if err := p.Cx.Get(&total, countQuery, params...); err != nil {
		return nil, err
	}

	if err := p.Cx.Select(items, query, params...); err != nil {
		return nil, err
	}

	meta := service.ContentResponseMeta{Size: p.PageSize, Page: p.Page, Total: total}
	resp := service.ToContentResponse(items, meta)

	return &resp, nil
}

func (p *Page) buildCountQuery() string {
	return buildQuery(`
	SELECT COUNT(*)
	FROM {{ .From }}
	{{if .Where -}} WHERE {{.Where}} {{- end}}
	`, p)
}

func (p *Page) buildQuery() string {
	return buildQuery(`
	SELECT {{.Attributes}}
	FROM {{.From}}
	{{if .Where -}} WHERE {{.Where}} {{- end}}
	{{if .Order -}} ORDER BY {{.Order}} {{- end}}
	LIMIT {{ .Limit }}
	OFFSET {{ .Offset }}
	`, p)
}

func (p *Page) Limit() int {
	return p.PageSize
}

func (p *Page) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func buildQuery(queryTpl string, data interface{}) string {
	var tpl bytes.Buffer

	t := template.New("query")
	t.Parse(queryTpl)

	if err := t.Execute(&tpl, data); err != nil {
		fmt.Printf("Build failed %v", err.Error())
		return ""
	}

	return tpl.String()
}
