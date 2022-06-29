package converter

// https://github.com/bollwarm/gooxml
import (
	"fmt"
	"io/ioutil"

	"baliance.com/gooxml"
	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/ofc/sharedTypes"
	"baliance.com/gooxml/schema/soo/wml"
	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xconditions"
	"gopkg.in/yaml.v2"
)

type builder struct {
	cfg *Config
}

func Build(cfile, ifile, ofile string) error {
	b := &builder{}
	spec, err := b.loadData(ifile)
	if err != nil {
		return eris.Wrapf(err, "failed to load %s", ifile)
	}
	b.buildSpec(spec, cfile, ofile)
	return nil
}

func (b *builder) loadData(infile string) (*ProgSpec, error) {
	yamlFile, err := ioutil.ReadFile(infile)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to read the file %s", infile)
	}
	var d ProgSpec
	if err := yaml.Unmarshal(yamlFile, &d); err != nil {
		return nil, eris.Wrapf(err, "failed to unmarshal the file %s", infile)
	}
	return &d, nil
}

func (b *builder) buildSpec(data *ProgSpec, cfile, ofile string) error {
	cfg, err := NewConfig(cfile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the configuration file %s", cfile)
	}
	b.cfg = cfg

	doc := document.New()
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

	nd := doc.Numbering.AddDefinition()
	nd = doc.Numbering.AddDefinition()
	// nd := doc.Numbering.Definitions()[0]
	for _, module := range data.Modules {
		for _, feature := range module.Features {
			p := doc.AddParagraph()
			p.SetNumberingDefinition(nd)
			p.SetNumberingLevel(0)
			p.AddRun().AddText(feature.Name)
			b.buildFeature(doc, feature)
		}
	}
	doc.SaveToFile(ofile)
	return nil
}

