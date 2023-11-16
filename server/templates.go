package server

import "github.com/seba-ban/seen-places/queries"

type TemplateInterface interface {
	TemplateName() string
	TemplateInputMap() (map[string]any, error)
}

type TemplateIndex struct {
}

func (t *TemplateIndex) TemplateName() string {
	return "index.html"
}

func (t *TemplateIndex) TemplateInputMap() map[string]any {
	return map[string]any{}
}

type TemplateDataSourceEls struct {
	DataSources []queries.GetDataSourcesFromPolygonRow
}

func (t *TemplateDataSourceEls) TemplateName() string {
	return "dataSourceEls.tmpl"
}

func (t *TemplateDataSourceEls) TemplateInputMap() map[string]any {
	return map[string]any{
		"DataSources": t.DataSources,
	}
}
