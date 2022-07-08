package converter

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"baliance.com/gooxml"
	"baliance.com/gooxml/color"
	"baliance.com/gooxml/common"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xconditions"
	"github.com/shomali11/util/xstrings"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

// Reference:
// https://github.com/bollwarm/gooxml   Office Open XML documents utility
// https://github.com/shomali11/util		A group of generic useful utility functions

type builder struct {
	cfile string // config file name
	ifile string // input file name
	ofile string // output file name
	dfile string // .docx file name
}

var cfg *Config
var ifilePath string

func Build(cfile, ifile, ofile string, tfile string) error {
	ncfg, err := NewConfig(cfile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the configuration file %s", cfile)
	}
	cfg = ncfg

	ifilePath = path.Dir(ifile)
	b := &builder{ifile: ifile, ofile: ofile, dfile: tfile}
	return b.buildSpec()
}

func (b *builder) buildSpec() error {
	var doc *document.Document
	if xstrings.IsBlank(b.dfile) {
		doc = document.New()
		// set to A4 size (https://stackoverflow.com/questions/57581695/detecting-and-setting-paper-size-in-word-js-api-or-ooxml)
		// doc.X().Body.SectPr = wml.NewCT_SectPr()
		// doc.X().Body.SectPr.PgSz = &wml.CT_PageSz{
		// 	WAttr: &sharedTypes.ST_TwipsMeasure{
		// 		ST_UnsignedDecimalNumber: gooxml.Uint64(uint64(11906)),
		// 	},
		// 	HAttr: &sharedTypes.ST_TwipsMeasure{
		// 		ST_UnsignedDecimalNumber: gooxml.Uint64(uint64(16838)),
		// 	},
		// }
	} else {
		ndoc, err := document.Open(b.dfile)
		if err != nil {
			return eris.Wrapf(err, "failed to open the document %s", b.dfile)
		}
		doc = ndoc
	}
	fixBulletIndentation(doc)

	// support multiple files
	files, err := b.resolveIFile(b.ifile)
	if err != nil {
		return eris.Wrap(err, "failed to resolve the source file")
	}
	for _, file := range *files {
		data, err := b.loadData(file)
		if err != nil {
			return eris.Wrap(err, "failed to load the file")
		}

		// header
		NewParaBuilder(doc).SetStyle("Heading1").SetText("PROGRAM DESCRIPTION").Build()

		for _, module := range data.Modules {
			NewParaBuilder(doc).SetStyle("Heading2").SetText(module.Name).Build()
			for _, feature := range module.Features {
				doc.AddParagraph()
				NewParaBuilder(doc).SetStyle("Heading3").SetText(feature.Name).Build()
				b.buildFeature(doc, &feature)
			}
			// if i < len(data.Modules)-1 {
			NewParaBuilder(doc).SetPageBreak().Build()
			// }
		}

	}

	doc.SaveToFile(b.ofile)
	return nil
}

