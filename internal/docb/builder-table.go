package docb

import "baliance.com/gooxml/document"

type TableBuilder struct {
	config       *Configuration
	document     *document.Document
	table        *document.Table
	rows         []builder
	widthPercent int
	borders      *Borders
}

func newTableBuilder(cfg *Configuration, doc *document.Document) *TableBuilder {
	t := doc.AddTable()
	return &TableBuilder{document: doc, table: &t, config: cfg}
}

func (t *TableBuilder) AddRow(nextBuilder func(*RowBuilder)) *TableBuilder {
	r := NewRowBuilder(t.config, t.document)
	t.rows = append(t.rows, r)
	nextBuilder(r)
	return t
}

func (t *TableBuilder) SetWidthPercent(value int) *TableBuilder {
	t.widthPercent = value
	return t
}

func (t *TableBuilder) SetBorders(callback func(*Borders)) *TableBuilder {
	t.borders = &Borders{}
	callback(t.borders)
	return t
}

func (t *TableBuilder) Build() {
	if t.borders != nil {
		b := t.table.Properties().Borders()
		if t.borders.Top != nil {
			b.SetTop(t.borders.Top.Style, t.borders.Top.Color, t.borders.Top.Thickness)
		}
		if t.borders.Right != nil {
			b.SetRight(t.borders.Right.Style, t.borders.Right.Color, t.borders.Right.Thickness)
		}
		if t.borders.Bottom != nil {
			b.SetBottom(t.borders.Bottom.Style, t.borders.Bottom.Color, t.borders.Bottom.Thickness)
		}
		if t.borders.Left != nil {
			b.SetLeft(t.borders.Left.Style, t.borders.Left.Color, t.borders.Bottom.Thickness)
		}
	}

	for _, builder := range t.rows {
		builder.Build()
	}
}
