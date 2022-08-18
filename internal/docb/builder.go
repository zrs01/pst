package docb

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
	"github.com/rotisserie/eris"
	"github.com/shomali11/util/xstrings"
	"github.com/thoas/go-funk"
	"github.com/zrs01/pst/internal/config"
	"gopkg.in/yaml.v2"
)

type Builder struct {
	cfile  string // config file name
	ifile  string // input file name
	ofile  string // output file name
	dfile  string // .docx file name
	config *config.Config
}

func Build(cfile, ifile, ofile string, tfile string) error {
	ncfg, err := config.NewConfig(cfile)
	if err != nil {
		return eris.Wrapf(err, "failed to load the configuration file %s", cfile)
	}

	// ifilePath = path.Dir(ifile)
	b := &Builder{ifile: ifile, ofile: ofile, dfile: tfile, config: ncfg}
	return b.construct()
}

func (b *Builder) construct() error {
	docb, err := NewDocumentBuilder(b.dfile, Configuration{
		FontFamily: b.config.FontFamily,
		FontSize:   b.config.FontSize,
		ImagePath:  filepath.Dir(b.ifile),
	})
	if err != nil {
		return eris.Wrap(err, "failed to create document builder")
	}

	// resolve wildcard
	files, err := b.resolveInputFile(b.ifile)
	if err != nil {
		return eris.Wrap(err, "failed to resolve the source file")
	}
	for _, file := range *files {
		data, err := b.loadData(file)
		if err != nil {
			return eris.Wrap(err, "failed to load the file")
		}

		// header
		docb.AddParagraph(func(p *ParagraphBuilder) {
			p.SetStyle("Heading1").SetText("PROGRAM DESCRIPTON")
		})

		for _, module := range data.Modules {
			docb.AddParagraph(func(p *ParagraphBuilder) {
				p.SetStyle("Heading2").SetText(module.Name)
			})
			for _, feature := range module.Features {
				docb.AddParagraph().
					AddParagraph(func(p *ParagraphBuilder) {
						p.SetStyle("Heading3").SetText(feature.Name)
					})
				if err := b.constructFeature(docb, &feature); err != nil {
					return eris.Wrap(err, "failed to build details")
				}
			}
			docb.AddParagraph(func(p *ParagraphBuilder) {
				p.SetPageBreak()
			})
		}
		docb.Build()

	}
	docb.Document.SaveToFile(b.ofile)
	return nil
}

