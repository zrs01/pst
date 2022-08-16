package docb

import "baliance.com/gooxml/document"

type RowBuilder struct {
	config      *Configuration
	doc         *document.Document
	row         document.Row
	cellBuilder []builder
}

func NewRowBuilder(cfg *Configuration, d *document.Document) *RowBuilder {
	return &RowBuilder{doc: d, config: cfg}
}

func (r *RowBuilder) AddCell(nextBuilder func(*CellBuilder)) *RowBuilder {
	c := newCellBuilder(r.config, r.doc)
	r.cellBuilder = append(r.cellBuilder, c)
	return r
}

func (r *RowBuilder) Build() {
	for _, builder := range r.cellBuilder {
		builder.Build()
	}
}