func (b *builder) buildFeature(doc *document.Document, feature *Feature) error {
	var createTable = func(addLineBreak bool) document.Table {
		if addLineBreak {
			doc.AddParagraph()
		}
		table := doc.AddTable()
		table.Properties().SetWidthPercent(100)
		borders := table.Properties().Borders()
		borders.SetTop(wml.ST_BorderSingle, color.Auto, 0.5*measurement.Point)
		borders.SetBottom(wml.ST_BorderSingle, color.Auto, 0.5*measurement.Point)
		borders.SetLeft(wml.ST_BorderSingle, color.Auto, 0.5*measurement.Point)
		borders.SetRight(wml.ST_BorderSingle, color.Auto, 0.5*measurement.Point)
		return table
	}
	var splitScenarioWord = func(s string) (string, string) {
		if xstrings.IsBlank(s) {
			return "", ""
		}
		s = strings.TrimSpace(s)
		parts := strings.Split(s, " ")
		if funk.Contains([]string{"given", "when", "then", "and", "but"}, strings.ToLower(parts[0])) {
			return parts[0], s[len(parts[0])+1:]
		}
		return "", s
	}

	/* ---------------------------------- main ---------------------------------- */

	w := 20
	gray := color.FromHex("ced4da")
	lightgray := color.FromHex("e9ecef")
	nd := doc.Numbering.Definitions()[0]
	rb := NewRowBuilder(doc)

	/* --------------------------------- PROGRAM -------------------------------- */
	tbMain := createTable(false)
	rb.Reset(tbMain.AddRow()).AddCell(
		rb.NewCellBuilder().SetText("Program ID").SetWidth(w).SetBold(true).SetBorderTopBottom(),
		rb.NewCellBuilder().SetText(feature.Id).SetBorderTopBottom()).Build()
	rb.Reset(tbMain.AddRow()).AddCell(
		rb.NewCellBuilder().SetText("Mode").SetWidth(w).SetBold(true).SetBorderTopBottom(),
		rb.NewCellBuilder().SetText(feature.Mode).SetBorderTopBottom()).Build()
	rb.Reset(tbMain.AddRow()).AddCell(
		rb.NewCellBuilder().SetText("Program Name").SetWidth(w).SetBold(true).SetBorderTopBottom(),
		rb.NewCellBuilder().SetText(feature.Name).SetBorderTopBottom()).Build()
	rb.Reset(tbMain.AddRow()).AddCell(
		rb.NewCellBuilder().SetText("Description").SetWidth(w).SetBold(true).SetBorderTopBottom(),
		rb.NewCellBuilder().SetText(feature.Desc).SetBorderTopBottom()).Build()

	if (feature.Env.Sources != nil && !reflect.ValueOf(feature.Env.Sources).IsZero()) ||
		(feature.Env.Languages != nil && !reflect.ValueOf(feature.Env.Languages).IsZero()) {
		rb.Reset(tbMain.AddRow()).AddCell(
			rb.NewCellBuilder().SetText("Program Environment:").SetBold(true).SetColspan(2).SetBackgroundColor(gray).SetBorderTopBottom()).Build()
		if feature.Env.Sources != nil && !reflect.ValueOf(feature.Env.Sources).IsZero() {
			rb.Reset(tbMain.AddRow()).AddCell(
				rb.NewCellBuilder().SetText("Program Source").SetWidth(w).SetBold(true).SetBorderTopBottom(),
				rb.NewCellBuilder().SetText(feature.Env.Sources).SetBorderTopBottom()).Build()
		}
		if feature.Env.Languages != nil && !reflect.ValueOf(feature.Env.Languages).IsZero() {
			rb.Reset(tbMain.AddRow()).AddCell(
				rb.NewCellBuilder().SetText("Language").SetWidth(w).SetBold(true).SetBorderTopBottom(),
				rb.NewCellBuilder().SetText(feature.Env.Languages).SetBorderTopBottom()).Build()
		}
	}

	/* -------------------------------- RESOURCE -------------------------------- */
	if len(feature.Resources) > 0 {
		tbRes := createTable(true)
		rb.Reset(tbRes.AddRow()).AddCell(
			rb.NewCellBuilder().SetText("Resources:").SetBold(true).SetColspan(2).SetBackgroundColor(gray).SetBorderTopBottom()).Build()
		rb.Reset(tbRes.AddRow()).AddCell(
			rb.NewCellBuilder().SetText("Table/File").SetWidth(w).SetBold(true).SetBackgroundColor(lightgray).SetBorderTopBottom(),
			rb.NewCellBuilder().SetText("Usage").SetBold(true).SetBackgroundColor(lightgray).SetBorderTopBottom()).Build()
		for _, res := range feature.Resources {
			rb.Reset(tbRes.AddRow()).AddCell(
				rb.NewCellBuilder().SetText(res.Name).SetWidth(w).SetBorderTopBottom(),
				rb.NewCellBuilder().SetText(res.Usage).SetBorderTopBottom()).Build()
		}

		/* --------------------------------- SCREEN --------------------------------- */
		if len(feature.Screens) > 0 {
			tbScr := createTable(true)
			rb.Reset(tbScr.AddRow()).AddCell(
				rb.NewCellBuilder().SetText("Screen:").SetBold(true).SetColspan(2).SetBackgroundColor(gray).SetBorderTopBottom()).Build()
			for i, scr := range feature.Screens {
				if i > 0 && xstrings.IsNotBlank(scr.Image.File) {
					tbScr = createTable(false)
				}
				rb.Reset(tbScr.AddRow()).AddCell(
					rb.NewCellBuilder().SetText("Screen ID").SetWidth(w).SetBold(true).SetBackgroundColor(lightgray).SetBorderTopBottom(),
					rb.NewCellBuilder().SetText("Name").SetBold(true).SetBackgroundColor(lightgray).SetBorderTopBottom()).Build()

				rb.Reset(tbScr.AddRow()).AddCell(
					rb.NewCellBuilder().SetText(scr.Id).SetWidth(w).SetBorderTopBottom(),
					rb.NewCellBuilder().SetText(scr.Name).SetBorderTopBottom()).Build()
				if xstrings.IsNotBlank(scr.Image.File) {
					NewParaBuilder(doc).SetImage(scr.Image.File).SetImageWidth(scr.Image.Width).Build()
				}
			}
		}

		/* ---------------------------------- INPUT --------------------------------- */
		if len(feature.Input) > 0 {
			tbInp := createTable(true)
			rb.Reset(tbInp.AddRow()).AddCell(
				rb.NewCellBuilder().SetText("Input:").SetBold(true).SetColspan(2).SetBackgroundColor(gray).SetBorderTopBottom()).Build()
			for i, inp := range feature.Input {
				rb.Reset(tbInp.AddRow()).AddCell(
					rb.NewCellBuilder().SetText(fmt.Sprintf("%d. %s", i+1, inp.Name)).SetBold(true).SetColspan(2).SetBackgroundColor(lightgray).SetBorderTopBottom(),
				).Build()
				if inp.Fields != nil && !reflect.ValueOf(inp.Fields).IsZero() {
					rb.Reset(tbInp.AddRow()).AddCell(
						rb.NewCellBuilder().SetText("Fields").SetBold(true).SetWidth(w),
						rb.NewCellBuilder().SetText(inp.Fields),
					).Build()
				}
				if inp.Constraints != nil && !reflect.ValueOf(inp.Constraints).IsZero() {
					rb.Reset(tbInp.AddRow()).AddCell(
						rb.NewCellBuilder().SetText("Constraints").SetBold(true).SetWidth(w),
						rb.NewCellBuilder().SetText(inp.Constraints).SetBullet(&nd),
					).Build()
				}
				if inp.Remarks != nil && !reflect.ValueOf(inp.Remarks).IsZero() {
					rb.Reset(tbInp.AddRow()).AddCell(
						rb.NewCellBuilder().SetText("Remarks").SetBold(true).SetWidth(w),
						rb.NewCellBuilder().SetText(inp.Remarks).SetBullet(&nd),
					).Build()
				}
			}
		}

		/* -------------------------------- SCENARIO -------------------------------- */
		if len(feature.Scenarios) > 0 {
			tbScn := createTable(true)
			rb.Reset(tbScn.AddRow()).AddCell(
				rb.NewCellBuilder().SetText("Scenarios:").SetBold(true).SetColspan(2).SetBackgroundColor(gray).SetBorderTopBottom()).Build()
			for i, scn := range feature.Scenarios {
				rb.Reset(tbScn.AddRow()).AddCell(
					rb.NewCellBuilder().SetText(fmt.Sprintf("%d. %s", i+1, scn.Name)).SetBold(true).SetBackgroundColor(lightgray).SetColspan(2).SetBorderTopBottom(),
				).Build()

				for _, action := range scn.Desc {
					keyword, others := splitScenarioWord(action)
					rb.Reset(tbScn.AddRow()).AddCell(
						rb.NewCellBuilder().SetText(keyword).SetBold(true).SetAlignment(wml.ST_JcLeft).SetWidth(10),
						rb.NewCellBuilder().SetText(others),
					).Build()
				}
			}
		}

		// Remarks
		if feature.Remarks != nil && !reflect.ValueOf(feature.Remarks).IsZero() {
			tbRmk := createTable(true)
			rb.Reset(tbRmk.AddRow()).AddCell(
				rb.NewCellBuilder().SetText("Remarks:").SetBold(true).SetBackgroundColor(gray).SetBorderTopBottom()).Build()
			rb.Reset(tbRmk.AddRow()).AddCell(
				rb.NewCellBuilder().SetText(feature.Remarks).SetBullet(&nd),
			).Build()
		}

	}
	return nil
}

