package docb

import "baliance.com/gooxml/document"

type TableBuilder struct {
	config          *Configuration
	document        *document.Document
	table           *document.Table
	rows            []builder
	cellSpacingAuto bool
	widthPercent    float64
	borders         *Borders
}

func newTableBuilder(cfg *Configuration, doc *document.Document, t document.Table) *TableBuilder {
	return &TableBuilder{config: cfg, document: doc, table: &t}
}

func (t *TableBuilder) AddRow(nextBuilder func(*RowBuilder)) *TableBuilder {
	r := NewRowBuilder(t.config, t.document, t.table.AddRow())
	t.rows = append(t.rows, r)
	nextBuilder(r)
	return t
}

func (t *TableBuilder) SetCellSpacingAuto() *TableBuilder {
	t.cellSpacingAuto = true
	return t
}

func (t *TableBuilder) SetWidthPercent(value float64) *TableBuilder {
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
			b.SetLeft(t.borders.Left.Style, t.borders.Left.Color, t.borders.Left.Thickness)
		}
		if t.borders.InsideHorizontal != nil {
			b.SetInsideHorizontal(t.borders.InsideHorizontal.Style, t.borders.InsideHorizontal.Color, t.borders.InsideHorizontal.Thickness)
		}
		if t.borders.InsideVertical != nil {
			b.SetInsideVertical(t.borders.InsideVertical.Style, t.borders.InsideVertical.Color, t.borders.InsideVertical.Thickness)
		}
	}

	if t.widthPercent > 0 {
		t.table.Properties().SetWidthPercent(t.widthPercent)
	}
	if t.cellSpacingAuto {
		t.table.Properties().SetCellSpacingAuto()
	}

	for _, builder := range t.rows {
		builder.Build()
	}
}
