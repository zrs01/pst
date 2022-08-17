// document constructor
package docb

import (
	"strings"

	"baliance.com/gooxml"
	"baliance.com/gooxml/document"
	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xstrings"
)

// Reference:
// https://github.com/bollwarm/gooxml   Office Open XML documents utility
// https://github.com/unidoc/unioffice-examples

type builder interface {
	Build()
}

type Configuration struct {
	FontFamily string
	FontSize   int
	ImagePath  string
}

/* -------------------------------------------------------------------------- */
/*                                  Document                                  */
/* -------------------------------------------------------------------------- */
type DocumentBuilder struct {
	Document *document.Document // allow outsider to custom document
	builder  []builder
	config   *Configuration
}

func NewDocumentBuilder(file string, cfg ...Configuration) (*DocumentBuilder, error) {
	// override the options if is is provided
	c := &Configuration{
		FontFamily: "Arial",
		FontSize:   10,
	}
	if len(cfg) > 0 {
		if xstrings.IsNotBlank(cfg[0].FontFamily) {
			c.FontFamily = cfg[0].FontFamily
		}
		if cfg[0].FontSize > 0 {
			c.FontSize = cfg[0].FontSize
		}
		c.ImagePath = cfg[0].ImagePath
	}

	// create document object at once for outsider customization
	var doc *document.Document
	if xstrings.IsBlank(file) {
		doc = document.New()
	} else {
		d, err := document.Open(file)
		if err != nil {
			return nil, eris.Wrapf(err, "failed to open file %s", file)
		}
		doc = d
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
	}
	// adjust the default settings
	fixBulletIndentation(doc)
	return &DocumentBuilder{Document: doc, config: c}, nil
}

func (d *DocumentBuilder) AddParagraph(nextBuilders ...func(*ParagraphBuilder)) *DocumentBuilder {
	p := newParagraphBuilder(d.config, d.Document, d.Document.AddParagraph())
	d.builder = append(d.builder, p)
	for _, nextBuilder := range nextBuilders {
		nextBuilder(p)
	}
	return d
}

func (d *DocumentBuilder) AddTable(nextBuilder func(*TableBuilder)) *DocumentBuilder {
	t := newTableBuilder(d.config, d.Document, d.Document.AddTable())
	d.builder = append(d.builder, t)
	nextBuilder(t)
	return d
}

func (d *DocumentBuilder) Build() {
	for _, builder := range d.builder {
		builder.Build()
	}
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
