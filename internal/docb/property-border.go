package docb

import (
	"baliance.com/gooxml/color"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
)

type BorderStyle struct {
	Style     wml.ST_Border
	Color     color.Color
	Thickness measurement.Distance
}

type Borders struct {
	Bottom           *BorderStyle
	Left             *BorderStyle
	Right            *BorderStyle
	Top              *BorderStyle
	InsideHorizontal *BorderStyle
	InsideVertical   *BorderStyle
}

func (b *Borders) SetBorderBottom(t wml.ST_Border, c color.Color, thickness measurement.Distance) {
	b.Bottom = &BorderStyle{Style: t, Color: c, Thickness: thickness}
}

func (b *Borders) SetBorderLeft(t wml.ST_Border, c color.Color, thickness measurement.Distance) {
	b.Left = &BorderStyle{Style: t, Color: c, Thickness: thickness}
}

func (b *Borders) SetBorderRight(t wml.ST_Border, c color.Color, thickness measurement.Distance) {
	b.Right = &BorderStyle{Style: t, Color: c, Thickness: thickness}
}

func (b *Borders) SetBorderTop(t wml.ST_Border, c color.Color, thickness measurement.Distance) {
	b.Top = &BorderStyle{Style: t, Color: c, Thickness: thickness}
}

func (b *Borders) SetBorderInsideHorizontal(t wml.ST_Border, c color.Color, thickness measurement.Distance) {
	b.InsideHorizontal = &BorderStyle{Style: t, Color: c, Thickness: thickness}
}

func (b *Borders) SetBorderInsideVertical(t wml.ST_Border, c color.Color, thickness measurement.Distance) {
	b.InsideVertical = &BorderStyle{Style: t, Color: c, Thickness: thickness}
}

func (b *Borders) SetBorderAll(t wml.ST_Border, c color.Color, thickness measurement.Distance) {
	b.Bottom = &BorderStyle{Style: t, Color: c, Thickness: thickness}
	b.Left = &BorderStyle{Style: t, Color: c, Thickness: thickness}
	b.Right = &BorderStyle{Style: t, Color: c, Thickness: thickness}
	b.Top = &BorderStyle{Style: t, Color: c, Thickness: thickness}
	b.InsideHorizontal = &BorderStyle{Style: t, Color: c, Thickness: thickness}
	b.InsideVertical = &BorderStyle{Style: t, Color: c, Thickness: thickness}
}

func (b *Borders) SetBorderTopBottom(t wml.ST_Border, c color.Color, thickness measurement.Distance) {
	b.Bottom = &BorderStyle{Style: t, Color: c, Thickness: thickness}
	b.Top = &BorderStyle{Style: t, Color: c, Thickness: thickness}
}
