package converter

import (
	"io/ioutil"

	"baliance.com/gooxml"
	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/ofc/sharedTypes"
	"baliance.com/gooxml/schema/soo/wml"
	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

func Build(ifile, ofile string) error {
	spec, err := loadData(ifile)
	if err != nil {
		return eris.Wrapf(err, "failed to load %s", ifile)
	}
	buildSpec(spec, ofile)
	return nil
}

func loadData(infile string) (*ProgSpec, error) {
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

func buildSpec(data *ProgSpec, ofile string) error {
	doc := document.New()
	// set to A4 size (https://stackoverflow.com/questions/57581695/detecting-and-setting-paper-size-in-word-js-api-or-ooxml)
	doc.X().Body = wml.NewCT_Body()
	doc.X().Body.SectPr = wml.NewCT_SectPr()
	doc.X().Body.SectPr.PgSz = &wml.CT_PageSz{
		WAttr: &sharedTypes.ST_TwipsMeasure{
			ST_UnsignedDecimalNumber: gooxml.Uint64(uint64(11906)),
		},
		HAttr: &sharedTypes.ST_TwipsMeasure{
			ST_UnsignedDecimalNumber: gooxml.Uint64(uint64(16838)),
		},
	}

	for _, module := range data.Modules {
		for _, feature := range module.Features {
			buildFeature(doc, feature)
		}
	}
	doc.SaveToFile(ofile)
	return nil
}

func buildFeature(doc *document.Document, feature Feature) error {
	var addLabelField = func(row document.Row, label, field string) {
		cell := row.AddCell()
		run := cell.AddParagraph().AddRun()
		run.Properties().SetBold(true)
		run.AddText(label)
		row.AddCell().AddParagraph().AddRun().AddText(field)
	}

	/* ---------------------------------- main ---------------------------------- */

	table := doc.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

	// program ID & mode
	addLabelField(table.AddRow(), "Program ID", feature.Id)
	addLabelField(table.AddRow(), "Mode", feature.Mode)

	return nil
}
