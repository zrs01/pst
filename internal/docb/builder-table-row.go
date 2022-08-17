package docb

import "baliance.com/gooxml/document"

type RowBuilder struct {
	config      *Configuration
	doc         *document.Document
	row         *document.Row
	cellBuilder []builder
}

func NewRowBuilder(cfg *Configuration, d *document.Document, r document.Row) *RowBuilder {
	return &RowBuilder{config: cfg, doc: d, row: &r}
}

func (r *RowBuilder) AddCell(nextBuilder func(*CellBuilder)) *RowBuilder {
	c := newCellBuilder(r.config, r.doc, r.row.AddCell())
	r.cellBuilder = append(r.cellBuilder, c)
	nextBuilder(c)
	return r
}

func (r *RowBuilder) Build() {
	for _, builder := range r.cellBuilder {
		builder.Build()
	}
}
