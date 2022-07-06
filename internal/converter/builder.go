package converter

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"baliance.com/gooxml"
	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/ofc/sharedTypes"
	"baliance.com/gooxml/schema/soo/wml"
	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xconditions"
	"github.com/shomali11/util/xstrings"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

// Reference:
// https://github.com/bollwarm/gooxml   Office Open XML documents utility
// https://github.com/shomali11/util		A group of generic useful utility functions

type builder struct {
	cfg   *Config
	cfile string // config file name
	ifile string // input file name
	ofile string // output file name
	dfile string // .docx file name
}

func Build(cfile, ifile, ofile string, tfile string) error {
	b := &builder{cfile: cfile, ifile: ifile, ofile: ofile, dfile: tfile}
	return b.buildSpec()
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

func (b *builder) resolveInFile(ifile string) (*[]string, error) {
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

func (b *builder) buildSpec() error {
	cfg, err := NewConfig(b.cfile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the configuration file %s", b.cfile)
	}
	b.cfg = cfg

	var doc *document.Document
	if xstrings.IsBlank(b.dfile) {
		doc = document.New()
		// set to A4 size (https://stackoverflow.com/questions/57581695/detecting-and-setting-paper-size-in-word-js-api-or-ooxml)
		doc.X().Body.SectPr = wml.NewCT_SectPr()
		doc.X().Body.SectPr.PgSz = &wml.CT_PageSz{
			WAttr: &sharedTypes.ST_TwipsMeasure{
				ST_UnsignedDecimalNumber: gooxml.Uint64(uint64(11906)),
			},
			HAttr: &sharedTypes.ST_TwipsMeasure{
				ST_UnsignedDecimalNumber: gooxml.Uint64(uint64(16838)),
			},
		}
	} else {
		if doc, err = document.Open(b.dfile); err != nil {
			return eris.Wrapf(err, "failed to open the document %s", b.dfile)
		}
	}
	fixBulletIndentation(doc)

	// support multiple files
	files, err := b.resolveInFile(b.ifile)
	if err != nil {
		return eris.Wrap(err, "failed to resolve the source file")
	}
	for _, file := range *files {
		data, err := b.loadData(file)
		if err != nil {
			return eris.Wrap(err, "failed to load the file")
		}

		// header
		newParaBuilder(doc.AddParagraph()).SetStyle("Heading1").SetText("PROGRAM DESCRIPTION").Build()

		for _, module := range data.Modules {
			newParaBuilder(doc.AddParagraph()).SetStyle("Heading2").SetText(module.Name).Build()
			for _, feature := range module.Features {
				doc.AddParagraph()
				newParaBuilder(doc.AddParagraph()).SetStyle("Heading3").SetText(feature.Name).Build()
				b.buildFeature(doc, &feature)
			}
			// if i < len(data.Modules)-1 {
			newParaBuilder(doc.AddParagraph()).SetPageBreak().Build()
			// }
		}

	}

	doc.SaveToFile(b.ofile)
	return nil
}

func (b *builder) buildFeature(doc *document.Document, feature *Feature) error {
	var createTable = func(p bool) document.Table {
		if p {
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
	rb := rowBuilder{cfg: b.cfg}

	tbMain := createTable(false)
	rb.Reset(tbMain.AddRow()).AddCell(
		newCellAttr().SetText("Program ID").SetWidth(w).SetBold(true).setBorderTopBottom(),
		newCellAttr().SetText(feature.Id).setBorderTopBottom()).Build()
	rb.Reset(tbMain.AddRow()).AddCell(
		newCellAttr().SetText("Mode").SetWidth(w).SetBold(true).setBorderTopBottom(),
		newCellAttr().SetText(feature.Mode).setBorderTopBottom()).Build()
	rb.Reset(tbMain.AddRow()).AddCell(
		newCellAttr().SetText("Program Name").SetWidth(w).SetBold(true).setBorderTopBottom(),
		newCellAttr().SetText(feature.Name).setBorderTopBottom()).Build()
	rb.Reset(tbMain.AddRow()).AddCell(
		newCellAttr().SetText("Description").SetWidth(w).SetBold(true).setBorderTopBottom(),
		newCellAttr().SetText(feature.Desc).setBorderTopBottom()).Build()

	if len(feature.Env.Sources) > 0 || len(feature.Env.Languages) > 0 {
		rb.Reset(tbMain.AddRow()).AddCell(
			newCellAttr().SetText("Program Environment:").SetBold(true).setColspan(2).setBackgroundColor(gray).setBorderTopBottom()).Build()
		rb.Reset(tbMain.AddRow()).AddCell(
			newCellAttr().SetText("Program Source").SetWidth(w).SetBold(true).setBorderTopBottom(),
			newCellAttr().SetText(feature.Env.Sources).setBorderTopBottom()).Build()
		rb.Reset(tbMain.AddRow()).AddCell(
			newCellAttr().SetText("Language").SetWidth(w).SetBold(true).setBorderTopBottom(),
			newCellAttr().SetText(feature.Env.Languages).setBorderTopBottom()).Build()
	}

	if len(feature.Resources) > 0 {
		// Resources
		tbRes := createTable(true)
		rb.Reset(tbRes.AddRow()).AddCell(
			newCellAttr().SetText("Resources:").SetBold(true).setColspan(2).setBackgroundColor(gray).setBorderTopBottom()).Build()
		rb.Reset(tbRes.AddRow()).AddCell(
			newCellAttr().SetText("Table/File").SetWidth(w).SetBold(true).setBackgroundColor(lightgray).setBorderTopBottom(),
			newCellAttr().SetText("Usage").SetBold(true).setBackgroundColor(lightgray).setBorderTopBottom()).Build()
		for _, res := range feature.Resources {
			rb.Reset(tbRes.AddRow()).AddCell(
				newCellAttr().SetText(res.Name).SetWidth(w).setBorderTopBottom(),
				newCellAttr().SetText(res.Usage).setBorderTopBottom()).Build()
		}

		// Input
		if len(feature.Input) > 0 {
			tbInp := createTable(true)
			rb.Reset(tbInp.AddRow()).AddCell(
				newCellAttr().SetText("Input:").SetBold(true).setColspan(2).setBackgroundColor(gray).setBorderTopBottom()).Build()
			for i, inp := range feature.Input {
				rb.Reset(tbInp.AddRow()).AddCell(
					newCellAttr().SetText(fmt.Sprintf("%d. %s", i+1, inp.Name)).SetBold(true).setColspan(2).setBackgroundColor(lightgray).setBorderTopBottom(),
				).Build()
				if len(inp.Fields) > 0 {
					rb.Reset(tbInp.AddRow()).AddCell(
						newCellAttr().SetText("Fields").SetBold(true).SetWidth(w),
						newCellAttr().SetText(inp.Fields),
					).Build()
				}
				if len(inp.Constraints) > 0 {
					rb.Reset(tbInp.AddRow()).AddCell(
						newCellAttr().SetText("Constraints").SetBold(true).SetWidth(w),
						newCellAttr().SetText(inp.Constraints).setBullet(&nd),
					).Build()
				}
				if len(inp.Remarks) > 0 {
					rb.Reset(tbInp.AddRow()).AddCell(
						newCellAttr().SetText("Remarks").SetBold(true).SetWidth(w),
						newCellAttr().SetText(inp.Remarks).setBullet(&nd),
					).Build()
				}
			}
		}

		// Scenarios
		if len(feature.Scenarios) > 0 {
			tbScn := createTable(true)
			rb.Reset(tbScn.AddRow()).AddCell(
				newCellAttr().SetText("Scenarios:").SetBold(true).setColspan(2).setBackgroundColor(gray).setBorderTopBottom()).Build()
			for i, scn := range feature.Scenarios {
				rb.Reset(tbScn.AddRow()).AddCell(
					newCellAttr().SetText(fmt.Sprintf("%d. %s", i+1, scn.Name)).SetBold(true).setBackgroundColor(lightgray).setColspan(2).setBorderTopBottom(),
				).Build()

				for _, action := range scn.Desc {
					keyword, others := splitScenarioWord(action)
					rb.Reset(tbScn.AddRow()).AddCell(
						newCellAttr().SetText(keyword).SetBold(true).SetAlignment(wml.ST_JcLeft).SetWidth(10),
						newCellAttr().SetText(others),
					).Build()
				}
			}
		}

		// Remarks
		if len(feature.Remarks) > 0 {
			tbRmk := createTable(true)
			rb.Reset(tbRmk.AddRow()).AddCell(
				newCellAttr().SetText("Remarks:").SetBold(true).setBackgroundColor(gray).setBorderTopBottom()).Build()
			rb.Reset(tbRmk.AddRow()).AddCell(
				newCellAttr().SetText(feature.Remarks).setBullet(&nd),
			).Build()
		}

	}
	return nil
}

/* -------------------------------- UTILITIES ------------------------------- */
func toStrArray(v interface{}) []string {
	text := []string{}
	switch t := v.(type) {
	case string:
		text = []string{t}
	case []string:
		text = t
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

type paraBuilder struct {
	para      *document.Paragraph
	style     string
	text      string
	pageBreak bool
	lineBreak bool
}

func newParaBuilder(p document.Paragraph) *paraBuilder {
	return &paraBuilder{para: &p}
}
func (p *paraBuilder) SetStyle(s string) *paraBuilder {
	p.style = s
	return p
}
func (p *paraBuilder) SetText(s string) *paraBuilder {
	p.text = s
	return p
}
func (p *paraBuilder) SetLineBreak() *paraBuilder {
	p.lineBreak = true
	return p
}
func (p *paraBuilder) SetPageBreak() *paraBuilder {
	p.pageBreak = true
	return p
}
func (p *paraBuilder) Build() {
	if xstrings.IsNotBlank(p.style) {
		p.para.SetStyle(p.style)
	}
	run := p.para.AddRun()
	if xstrings.IsNotBlank(p.text) {
		run.AddText(p.text)
	}
	if p.lineBreak {
		run.AddBreak()
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

type rowBuilder struct {
	cfg       *Config
	row       document.Row
	cellAttrs []*cellAttr
}

func (r *rowBuilder) Reset(row document.Row) *rowBuilder {
	r.row = row
	r.cellAttrs = []*cellAttr{}
	return r
}

func (r *rowBuilder) AddCell(attrs ...*cellAttr) *rowBuilder {
	r.cellAttrs = append(r.cellAttrs, attrs...)
	return r
}

func (c *rowBuilder) Build() {
	c.BuildCustom(nil)
}
func (c *rowBuilder) BuildCustom(customContent func(cell *document.Cell)) {
	bSty := wml.ST_BorderSingle
	bCol := color.Auto
	bThk := measurement.Distance(0.5 * measurement.Point)

	for _, attr := range c.cellAttrs {
		nc := c.row.AddCell()

		// border
		if attr.borderBottom {
			nc.Properties().Borders().SetBottom(bSty, bCol, bThk)
		}
		if attr.borderLeft {
			nc.Properties().Borders().SetLeft(bSty, bCol, bThk)
		}
		if attr.borderRight {
			nc.Properties().Borders().SetRight(bSty, bCol, bThk)
		}
		if attr.borderTop {
			nc.Properties().Borders().SetTop(bSty, bCol, bThk)
		}
		if attr.borderInsideHorizontal {
			nc.Properties().Borders().SetInsideHorizontal(bSty, bCol, bThk)
		}
		if attr.borderInsideVertical {
			nc.Properties().Borders().SetInsideVertical(bSty, bCol, bThk)
		}
		// width
		if attr.width > 0 {
			nc.Properties().SetWidthPercent(float64(attr.width))
		}
		// colspan
		if attr.colspan > 0 {
			nc.Properties().SetColumnSpan(attr.colspan)
		}
		// background color
		if attr.backgroundColor != nil {
			nc.Properties().SetShading(wml.ST_ShdSolid, *attr.backgroundColor, color.Auto)
		}
		if customContent != nil {
			customContent(&nc)
		} else {
			if len(attr.text) == 0 {
				p := nc.AddParagraph()
				p.AddRun().AddText("")
			} else {
				for _, t := range attr.text {
					p := nc.AddParagraph()
					if attr.alignment != wml.ST_JcUnset {
						p.Properties().SetAlignment(attr.alignment)
					}
					tab := ""
					if attr.bullet != nil && len(attr.text) > 1 {
						p.SetNumberingLevel(0)
						p.SetNumberingDefinition(*attr.bullet)
						tab = "\t"
					}
					lines := strings.Split(t, "\\n")
					for i, line := range lines {
						line := strings.ReplaceAll(line, "\\t", "\t")
						run := p.AddRun()
						run.Properties().SetBold(attr.bold)
						run.Properties().SetFontFamily(xconditions.IfThenElse(attr.fontFamily != "", attr.fontFamily, c.cfg.FontFamily).(string))
						run.Properties().SetSize(measurement.Distance(xconditions.IfThenElse(attr.fontSize > 0, attr.fontSize, c.cfg.FontSize).(int)))
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
	}
}

/* -------------------------------------------------------------------------- */
/*                                CELL BUILDER                                */
/* -------------------------------------------------------------------------- */

type cellAttr struct {
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
}

func newCellAttr() *cellAttr {
	return &cellAttr{}
}
func (c *cellAttr) SetFontFamily(ff string) *cellAttr {
	c.fontFamily = ff
	return c
}
func (c *cellAttr) SetFontSize(fs int) *cellAttr {
	c.fontSize = fs
	return c
}
func (c *cellAttr) SetBold(b bool) *cellAttr {
	c.bold = b
	return c
}
func (c *cellAttr) SetWidth(w int) *cellAttr {
	c.width = w
	return c
}
func (c *cellAttr) setColspan(cols int) *cellAttr {
	c.colspan = cols
	return c
}
func (c *cellAttr) SetText(v interface{}) *cellAttr {
	c.text = toStrArray(v)
	return c
}
func (c *cellAttr) setBullet(b *document.NumberingDefinition) *cellAttr {
	c.bullet = b
	return c
}
func (c *cellAttr) setBackgroundColor(color color.Color) *cellAttr {
	c.backgroundColor = &color
	return c
}
func (c *cellAttr) setBorderBottom() *cellAttr {
	c.borderBottom = true
	return c
}
func (c *cellAttr) setBorderLeft() *cellAttr {
	c.borderLeft = true
	return c
}
func (c *cellAttr) setBorderRight() *cellAttr {
	c.borderRight = true
	return c
}
func (c *cellAttr) setBorderTop() *cellAttr {
	c.borderTop = true
	return c
}
func (c *cellAttr) setBorderInsideHorizontal() *cellAttr {
	c.borderInsideHorizontal = true
	return c
}
func (c *cellAttr) setBorderInsideVertical() *cellAttr {
	c.borderInsideVertical = true
	return c
}
func (c *cellAttr) setBorderAll() *cellAttr {
	c.borderBottom = true
	c.borderLeft = true
	c.borderRight = true
	c.borderTop = true
	c.borderInsideHorizontal = true
	c.borderInsideVertical = true
	return c
}
func (c *cellAttr) setBorderTopBottom() *cellAttr {
	c.borderBottom = true
	c.borderTop = true
	return c
}
func (c *cellAttr) SetAlignment(align wml.ST_Jc) *cellAttr {
	c.alignment = align
	return c
}
