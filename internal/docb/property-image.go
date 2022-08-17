package docb

import "baliance.com/gooxml/schema/soo/wml"

type ImageProperty struct {
	FilePath  string
	Width     float64
	Height    float64
	Alignment wml.ST_Jc
}

func newImageProperty() *ImageProperty {
	return &ImageProperty{Width: 400.0}
}

func (p *ImageProperty) SetFile(path string) *ImageProperty {
	p.FilePath = path
	return p
}

func (p *ImageProperty) SetWidth(value float64) *ImageProperty {
	p.Width = value
	return p
}

func (p *ImageProperty) SetHeight(value float64) *ImageProperty {
	p.Height = value
	return p
}

func (p *ImageProperty) SetAlignment(align wml.ST_Jc) *ImageProperty {
	p.Alignment = align
	return p
}
