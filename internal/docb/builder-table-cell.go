package docb

import (
	"strings"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
	"github.com/shomali11/util/xconditions"
)

type CellBuilder struct {
	config          *Configuration
	document        *document.Document
	cell            *document.Cell
	builder         []builder
	fontFamily      string
	fontSize        int
	bold            bool
	width           int
	colspan         int
	bullet          *document.NumberingDefinition
	text            []string
	backgroundColor *color.Color
	borders         *Borders
	alignment       wml.ST_Jc
}

func newCellBuilder(cfg *Configuration, doc *document.Document, c document.Cell) *CellBuilder {
	return &CellBuilder{config: cfg, document: doc, cell: &c}
}

func (c *CellBuilder) AddParagraph(nextBuilder ...func(*ParagraphBuilder)) *CellBuilder {
	p := newParagraphBuilder(c.config, c.document, c.cell.AddParagraph())
	c.builder = append(c.builder, p)
	if len(nextBuilder) > 0 {
		nextBuilder[0](p)
	}
	return c
}

func (c *CellBuilder) SetFontFamily(ff string) *CellBuilder {
	c.fontFamily = ff
	return c
}

func (c *CellBuilder) SetFontSize(fs int) *CellBuilder {
	c.fontSize = fs
	return c
}

func (c *CellBuilder) SetBold() *CellBuilder {
	c.bold = true
	return c
}

func (c *CellBuilder) SetWidth(w int) *CellBuilder {
	c.width = w
	return c
}

func (c *CellBuilder) SetColspan(cols int) *CellBuilder {
	c.colspan = cols
	return c
}

func (c *CellBuilder) SetText(v interface{}) *CellBuilder {
	c.text = toStrArray(v)
	return c
}

func (c *CellBuilder) SetBullet(b *document.NumberingDefinition) *CellBuilder {
	c.bullet = b
	return c
}

func (c *CellBuilder) SetBackgroundColor(color color.Color) *CellBuilder {
	c.backgroundColor = &color
	return c
}

func (c *CellBuilder) SetBorders(callback func(*Borders)) *CellBuilder {
	c.borders = &Borders{}
	callback(c.borders)
	return c
}

func (c *CellBuilder) SetAlignment(align wml.ST_Jc) *CellBuilder {
	c.alignment = align
	return c
}

func (c *CellBuilder) Build() {
	if c.borders != nil {
		b := c.cell.Properties().Borders()
		if c.borders.Top != nil {
			b.SetTop(c.borders.Top.Style, c.borders.Top.Color, c.borders.Top.Thickness)
		}
		if c.borders.Right != nil {
			b.SetRight(c.borders.Right.Style, c.borders.Right.Color, c.borders.Right.Thickness)
		}
		if c.borders.Bottom != nil {
			b.SetBottom(c.borders.Bottom.Style, c.borders.Bottom.Color, c.borders.Bottom.Thickness)
		}
		if c.borders.Left != nil {
			b.SetLeft(c.borders.Left.Style, c.borders.Left.Color, c.borders.Left.Thickness)
		}
		if c.borders.InsideHorizontal != nil {
			b.SetInsideHorizontal(c.borders.InsideHorizontal.Style, c.borders.InsideHorizontal.Color, c.borders.InsideHorizontal.Thickness)
		}
		if c.borders.InsideVertical != nil {
			b.SetInsideVertical(c.borders.InsideVertical.Style, c.borders.InsideVertical.Color, c.borders.InsideVertical.Thickness)
		}
	}

	if c.width > 0 {
		c.cell.Properties().SetWidthPercent(float64(c.width))
	}

	if c.colspan > 0 {
		c.cell.Properties().SetColumnSpan(c.colspan)
	}

	if c.backgroundColor != nil {
		c.cell.Properties().SetShading(wml.ST_ShdSolid, *c.backgroundColor, color.Auto)
	}

	if len(c.text) == 0 {
		p := c.cell.AddParagraph()
		p.AddRun().AddText("")
	} else {
		for _, t := range c.text {
			p := c.cell.AddParagraph()
			if c.alignment != wml.ST_JcUnset {
				p.Properties().SetAlignment(c.alignment)
			}
			tab := ""
			if c.bullet != nil && len(c.text) > 1 {
				p.SetNumberingLevel(0)
				p.SetNumberingDefinition(*c.bullet)
				tab = "\t"
			}
			lines := strings.Split(t, "\\n")
			for i, line := range lines {
				line := strings.ReplaceAll(line, "\\t", "\t")
				run := p.AddRun()
				run.Properties().SetBold(c.bold)
				run.Properties().SetFontFamily(xconditions.IfThenElse(c.fontFamily != "", c.fontFamily, c.config.FontFamily).(string))
				run.Properties().SetSize(measurement.Distance(xconditions.IfThenElse(c.fontSize > 0, c.fontSize, c.config.FontSize).(int)))
				if i == 0 {
					run.AddText(tab + line)
				} else {
					// multi-line text
					run.AddBreak()
					run.AddText(line)
				}
			}
		}
	}

	for _, builder := range c.builder {
		builder.Build()
	}
}
