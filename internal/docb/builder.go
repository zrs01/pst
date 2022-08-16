package docb

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rotisserie/eris"
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
	docb, err := NewDocumentBuilder(b.dfile)
	if err != nil {
		return eris.Wrap(err, "failed to create document builder")
	}

	// resolve wildcard
	files, err := b.resolveInputFile(b.ifile)
	if err != nil {
		return eris.Wrap(err, "failed to resolve the source file")
	}
	for _, file := range *files {
		_, err := b.loadData(file)
		if err != nil {
			return eris.Wrap(err, "failed to load the file")
		}

		// header
		docb.AddParagraph(func(p *ParagraphBuilder) {
			p.SetStyle("Heading1").SetText("PROGRAM DESCRIPTON")
		}).Build()

	}
	docb.Document.SaveToFile(b.ofile)
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