func (b *builder) loadData(file string) (*ProgSpec, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, eris.Wrap(err, "failed to read the file")
	}
	var d ProgSpec
	if err := yaml.Unmarshal(yamlFile, &d); err != nil {
		return nil, eris.Wrapf(err, "failed to unmarshal the file %s", file)
	}
	return &d, nil
}

func (b *builder) resolveIFile(ifile string) (*[]string, error) {
	files := []string{}
	ifiles := strings.Split(ifile, ",")
	for _, ifile := range ifiles {
		fs, err := filepath.Glob(ifile)
		if err != nil {
			return nil, eris.Wrap(err, "failed to glob the file")
		}
		if len(fs) == 0 {
			return nil, eris.Errorf("no such file: %s", ifile)
		}
		files = append(files, fs...)
	}
	sort.Strings(files)
	return &files, nil
}

/* -------------------------------- UTILITIES ------------------------------- */
func toStrArray(v interface{}) []string {
	text := []string{}
	switch t := v.(type) {
	case string:
		text = []string{t}
	case []string:
		text = append(text, t...)
	case []interface{}:
		for _, value := range t {
			content := strings.TrimSpace(value.(string))
			text = append(text, content)
		}
	}
	return text
}

func fixBulletIndentation(doc *document.Document) {
	const indentStart = 400
	const indentDelta = 400
	const hangingIndent = 360

	abns := doc.Numbering.X().AbstractNum[0]
	for i, lvl := range abns.Lvl {
		indent := int64(i*indentDelta + indentStart)
		lvl.PPr.Ind.LeftAttr.Int64 = gooxml.Int64(indent)
		lvl.PPr.Ind.HangingAttr.ST_UnsignedDecimalNumber = gooxml.Uint64(uint64(hangingIndent))
	}
}

