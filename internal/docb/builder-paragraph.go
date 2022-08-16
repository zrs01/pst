package docb

import (
	"path"

	"baliance.com/gooxml/common"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
	"github.com/shomali11/util/xstrings"
	"github.com/sirupsen/logrus"
)

type ParagraphBuilder struct {
	config    *Configuration
	document  *document.Document
	style     string
	text      []string
	alignment wml.ST_Jc
	pageBreak bool
	lineBreak bool
	images    []*ImageProperty
	// 	imageFilePath string
	// 	imageWidth    int
}

func newParagraphBuilder(cfg *Configuration, d *document.Document) *ParagraphBuilder {
	return &ParagraphBuilder{document: d, config: cfg}
}

func (p *ParagraphBuilder) SetStyle(s string) *ParagraphBuilder {
	p.style = s
	return p
}

func (p *ParagraphBuilder) SetText(s interface{}) *ParagraphBuilder {
	p.text = toStrArray(s)
	return p
}

func (p *ParagraphBuilder) SetLineBreak() *ParagraphBuilder {
	p.lineBreak = true
	return p
}

func (p *ParagraphBuilder) SetPageBreak() *ParagraphBuilder {
	p.pageBreak = true
	return p
}

func (p *ParagraphBuilder) SetAlignment(align wml.ST_Jc) *ParagraphBuilder {
	p.alignment = align
	return p
}

func (p *ParagraphBuilder) AddImage(set func(*ImageProperty)) *ParagraphBuilder {
	i := newImageProperty()
	set(i)
	p.images = append(p.images, i)
	return p
}

func (p *ParagraphBuilder) Build() {
	paragraph := p.document.AddParagraph()

	if xstrings.IsNotBlank(p.style) {
		paragraph.SetStyle(p.style)
	}

	if p.alignment != wml.ST_JcUnset {
		paragraph.Properties().SetAlignment(p.alignment)
	}

	for i, s := range p.text {
		run := paragraph.AddRun()
		run.AddText(s)
		if i == len(p.text)-1 && p.lineBreak {
			run.AddBreak()
		}
	}

	for _, img := range p.images {
		// image should relative to input file
		imgFilePath := path.Join(p.config.ImagePath, img.FilePath)
		imgFile, err := common.ImageFromFile(imgFilePath)
		if err != nil {
			logrus.Warn(err)
			return
		}
		iref, err := p.document.AddImage(imgFile)
		if err != nil {
			logrus.Warn(err)
			return
		}
		para := p.document.AddParagraph()
		if img.Alignment != wml.ST_JcUnset {
			para.Properties().SetAlignment(img.Alignment)
		}
		inl, err := para.AddRun().AddDrawingInline(iref)
		if err != nil {
			logrus.Warn(err)
			return
		}
		w := img.Width
		h := w / float64(imgFile.Size.X) * float64(imgFile.Size.Y)
		inl.SetSize(measurement.Distance(w*measurement.Point), measurement.Distance(h*measurement.Point))
	}
	if p.pageBreak {
		paragraph.Properties().AddSection(wml.ST_SectionMarkNextPage)

		// have to set the header on the section as well, Word doesn't automatically associate them
		// http://github.com/unidoc/unioffice/issues/173

		// sect := p.para.Properties().AddSection(wml.ST_SectionMarkNextPage)
		// sect.SetHeader(hdr, wml.ST_HdrFtrDefault)
		// sect.SetFooter(ftr, wml.ST_HdrFtrDefault)
	}
}
