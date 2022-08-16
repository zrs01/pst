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

func (p *ImageProperty) SetFile(path string) {
	p.FilePath = path
}

func (p *ImageProperty) SetWidth(value float64) {
	p.Width = value
}

func (p *ImageProperty) SetHeight(value float64) {
	p.Height = value
}

func (p *ImageProperty) SetAlignment(align wml.ST_Jc) {
	p.Alignment = align
}
