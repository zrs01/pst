package converter

import (
	"fmt"
	"log"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
)

func GenSimpleDoc() {
	doc := document.New()

	para := doc.AddParagraph()
	run := para.AddRun()
	para.SetStyle("Title")
	run.AddText("Simple Document Formatting")

	para = doc.AddParagraph()
	para.SetStyle("Heading1")
	run = para.AddRun()
	run.AddText("Some Heading Text")

	para = doc.AddParagraph()
	para.SetStyle("Heading2")
	run = para.AddRun()
	run.AddText("Some Heading Text")

	para = doc.AddParagraph()
	para.SetStyle("Heading3")
	run = para.AddRun()
	run.AddText("Some Heading Text")

	para = doc.AddParagraph()
	para.Properties().SetFirstLineIndent(0.5 * measurement.Inch)

	run = para.AddRun()
	run.AddText("A run is a string of characters with the same formatting. ")

	run = para.AddRun()
	run.Properties().SetBold(true)
	run.Properties().SetFontFamily("Courier")
	run.Properties().SetSize(15)
	run.Properties().SetColor(color.Red)
	run.AddText("Multiple runs with different formatting can exist in the same paragraph. ")

	run = para.AddRun()
	run.AddText("Adding breaks to a run will insert line breaks after the run. ")
	run.AddBreak()
	run.AddBreak()

	run = createParaRun(doc, "Runs support styling options:")

	run = createParaRun(doc, "small caps")
	run.Properties().SetSmallCaps(true)

	run = createParaRun(doc, "strike")
	run.Properties().SetStrikeThrough(true)

	run = createParaRun(doc, "double strike")
	run.Properties().SetDoubleStrikeThrough(true)

	run = createParaRun(doc, "outline")
	run.Properties().SetOutline(true)

	run = createParaRun(doc, "emboss")
	run.Properties().SetEmboss(true)

	run = createParaRun(doc, "shadow")
	run.Properties().SetShadow(true)

	run = createParaRun(doc, "imprint")
	run.Properties().SetImprint(true)

	run = createParaRun(doc, "highlighting")
	run.Properties().SetHighlight(wml.ST_HighlightColorYellow)

	run = createParaRun(doc, "underline")
	run.Properties().SetUnderline(wml.ST_UnderlineWavyDouble, color.Red)

	run = createParaRun(doc, "text effects")
	run.Properties().SetEffect(wml.ST_TextEffectAntsRed)

	nd := doc.Numbering.Definitions()[0]

	for i := 1; i < 5; i++ {
		p := doc.AddParagraph()
		p.SetNumberingLevel(i - 1)
		p.SetNumberingDefinition(nd)
		run := p.AddRun()
		run.AddText(fmt.Sprintf("Level %d", i))
	}
	doc.SaveToFile("simple.docx")
}

func createParaRun(doc *document.Document, s string) document.Run {
	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText(s)
	return run
}