/* -------------------------------------------------------------------------- */
/*                             PARAGRAPHS BUILDER                             */
/* -------------------------------------------------------------------------- */

type ParaBuilder struct {
	doc        *document.Document
	para       *document.Paragraph
	style      string
	text       []string
	pageBreak  bool
	lineBreak  bool
	imagePath  string
	imageWidth int
}

func NewParaBuilder(d *document.Document) *ParaBuilder {
	p := d.AddParagraph()
	return &ParaBuilder{doc: d, para: &p}
}
func (p *ParaBuilder) SetStyle(s string) *ParaBuilder {
	p.style = s
	return p
}
func (p *ParaBuilder) SetText(s interface{}) *ParaBuilder {
	p.text = toStrArray(s)
	return p
}
func (p *ParaBuilder) SetLineBreak() *ParaBuilder {
	p.lineBreak = true
	return p
}
func (p *ParaBuilder) SetPageBreak() *ParaBuilder {
	p.pageBreak = true
	return p
}
func (p *ParaBuilder) SetImage(path string) *ParaBuilder {
	p.imagePath = path
	return p
}
func (p *ParaBuilder) SetImageWidth(width int) *ParaBuilder {
	p.imageWidth = width
	return p
}
func (p *ParaBuilder) Build() {
	if xstrings.IsNotBlank(p.style) {
		p.para.SetStyle(p.style)
	}
	for i, s := range p.text {
		run := p.para.AddRun()
		run.AddText(s)
		if i == len(p.text)-1 && p.lineBreak {
			run.AddBreak()
		}
	}
	if xstrings.IsNotBlank(p.imagePath) {
		// image should relative to input file
		imgPath := path.Join(ifilePath, p.imagePath)
		img, err := common.ImageFromFile(imgPath)
		if err != nil {
			logrus.Warn(err)
			return
		}
		iref, err := p.doc.AddImage(img)
		if err != nil {
			logrus.Warn(err)
			return
		}
		para := p.doc.AddParagraph()
		para.Properties().SetAlignment(wml.ST_JcCenter)
		inl, err := para.AddRun().AddDrawingInline(iref)
		if err != nil {
			logrus.Warn(err)
			return
		}
		w := 400.0
		if p.imageWidth > 0 {
			w = float64(p.imageWidth)
		}
		h := w / float64(img.Size.X) * float64(img.Size.Y)
		inl.SetSize(measurement.Distance(w*measurement.Point), measurement.Distance(h*measurement.Point))
	}
	if p.pageBreak {
		p.para.Properties().AddSection(wml.ST_SectionMarkNextPage)

		// have to set the header on the section as well, Word doesn't automatically associate them
		// http://github.com/unidoc/unioffice/issues/173

		// sect := p.para.Properties().AddSection(wml.ST_SectionMarkNextPage)
		// sect.SetHeader(hdr, wml.ST_HdrFtrDefault)
		// sect.SetFooter(ftr, wml.ST_HdrFtrDefault)
	}
}

/* -------------------------------------------------------------------------- */
/*                                 ROW BUILDER                                */
/* -------------------------------------------------------------------------- */