func (b *builder) buildFeature(doc *document.Document, feature Feature) error {
	// var setRunProperty = func(run *document.Run) {
	// 	run.Properties().SetFontFamily(b.cfg.FontFamily)
	// 	run.Properties().SetSize(measurement.Distance(b.cfg.FontSize))
	// }
	var createTable = func(p bool) document.Table {
		if p {
			doc.AddParagraph()
		}
		table := doc.AddTable()
		table.Properties().SetWidthPercent(100)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)
		return table
	}
	// var createVerticalBullet = func(cell *document.Cell, title string, content interface{}) {
	// 	text := toStrArray(content)
	// 	if len(text) > 0 {
	// 		{
	// 			par := cell.AddParagraph()
	// 			run := par.AddRun()
	// 			setRunProperty(&run)
	// 			run.Properties().SetUnderline(wml.ST_UnderlineSingle, color.Black)
	// 			run.AddText(title)
	// 		}
	// 		{
	// 			nd := doc.Numbering.Definitions()[0]
	// 			for _, t := range text {
	// 				par := cell.AddParagraph()
	// 				par.SetNumberingLevel(0)
	// 				par.SetNumberingDefinition(nd)
	// 				run := par.AddRun()
	// 				setRunProperty(&run)
	// 				run.AddText("\t" + t)
	// 			}
	// 		}
	// 	}
	// }

	/* ---------------------------------- main ---------------------------------- */

	w := 25
	gray := color.FromHex("ced4da")
	lightgray := color.FromHex("e9ecef")
	nd := doc.Numbering.Definitions()[0]
	rb := rowBuilder{cfg: b.cfg}

	tbMain := createTable(false)
	rb.Reset(tbMain.AddRow()).AddCell((&cellAttr{}).SetText("Program ID").SetWidth(w).SetBold(true), (&cellAttr{}).SetText(feature.Id)).Build()
	rb.Reset(tbMain.AddRow()).AddCell((&cellAttr{}).SetText("Mode").SetWidth(w).SetBold(true), (&cellAttr{}).SetText(feature.Mode)).Build()
	rb.Reset(tbMain.AddRow()).AddCell((&cellAttr{}).SetText("Program Name").SetWidth(w).SetBold(true), (&cellAttr{}).SetText(feature.Name)).Build()
	rb.Reset(tbMain.AddRow()).AddCell((&cellAttr{}).SetText("Description").SetWidth(w).SetBold(true), (&cellAttr{}).SetText(feature.Desc)).Build()

	rb.Reset(tbMain.AddRow()).AddCell((&cellAttr{}).SetText("Program Environment:").SetBold(true).setColspan(2).setBackgroundColor(gray)).Build()
	rb.Reset(tbMain.AddRow()).AddCell((&cellAttr{}).SetText("Program Source").SetWidth(w).SetBold(true), (&cellAttr{}).SetText(feature.Env.Sources)).Build()
	rb.Reset(tbMain.AddRow()).AddCell((&cellAttr{}).SetText("Language").SetWidth(w).SetBold(true), (&cellAttr{}).SetText(feature.Env.Languages)).Build()

	if len(feature.Resources) > 0 {
		tbRes := createTable(true)
		rb.Reset(tbRes.AddRow()).AddCell((&cellAttr{}).SetText("Resources:").SetBold(true).setColspan(2).setBackgroundColor(gray)).Build()
		rb.Reset(tbRes.AddRow()).AddCell((&cellAttr{}).SetText("Table/File").SetWidth(w).SetBold(true).setBackgroundColor(lightgray), (&cellAttr{}).SetText("Usage").SetBold(true).setBackgroundColor(lightgray)).Build()
		for _, res := range feature.Resources {
			rb.Reset(tbRes.AddRow()).AddCell((&cellAttr{}).SetText(res.Name).SetWidth(w), (&cellAttr{}).SetText(res.Usage)).Build()
		}

		tbInp := createTable(true)
		rb.Reset(tbInp.AddRow()).AddCell((&cellAttr{}).SetText("Input:").SetBold(true).setColspan(4).setBackgroundColor(gray)).Build()
		for i, inp := range feature.Input {
			rb.Reset(tbInp.AddRow()).AddCell(
				(&cellAttr{}).SetText(fmt.Sprintf("%d. %s", i+1, inp.Name)).SetBold(true).setBackgroundColor(lightgray),
			).Build()
			if len(inp.Fields) > 0 {
				rb.Reset(tbInp.AddRow()).AddCell(
					(&cellAttr{}).SetText("Fields").SetBold(true),
					(&cellAttr{}).SetText(inp.Fields).setColspan(3).setBullet(&nd),
				).Build()
			}
			if len(inp.Constraints) > 0 {
				rb.Reset(tbInp.AddRow()).AddCell(
					(&cellAttr{}).SetText("Constraints").SetBold(true),
					(&cellAttr{}).SetText(inp.Constraints).setColspan(3).setBullet(&nd),
				).Build()
			}
			if len(inp.Remarks) > 0 {
				rb.Reset(tbInp.AddRow()).AddCell(
					(&cellAttr{}).SetText("Remarks").SetBold(true),
					(&cellAttr{}).SetText(inp.Remarks).setColspan(3).setBullet(&nd),
				).Build()
			}
		}

		tbScn := createTable(true)
		rb.Reset(tbScn.AddRow()).AddCell((&cellAttr{}).SetText("Scenarios:").SetBold(true).setColspan(4).setBackgroundColor(gray)).Build()
		for i, scn := range feature.Scenarios {
			rb.Reset(tbScn.AddRow()).AddCell(
				(&cellAttr{}).SetText(fmt.Sprintf("%d. %s", i+1, scn.Name)).SetBold(true).setBackgroundColor(lightgray).setColspan(2),
			).Build()
			if len(scn.Given) > 0 {
				rb.Reset(tbScn.AddRow()).AddCell(
					(&cellAttr{}).SetText("Given").SetBold(true),
					(&cellAttr{}).SetText(scn.Given).setBullet(&nd),
				).Build()
			}
			if len(scn.When) > 0 {
				rb.Reset(tbScn.AddRow()).AddCell(
					(&cellAttr{}).SetText("When").SetBold(true),
					(&cellAttr{}).SetText(scn.When).setBullet(&nd),
				).Build()
			}
			if len(scn.And) > 0 {
				rb.Reset(tbScn.AddRow()).AddCell(
					(&cellAttr{}).SetText("And").SetBold(true),
					(&cellAttr{}).SetText(scn.And).setBullet(&nd),
				).Build()
			}
			if len(scn.But) > 0 {
				rb.Reset(tbScn.AddRow()).AddCell(
					(&cellAttr{}).SetText("But").SetBold(true),
					(&cellAttr{}).SetText(scn.But).setBullet(&nd),
				).Build()
			}
			if len(scn.Then) > 0 {
				rb.Reset(tbScn.AddRow()).AddCell(
					(&cellAttr{}).SetText("Then").SetBold(true),
					(&cellAttr{}).SetText(scn.Then).setBullet(&nd),
				).Build()
			}
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
func (c *rowBuilder) BuildCustom(setContent func(cell *document.Cell)) {
	for _, attr := range c.cellAttrs {
		nc := c.row.AddCell()
		// nc.Properties().Borders().SetBottom(wml.ST_BorderDouble, color.Black, 0.5*measurement.Point)
		if attr.width > 0 {
			nc.Properties().SetWidthPercent(float64(attr.width))
		}
		if attr.colspan > 0 {
			nc.Properties().SetColumnSpan(attr.colspan)
		}
		if attr.backgroundColor != nil {
			nc.Properties().SetShading(wml.ST_ShdSolid, *attr.backgroundColor, color.Auto)
		}
		if setContent != nil {
			setContent(&nc)
		} else {
			for _, t := range attr.text {
				p := nc.AddParagraph()
				tab := ""
				if attr.bullet != nil {
					p.SetNumberingLevel(0)
					p.SetNumberingDefinition(*attr.bullet)
					tab = "\t"
				}
				run := p.AddRun()
				run.Properties().SetBold(attr.isBold)
				run.Properties().SetFontFamily(xconditions.IfThenElse(attr.fontFamily != "", attr.fontFamily, c.cfg.FontFamily).(string))
				run.Properties().SetSize(measurement.Distance(xconditions.IfThenElse(attr.fontSize > 0, attr.fontSize, c.cfg.FontSize).(int)))
				run.AddText(tab + t)
			}
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                                CELL BUILDER                                */
/* -------------------------------------------------------------------------- */

type cellAttr struct {
	fontFamily      string
	fontSize        int
	isBold          bool
	width           int
	colspan         int
	bullet          *document.NumberingDefinition
	text            []string
	backgroundColor *color.Color
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
	c.isBold = b
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