func (b *Builder) constructFeature(docb *DocumentBuilder, feature *Feature) error {
	c1 := color.FromHex("ced4da") // gray
	c2 := color.FromHex("e9ecef") // light gray
	wd := float64(20)
	nd := docb.Document.Numbering.Definitions()[0]
	bs := wml.ST_BorderSingle
	bc := color.Auto
	bt := measurement.Distance(0.5 * measurement.Point)

	// generic table
	type xcol struct {
		value        interface{}
		bold         bool
		colspan      int
		widthPercent float64
		numbering    document.NumberingDefinition
		alignment    wml.ST_Jc
		allowEmpty   bool
	}
	type xrow struct {
		cols     []xcol
		bgColor  color.Color
		hasValue bool
	}
	var createTable = func(newParagraph bool, rf func() []xrow) {
		rows := rf()
		isContentBlank := true
		for i := 0; i < len(rows); i++ {
			for j := 0; j < len(rows[i].cols); j++ {
				if !b.isValueBlank(rows[i].cols[j].value) {
					isContentBlank = false
					break
				}
			}
		}
		if !isContentBlank {
			if newParagraph {
				docb.AddParagraph()
			}
			docb.AddTable(func(tb *TableBuilder) {
				tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) })
				for _, row := range rows {
					// check if row contains value
					isRowHasSomeValue := true
					isRowHasValue := true
					{
						count := 0
						for i := 0; i < len(row.cols); i++ {
							if !b.isValueBlank(row.cols[i].value) {
								count++
							}
						}
						isRowHasSomeValue = count > 0
						isRowHasValue = count == len(row.cols)
					}

					if (row.hasValue && isRowHasValue) || (!row.hasValue && isRowHasSomeValue) {
						tb.AddRow(func(rb *RowBuilder) {
							for _, col := range row.cols {
								if !b.isValueBlank(col.value) || col.allowEmpty {
									rb.AddCell(func(cb *CellBuilder) {
										cb.SetText(col.value)
										if (row.bgColor != color.Color{}) {
											cb.SetBackgroundColor(row.bgColor)
										}
										if col.bold {
											cb.SetBold()
										}
										if col.colspan > 0 {
											cb.SetColspan(col.colspan)
										}
										if col.widthPercent > 0 {
											cb.SetWidthPercent(col.widthPercent)
										}
										if (col.numbering != document.NumberingDefinition{}) {
											cb.SetBullet(&col.numbering)
										}
										if col.alignment != wml.ST_JcUnset {
											cb.SetAlignment(col.alignment)
										}
									})
								}
							}
						})

					}
				}
			})
		}
	}

	/* --------------------------------- PROGRAM -------------------------------- */
	createTable(false, func() []xrow {
		return []xrow{
			{cols: []xcol{{value: "Program ID", bold: true, widthPercent: wd}, {value: feature.Id}}, bgColor: c1},
			{cols: []xcol{{value: "Mode", bold: true}, {value: feature.Mode}}, hasValue: true},
			{cols: []xcol{{value: "Program Name", bold: true}, {value: feature.Name}}, hasValue: true},
			{cols: []xcol{{value: "Description", bold: true}, {value: feature.Desc}}, hasValue: true},

			{cols: []xcol{{value: "Program Environment", bold: true, colspan: 2}}, bgColor: c1},
			{cols: []xcol{{value: "Program Source", bold: true}, {value: feature.Mode}}, hasValue: true},
			{cols: []xcol{{value: "Language", bold: true}, {value: feature.Env.Languages}}, hasValue: true},
		}
	})

	/* -------------------------------- RESOURCE -------------------------------- */
	createTable(true, func() []xrow {
		var content []xrow
		if len(feature.Resources) > 0 {
			content = append(content,
				xrow{cols: []xcol{{value: "Resource", bold: true, colspan: 2}}, bgColor: c1},
				xrow{cols: []xcol{{value: "Table/File", bold: true}, {value: "Usage", bold: true}}, bgColor: c2},
			)
		}
		for _, res := range feature.Resources {
			content = append(content, xrow{cols: []xcol{{value: res.Name}, {value: res.Usage}}})
		}
		return content
	})

	/* --------------------------------- SCREEN --------------------------------- */
	if len(feature.Screens) > 0 {
		docb.AddParagraph().AddTable(func(tb *TableBuilder) {
			tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) }).
				AddRow(func(rb *RowBuilder) {
					rb.AddCell(func(cb *CellBuilder) {
						cb.SetText("Screen:").SetBold().SetColspan(2).SetBackgroundColor(c1)
					})
				})
			for _, scr := range feature.Screens {
				docb.AddTable(func(tb *TableBuilder) {
					tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) }).
						AddRow(func(rb *RowBuilder) {
							rb.
								AddCell(func(cb *CellBuilder) { cb.SetText("Screen ID").SetBold().SetWidthPercent(wd).SetBackgroundColor(c2) }).
								AddCell(func(cb *CellBuilder) { cb.SetText("Name").SetBold().SetBackgroundColor(c2) })
						})
					tb.AddRow(func(rb *RowBuilder) {
						rb.
							AddCell(func(cb *CellBuilder) { cb.SetText(scr.Id).SetWidthPercent(wd) }).
							AddCell(func(cb *CellBuilder) { cb.SetText(scr.Name) })
					})
					if xstrings.IsNotBlank(scr.Image.File) {
						tb.AddRow(func(rb *RowBuilder) {
							rb.AddCell(func(cb *CellBuilder) {
								cb.SetColspan(2).
									AddParagraph().AddParagraph(func(pb *ParagraphBuilder) {
									pb.SetAlignment(wml.ST_JcCenter).AddImage(func(ip *ImageProperty) { ip.SetFile(scr.Image.File).SetWidth(float64(scr.Image.Width)) })
								})
							})
						})
					}

				})
			}
		})
	}

	/* ---------------------------------- INPUT --------------------------------- */
	createTable(true, func() []xrow {
		var content []xrow
		if len(feature.Input) > 0 {
			content = append(content, xrow{cols: []xcol{{value: "Input", bold: true, colspan: 2, widthPercent: wd}}, bgColor: c1})
		}
		for i, input := range feature.Input {
			content = append(content,
				xrow{cols: []xcol{{value: fmt.Sprintf("%d. %s", i+1, input.Name), bold: true, colspan: 2}}, bgColor: c2},
				xrow{cols: []xcol{{value: "Fields", bold: true}, {value: input.Fields}}, hasValue: true},
				xrow{cols: []xcol{{value: "Constraints", bold: true}, {value: input.Constraints}}, hasValue: true},
				xrow{cols: []xcol{{value: "Remarks", bold: true}, {value: input.Remarks}}, hasValue: true},
			)
		}
		return content
	})

	/* ------------------------------- PARAMETERS ------------------------------- */
	createTable(true, func() []xrow {
		var content []xrow
		if len(feature.Parameters) > 0 {
			content = append(content,
				xrow{cols: []xcol{{value: "Parameters", bold: true, colspan: 5}}, bgColor: c1},
				xrow{cols: []xcol{
					{value: "ID", bold: true},
					{value: "Fields", bold: true}, {value: "Data Items"}, {value: "I/O", bold: true, alignment: wml.ST_JcCenter}, {value: "Processing Remarks", bold: true},
				}, bgColor: c2},
			)
			for i, param := range feature.Parameters {
				content = append(content, xrow{cols: []xcol{
					{value: fmt.Sprintf("%d", i+1)},
					{value: param.Field, allowEmpty: true},
					{value: param.Data, allowEmpty: true},
					{value: param.IO, allowEmpty: true},
					{value: param.Remarks, allowEmpty: true},
				}})
			}
		}
		return content
	})

	/* -------------------------------- SCENARIO -------------------------------- */
	createTable(true, func() []xrow {
		var content []xrow
		if len(feature.Scenarios) > 0 {
			content = append(content, xrow{cols: []xcol{{value: "Scenarios and Processign Logic", bold: true, colspan: 2}}, bgColor: c1})
		}
		for i, scn := range feature.Scenarios {
			content = append(content, xrow{cols: []xcol{{value: fmt.Sprintf("%d. %s", i+1, scn.Name), bold: true, colspan: 2}}, bgColor: c2})
			for _, action := range scn.Desc {
				keyword, others := b.splitGherkinWord(action)
				content = append(content, xrow{cols: []xcol{{value: keyword, bold: true, widthPercent: 10}, {value: others}}})
			}
		}
		return content
	})

	/* --------------------------------- OTHERS --------------------------------- */
	createTable(true, func() []xrow {
		var content []xrow

		textMap := map[string]string{"Reference": "External Reference", "Limits": "Program Limits", "Remarks": "Remarks"}
		flds := reflect.VisibleFields(reflect.TypeOf(feature.Others))
		for _, fld := range flds {
			for _, index := range fld.Index {
				value := reflect.ValueOf(feature.Others).Field(index)
				if !value.IsNil() {
					content = append(content,
						xrow{cols: []xcol{{value: textMap[fld.Name]}}, bgColor: c1},
						xrow{cols: []xcol{{value: value.Interface(), numbering: nd}}},
					)
				}
			}
		}
		return content
	})

	return nil
}

func (b *Builder) loadData(file string) (*ProgSpec, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, eris.Wrap(err, "failed to read the file")
	}
	var d ProgSpec
	if err := yaml.Unmarshal(yamlFile, &d); err != nil {
		return nil, eris.Wrapf(err, "failed to unmarshal the file %s", file)
	}
	return &d, nil
}

func (b *Builder) resolveInputFile(ifile string) (*[]string, error) {
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

func (b *Builder) splitGherkinWord(s string) (string, string) {
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

func (b *Builder) isValueBlank(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsZero()
}

func (b *Builder) isAllFieldsInStructEmpty(value interface{}) bool {
	flds := reflect.VisibleFields(reflect.TypeOf(value))
	for _, fld := range flds {
		for _, index := range fld.Index {
			value := reflect.ValueOf(value).Field(index)
			if !value.IsNil() {
				return false
			}
		}
	}
	return true
}