type RowBuilder struct {
	doc         *document.Document
	row         document.Row
	cellBuilder []*CellBuilder
}

func NewRowBuilder(d *document.Document) *RowBuilder {
	return &RowBuilder{doc: d}
}
func (r *RowBuilder) Reset(row document.Row) *RowBuilder {
	r.row = row
	r.cellBuilder = []*CellBuilder{}
	return r
}
func (r *RowBuilder) NewCellBuilder() *CellBuilder {
	return &CellBuilder{doc: r.doc}
}
func (r *RowBuilder) AddCell(attrs ...*CellBuilder) *RowBuilder {
	r.cellBuilder = append(r.cellBuilder, attrs...)
	return r
}
func (r *RowBuilder) Build() {
	for _, attr := range r.cellBuilder {
		nc := r.row.AddCell()
		attr.SetCell(&nc).Build()
	}
}

/* -------------------------------------------------------------------------- */
/*                                CELL BUILDER                                */
/* -------------------------------------------------------------------------- */

type CellBuilder struct {
	doc                    *document.Document
	cell                   *document.Cell
	fontFamily             string
	fontSize               int
	bold                   bool
	width                  int
	colspan                int
	bullet                 *document.NumberingDefinition
	text                   []string
	backgroundColor        *color.Color
	borderBottom           bool
	borderLeft             bool
	borderRight            bool
	borderTop              bool
	borderInsideHorizontal bool
	borderInsideVertical   bool
	alignment              wml.ST_Jc
	imagePath              string
}

func (c *CellBuilder) SetCell(dc *document.Cell) *CellBuilder {
	c.cell = dc
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
func (c *CellBuilder) SetBold(b bool) *CellBuilder {
	c.bold = b
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
func (c *CellBuilder) SetBorderBottom() *CellBuilder {
	c.borderBottom = true
	return c
}
func (c *CellBuilder) SetBorderLeft() *CellBuilder {
	c.borderLeft = true
	return c
}
func (c *CellBuilder) SetBorderRight() *CellBuilder {
	c.borderRight = true
	return c
}
func (c *CellBuilder) SetBorderTop() *CellBuilder {
	c.borderTop = true
	return c
}
func (c *CellBuilder) SetBorderInsideHorizontal() *CellBuilder {
	c.borderInsideHorizontal = true
	return c
}
func (c *CellBuilder) SetBorderInsideVertical() *CellBuilder {
	c.borderInsideVertical = true
	return c
}
func (c *CellBuilder) SetBorderAll() *CellBuilder {
	c.borderBottom = true
	c.borderLeft = true
	c.borderRight = true
	c.borderTop = true
	c.borderInsideHorizontal = true
	c.borderInsideVertical = true
	return c
}
func (c *CellBuilder) SetBorderTopBottom() *CellBuilder {
	c.borderBottom = true
	c.borderTop = true
	return c
}
func (c *CellBuilder) SetAlignment(align wml.ST_Jc) *CellBuilder {
	c.alignment = align
	return c
}
func (c *CellBuilder) Build() {
	if c.cell == nil {
		logrus.Warn("failed to build the cell, missing cell instance")
		return
	}

	bSty := wml.ST_BorderSingle
	bCol := color.Auto
	bThk := measurement.Distance(0.5 * measurement.Point)

	if c.borderBottom {
		c.cell.Properties().Borders().SetBottom(bSty, bCol, bThk)
	}
	if c.borderLeft {
		c.cell.Properties().Borders().SetLeft(bSty, bCol, bThk)
	}
	if c.borderRight {
		c.cell.Properties().Borders().SetRight(bSty, bCol, bThk)
	}
	if c.borderTop {
		c.cell.Properties().Borders().SetTop(bSty, bCol, bThk)
	}
	if c.borderInsideHorizontal {
		c.cell.Properties().Borders().SetInsideHorizontal(bSty, bCol, bThk)
	}
	if c.borderInsideVertical {
		c.cell.Properties().Borders().SetInsideVertical(bSty, bCol, bThk)
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
				run.Properties().SetFontFamily(xconditions.IfThenElse(c.fontFamily != "", c.fontFamily, cfg.FontFamily).(string))
				run.Properties().SetSize(measurement.Distance(xconditions.IfThenElse(c.fontSize > 0, c.fontSize, cfg.FontSize).(int)))
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
}
