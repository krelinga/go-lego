package main

import (
	"bytes"
	"fmt"
)

var equalTemplateBodies = map[string]string{
	"comparer": `func ({{.ReceiverVar}} {{.ReceiverType}}) Equal(other {{.ReceiverType}}) bool {
	return {{.Prefix}}EqualComparer({{.ReceiverVar}}, other)
}
`,
	"viewer": `func ({{.ReceiverVar}} {{.ReceiverType}}) Equal(other {{.ReceiverType}}) bool {
	return {{.Prefix}}EqualViewer({{.ReceiverVar}}, other)
}
`,
	"slice": `func ({{.ReceiverVar}} {{.ReceiverType}}) Equal(other {{.ReceiverType}}) bool {
	return {{.Prefix}}EqualSlice({{.ReceiverVar}}, other)
}
`,
	"slicego": `func ({{.ReceiverVar}} {{.ReceiverType}}) Equal(other {{.ReceiverType}}) bool {
	return {{.Prefix}}EqualSliceGo({{.ReceiverVar}}, other)
}
`,
}

type equalTemplateData struct {
	ReceiverVar  string
	ReceiverType string
	Prefix       string // "lego." or "" when inside the lego package
}

func generateEqual(cfg config) (string, error) {
	tmpl := equalTemplates.Lookup(cfg.Equal)
	if tmpl == nil {
		return "", fmt.Errorf("unknown equal strategy %q", cfg.Equal)
	}

	prefix := "lego."
	if cfg.Package == "lego" {
		prefix = ""
	}

	data := equalTemplateData{
		ReceiverVar:  cfg.ReceiverVar,
		ReceiverType: cfg.receiverType(),
		Prefix:       prefix,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %q: %w", cfg.Equal, err)
	}
	return buf.String(), nil
}
