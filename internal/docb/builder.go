package docb

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"baliance.com/gooxml/color"
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
				if err := b.buildFeature(docb, &feature); err != nil {
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

func (b *Builder) buildFeature(docb *DocumentBuilder, feature *Feature) error {
	c1 := color.FromHex("ced4da") // gray
	c2 := color.FromHex("e9ecef") // light gray
	wd := 20
	nd := docb.Document.Numbering.Definitions()[0]
	bs := wml.ST_BorderSingle
	bc := color.Auto
	bt := measurement.Distance(0.5 * measurement.Point)

	/* --------------------------------- PROGRAM -------------------------------- */
	docb.AddTable(func(tb *TableBuilder) {
		tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) }).
			AddRow(func(rb *RowBuilder) {
				rb.
					AddCell(func(cb *CellBuilder) { cb.SetText("Program ID").SetBold().SetWidth(wd).SetBackgroundColor(c1) }).
					AddCell(func(cb *CellBuilder) { cb.SetText(feature.Id).SetBackgroundColor(c1) })
			}).
			AddRow(func(rb *RowBuilder) {
				rb.
					AddCell(func(cb *CellBuilder) { cb.SetText("Mode").SetBold().SetWidth(wd) }).
					AddCell(func(cb *CellBuilder) { cb.SetText(feature.Mode) })
			}).
			AddRow(func(rb *RowBuilder) {
				rb.
					AddCell(func(cb *CellBuilder) { cb.SetText("Program Name").SetBold().SetWidth(wd) }).
					AddCell(func(cb *CellBuilder) { cb.SetText(feature.Name) })
			}).
			AddRow(func(rb *RowBuilder) {
				rb.
					AddCell(func(cb *CellBuilder) { cb.SetText("Description").SetBold().SetWidth(wd) }).
					AddCell(func(cb *CellBuilder) { cb.SetText(feature.Desc) })
			})
	})

	if (feature.Env.Sources != nil && !reflect.ValueOf(feature.Env.Sources).IsZero()) || (feature.Env.Languages != nil && !reflect.ValueOf(feature.Env.Languages).IsZero()) {
		docb.AddTable(func(tb *TableBuilder) {
			tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) }).
				AddRow(func(rb *RowBuilder) {
					rb.AddCell(func(cb *CellBuilder) {
						cb.SetText("Program Environment:").SetBold().SetColspan(2).SetBackgroundColor(c1)
					})
				})
			if feature.Env.Sources != nil && !reflect.ValueOf(feature.Env.Sources).IsZero() {
				tb.AddRow(func(rb *RowBuilder) {
					rb.
						AddCell(func(cb *CellBuilder) { cb.SetText("Program Source").SetBold().SetWidth(wd) }).
						AddCell(func(cb *CellBuilder) { cb.SetText(feature.Env.Sources) })
				})
			}
			if feature.Env.Languages != nil && !reflect.ValueOf(feature.Env.Languages).IsZero() {
				tb.AddRow(func(rb *RowBuilder) {
					rb.
						AddCell(func(cb *CellBuilder) { cb.SetText("Language").SetBold().SetWidth(wd) }).
						AddCell(func(cb *CellBuilder) { cb.SetText(feature.Env.Languages) })
				})
			}
		})
	}

	/* -------------------------------- RESOURCE -------------------------------- */
	if len(feature.Resources) > 0 {
		docb.AddParagraph().AddTable(func(tb *TableBuilder) {
			tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) }).
				AddRow(func(rb *RowBuilder) {
					rb.AddCell(func(cb *CellBuilder) {
						cb.SetText("Resource:").SetBold().SetColspan(2).SetBackgroundColor(c1)
					})
				})
			tb.AddRow(func(rb *RowBuilder) {
				rb.
					AddCell(func(cb *CellBuilder) { cb.SetText("Table/File").SetBold().SetWidth(wd).SetBackgroundColor(c2) }).
					AddCell(func(cb *CellBuilder) { cb.SetText("Usage").SetBold().SetBackgroundColor(c2) })
			})
			for _, res := range feature.Resources {
				tb.AddRow(func(rb *RowBuilder) {
					rb.
						AddCell(func(cb *CellBuilder) { cb.SetText(res.Name).SetWidth(wd) }).
						AddCell(func(cb *CellBuilder) { cb.SetText(res.Usage) })
				})
			}
		})
	}

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
								AddCell(func(cb *CellBuilder) { cb.SetText("Screen ID").SetBold().SetWidth(wd).SetBackgroundColor(c2) }).
								AddCell(func(cb *CellBuilder) { cb.SetText("Name").SetBold().SetBackgroundColor(c2) })
						})
					tb.AddRow(func(rb *RowBuilder) {
						rb.
							AddCell(func(cb *CellBuilder) { cb.SetText(scr.Id).SetWidth(wd) }).
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
	if len(feature.Input) > 0 {
		docb.AddParagraph().AddTable(func(tb *TableBuilder) {
			tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) }).
				AddRow(func(rb *RowBuilder) {
					rb.AddCell(func(cb *CellBuilder) {
						cb.SetText("Input:").SetBold().SetColspan(2).SetBackgroundColor(c1)
					})
				})
			for i, inp := range feature.Input {
				tb.AddRow(func(rb *RowBuilder) {
					rb.AddCell(func(cb *CellBuilder) {
						cb.SetText(fmt.Sprintf("%d. %s", i+1, inp.Name)).SetBold().SetColspan(2).SetBackgroundColor(c2)
					})
				})
				if inp.Fields != nil && !reflect.ValueOf(inp.Fields).IsZero() {
					tb.AddRow(func(rb *RowBuilder) {
						rb.
							AddCell(func(cb *CellBuilder) { cb.SetText("Fields").SetBold().SetWidth(wd) }).
							AddCell(func(cb *CellBuilder) { cb.SetText(inp.Fields) })
					})
				}
				if inp.Constraints != nil && !reflect.ValueOf(inp.Constraints).IsZero() {
					tb.AddRow(func(rb *RowBuilder) {
						rb.
							AddCell(func(cb *CellBuilder) { cb.SetText("Constraints").SetBold().SetWidth(wd) }).
							AddCell(func(cb *CellBuilder) { cb.SetText(inp.Constraints) })
					})
				}
				if inp.Remarks != nil && !reflect.ValueOf(inp.Remarks).IsZero() {
					tb.AddRow(func(rb *RowBuilder) {
						rb.
							AddCell(func(cb *CellBuilder) { cb.SetText("Remarks").SetBold().SetWidth(wd) }).
							AddCell(func(cb *CellBuilder) { cb.SetText(inp.Remarks) })
					})
				}
			}
		})
	}

	/* -------------------------------- SCENARIO -------------------------------- */
	if len(feature.Scenarios) > 0 {
		docb.AddParagraph().AddTable(func(tb *TableBuilder) {
			tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) }).
				AddRow(func(rb *RowBuilder) {
					rb.AddCell(func(cb *CellBuilder) {
						cb.SetText("Scenarios and Processing Logic:").SetBold().SetColspan(2).SetBackgroundColor(c1)
					})
				})
			for i, scn := range feature.Scenarios {
				tb.AddRow(func(rb *RowBuilder) {
					rb.AddCell(func(cb *CellBuilder) {
						cb.SetText(fmt.Sprintf("%d. %s", i+1, scn.Name)).SetBold().SetColspan(2).SetBackgroundColor(c2)
					})
				})
				for _, action := range scn.Desc {
					keyword, others := b.splitGherkinWord(action)
					tb.AddRow(func(rb *RowBuilder) {
						rb.
							AddCell(func(cb *CellBuilder) { cb.SetText(keyword).SetBold().SetAlignment(wml.ST_JcLeft).SetWidth(10) }).
							AddCell(func(cb *CellBuilder) { cb.SetText(others) })
					})
				}
			}
		})

	}

	/* --------------------------------- REMARKS -------------------------------- */
	if feature.Remarks != nil && !reflect.ValueOf(feature.Remarks).IsZero() {
		docb.AddParagraph().AddTable(func(tb *TableBuilder) {
			tb.SetWidthPercent(100).SetBorders(func(b *Borders) { b.SetBorderAll(bs, bc, bt) }).
				AddRow(func(rb *RowBuilder) {
					rb.AddCell(func(cb *CellBuilder) {
						cb.SetText("Remarks:").SetBold().SetColspan(2).SetBackgroundColor(c1)
					})
				}).AddRow(func(rb *RowBuilder) {
				rb.AddCell(func(cb *CellBuilder) {
					cb.SetText(feature.Remarks).SetBullet(&nd)
				})
			})
		})
	}

	/* --------------------------------- OTHERS --------------------------------- */
	// if !b.isInterfaceBlank(feature.Others.Limits) || !b.isInterfaceBlank(feature.Others.Reference) || !b.isInterfaceBlank(feature.Others.Remarks) {
	// 	docb.AddParagraph()
	// 	if !b.isInterfaceBlank(feature.Others.Reference) {
	// 		docb.AddTable(func(tb *TableBuilder) {

	// 		})
	// 	}
	// }

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

func (b *Builder) isInterfaceBlank(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsZero()
}