func GenSimpleTable() {
	doc := document.New()

	// First Table
	{
		table := doc.AddTable()
		// width of the page
		table.Properties().SetWidthPercent(100)
		// with thick borers
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, 2*measurement.Point)

		row := table.AddRow()
		run := row.AddCell().AddParagraph().AddRun()
		run.AddText("Name")
		run.Properties().SetHighlight(wml.ST_HighlightColorYellow)
		row.AddCell().AddParagraph().AddRun().AddText("John Smith")
		row = table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText("Street Address")
		row.AddCell().AddParagraph().AddRun().AddText("111 Country Road")
	}

	doc.AddParagraph() // break up the consecutive tables

	// Second Table
	{
		table := doc.AddTable()
		// 4 inches wide
		table.Properties().SetWidth(4 * measurement.Inch)
		borders := table.Properties().Borders()
		// thin borders
		borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)

		row := table.AddRow()
		cell := row.AddCell()
		// column span / merged cells
		cell.Properties().SetColumnSpan(2)

		run := cell.AddParagraph().AddRun()
		run.AddText("Cells can span multiple columns")

		row = table.AddRow()
		cell = row.AddCell()
		cell.Properties().SetVerticalMerge(wml.ST_MergeRestart)
		cell.AddParagraph().AddRun().AddText("Vertical Merge")
		row.AddCell().AddParagraph().AddRun().AddText("")

		row = table.AddRow()
		cell = row.AddCell()
		cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
		cell.AddParagraph().AddRun().AddText("Vertical Merge 2")
		row.AddCell().AddParagraph().AddRun().AddText("")

		row = table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText("Street Address")
		row.AddCell().AddParagraph().AddRun().AddText("111 Country Road")
	}

	doc.AddParagraph()

	// Third Table
	{
		table := doc.AddTable()
		table.Properties().SetWidthPercent(100)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

		hdrRow := table.AddRow()

		cell := hdrRow.AddCell()
		cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
		cellPara := cell.AddParagraph()
		cellPara.Properties().SetAlignment(wml.ST_JcLeft)
		cellPara.AddRun().AddText("Left Align")

		cell = hdrRow.AddCell()
		cell.Properties().SetShading(wml.ST_ShdThinDiagStripe, color.Red, color.LightGray)
		cellPara = cell.AddParagraph()
		cellPara.Properties().SetAlignment(wml.ST_JcCenter)
		cellPara.AddRun().AddText("Center Align")

		cell = hdrRow.AddCell()
		cell.Properties().SetShading(wml.ST_ShdPct20, color.Red, color.LightGray)
		cellPara = cell.AddParagraph()
		cellPara.Properties().SetAlignment(wml.ST_JcRight)
		cellPara.AddRun().AddText("Right Align")

		veryLightGray := color.RGB(240, 240, 240)
		for i := 0; i < 5; i++ {
			row := table.AddRow()
			for j := 0; j < 3; j++ {
				cell := row.AddCell()
				// shade every other row
				if i%2 == 0 {
					cell.Properties().SetShading(wml.ST_ShdSolid, veryLightGray, color.Auto)
				}
				cell.AddParagraph().AddRun()
			}
		}
	}

	doc.AddParagraph()
	// Fourth Table
	{
		table := doc.AddTable()
		table.Properties().SetWidthPercent(50)
		table.Properties().SetAlignment(wml.ST_JcTableCenter)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

		row := table.AddRow()
		row.Properties().SetHeight(2*measurement.Inch, wml.ST_HeightRuleExact)

		cell := row.AddCell()
		cell.Properties().SetVerticalAlignment(wml.ST_VerticalJcCenter)

		para := cell.AddParagraph()
		para.Properties().SetAlignment(wml.ST_JcCenter)
		run := para.AddRun()
		run.AddText("hello world")
	}

	doc.AddParagraph()
	// Fifth Table
	{
		table := doc.AddTable()
		table.Properties().SetWidthPercent(90)
		table.Properties().SetAlignment(wml.ST_JcTableCenter)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

		row := table.AddRow()

		cell := row.AddCell()
		cell.Properties().SetWidthPercent(25)
		para := cell.AddParagraph()
		run := para.AddRun()
		run.AddText("hello")

		cell = row.AddCell()
		cell.Properties().SetWidthPercent(75)
		para = cell.AddParagraph()
		run = para.AddRun()
		run.AddText("world")

		// start a new table
		doc.AddParagraph()
		table = doc.AddTable()
		table.Properties().SetWidthPercent(90)
		table.Properties().SetAlignment(wml.ST_JcTableCenter)
		borders = table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

		row = table.AddRow()

		cell = row.AddCell()
		cell.Properties().SetWidth(0.25 * measurement.Inch)
		para = cell.AddParagraph()
		run = para.AddRun()
		run.AddText("hello")

		cell = row.AddCell()
		cell.Properties().SetWidth(2.5 * measurement.Inch)
		para = cell.AddParagraph()
		run = para.AddRun()
		run.AddText("world")

	}
	doc.AddParagraph()
	// Seventh Table - Styled
	{
		// construct a table style
		ts := doc.Styles.AddStyle("MyTableStyle", wml.ST_StyleTypeTable, false)
		tp := ts.TableProperties()
		tp.SetRowBandSize(1)
		tp.SetColumnBandSize(1)
		tp.SetTableIndent(measurement.Zero)

		// horizomntal banding
		s := ts.TableConditionalFormatting(wml.ST_TblStyleOverrideTypeBand1Horz)
		s.CellProperties().SetShading(wml.ST_ShdSolid, color.LightBlue, color.Red)

		// first row bold
		s = ts.TableConditionalFormatting(wml.ST_TblStyleOverrideTypeFirstRow)
		s.RunProperties().SetBold(true)

		// last row bold
		s = ts.TableConditionalFormatting(wml.ST_TblStyleOverrideTypeLastRow)
		s.RunProperties().SetBold(true)
		cb := s.CellProperties().Borders()
		cb.SetTop(wml.ST_BorderDouble, color.Black, 0.5*measurement.Point)

		tp.Borders().SetAll(wml.ST_BorderSingle, color.Blue, 0.5*measurement.Point)

		table := doc.AddTable()

		table.Properties().SetLayout(wml.ST_TblLayoutTypeFixed)
		table.Properties().SetWidthPercent(90)
		table.Properties().SetStyle("MyTableStyle")
		look := table.Properties().TableLook()
		look.SetFirstColumn(true)
		look.SetFirstRow(true)
		look.SetHorizontalBanding(true)

		for r := 0; r < 5; r++ {
			row := table.AddRow()
			for c := 0; c < 5; c++ {
				cell := row.AddCell()
				cell.AddParagraph().AddRun().AddText(fmt.Sprintf("row %d col %d", r+1, c+1))
			}
		}
	}

	// Sixth Table - Insert
	{
		doc.InsertParagraphBefore(doc.Paragraphs()[5])
		table := doc.InsertTableBefore(doc.Paragraphs()[5])
		table.Properties().SetWidthPercent(90)
		table.Properties().SetAlignment(wml.ST_JcTableCenter)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

		row := table.AddRow()

		cell := row.AddCell()
		cell.Properties().SetWidth(0.25 * measurement.Inch)
		para := cell.AddParagraph()
		run := para.AddRun()
		run.AddText("Insert")
		cell = row.AddCell()
		cell.Properties().SetWidth(0.25 * measurement.Inch)
		para = cell.AddParagraph()
		run = para.AddRun()
		run.AddText("new")
		cell = row.AddCell()
		cell.Properties().SetWidth(2.75 * measurement.Inch)
		para = cell.AddParagraph()
		run = para.AddRun()
		run.AddText("table")
	}
	if err := doc.Validate(); err != nil {
		log.Fatalf("error during validation: %s", err)
	}
	doc.SaveToFile("tables.docx")

}
