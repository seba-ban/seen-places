package server

type TemplateInterface interface {
	TemplateName() string
	TemplateInputMap() (map[string]any, error)
}

type TemplateIndex struct {
	GeoJSON *string
}

func (t *TemplateIndex) TemplateName() string {
	return "index.tmpl"
}

func (t *TemplateIndex) TemplateInputMap() map[string]any {
	return map[string]any{
		"GeoJSON": t.GeoJSON,
	}
}
