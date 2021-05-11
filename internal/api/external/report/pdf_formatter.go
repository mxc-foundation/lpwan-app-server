package report

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type pdfFormat struct {
	gridWidth          float64
	gridHeight         float64
	pageWidth          float64
	pageHeight         float64
	indentationUp      float64
	indentationBottom  float64
	indentationLeft    float64
	indentationRight   float64
	lineSpacing        float64
	charSpacing        float64
	titleFontSize      float64
	contentFontSize    float64
	tableWidth         float64
	disclaimerFontSize float64
	bannerHeight       float64
}

func defaultPDFConfiguration(pdf *gofpdf.Fpdf) pdfFormat {
	w, h := pdf.GetPageSize()
	format := pdfFormat{
		pageWidth:  w,
		pageHeight: h,
		gridWidth:  h / 40.0,
		gridHeight: h / 40.0,
	}
	format.indentationUp = format.gridHeight * 3
	format.indentationBottom = format.gridHeight * 2
	format.indentationLeft = format.gridWidth * 2
	format.indentationRight = format.indentationLeft
	format.lineSpacing = format.gridHeight / 2
	format.titleFontSize = format.gridHeight
	format.contentFontSize = format.gridHeight / 2
	format.charSpacing = format.contentFontSize / 2
	format.disclaimerFontSize = format.contentFontSize * 0.8
	format.bannerHeight = format.indentationUp
	format.tableWidth = format.pageWidth - format.indentationLeft/2 - format.indentationRight/2

	return format
}

/*// drawGrid is extremely useful for designing pdf pages' layout in the beginning
// commented out for passing lint check
func drawGrid(pdf *gofpdf.Fpdf, f pdfFormat) {
	fontSize := 6.0
	pdf.SetFont("courier", "", fontSize)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetDrawColor(200, 200, 200)
	for x := 0.0; x < f.pageWidth; x += f.gridWidth {
		pdf.Line(x, 0, x, f.pageHeight)
		pdf.Text(x, fontSize, fmt.Sprintf("%d", int(x)))
	}
	for y := 0.0; y < f.pageHeight; y += f.gridHeight {
		pdf.Line(0, y, f.pageWidth, y)
		pdf.Text(0, y+fontSize, fmt.Sprintf("%d", int(y)))
	}
}*/

func addNewPageWithCustomization(pdf *gofpdf.Fpdf, f pdfFormat) {
	pdf.AddPage()
	// add page number
	pdf.SetFont("times", "", f.disclaimerFontSize)
	pdf.SetTextColor(51, 51, 51)
	pdf.Text(f.pageWidth/2, f.pageHeight-2*f.disclaimerFontSize, fmt.Sprintf("%d", pdf.PageNo()))
	// add disclaimer on bottom of every page
	if pdf.PageNo() == 1 {
		return
	}
	pdf.MoveTo(f.gridWidth, f.pageHeight-f.indentationBottom-4*f.disclaimerFontSize)
	pdf.MultiCell(f.tableWidth, f.disclaimerFontSize,
		"*This information is provided to the best of our current knowledge & ability. "+
			"The MXC Foundation takes no legal responsibility for the accuracy or timeliness of this data. "+
			"On-chain data is used to compile this information", gofpdf.BorderNone, gofpdf.AlignLeft, false)
}

func addReportBanner(pdf *gofpdf.Fpdf, f pdfFormat, supernode, username string) {
	// draw banner
	pdf.SetFillColor(28, 20, 120)
	pdf.Rect(0, 0, f.pageWidth, f.bannerHeight, "F")
	// add title
	pdf.SetFont("arial", "B", f.titleFontSize)
	pdf.SetTextColor(255, 255, 255)
	pdf.Text(f.gridWidth, f.bannerHeight/2.0+f.titleFontSize/2.0, "Mining Income Report ")

	// add supernode name, username
	infoCellWidth := 10 * f.gridWidth
	pdf.SetFont("arial", "", f.contentFontSize)
	pdf.SetTextColor(255, 255, 255)
	pdf.MoveTo(f.pageWidth-f.indentationRight-infoCellWidth, 0.5*f.gridHeight)
	pdf.MultiCell(infoCellWidth, f.contentFontSize+f.lineSpacing, fmt.Sprintf("Supernode: %s\n User: %s",
		supernode, username), gofpdf.BorderNone, gofpdf.AlignRight, false)
}

func addReportTable(pdf *gofpdf.Fpdf, f pdfFormat, table [][]string, cellWidth []float64) error {
	// insanity check
	if len(table[0]) != len(cellWidth) {
		return fmt.Errorf("length of cellWidth must be same as length of table columns")
	}

	recHeight := 3.0
	tableX := f.indentationLeft / 2
	tableY := f.bannerHeight + f.lineSpacing
	tableWidth := f.tableWidth
	tableHeight := f.pageHeight - tableY - f.indentationBottom
	cellHeight := 1.2 * f.contentFontSize
	rowHeight := 2 * cellHeight
	rowMaxPerPage := int(tableHeight / rowHeight)

	// truncate table into pages
	var pages int
	if (len(table) - (len(table)/rowMaxPerPage)*rowMaxPerPage) == 0 {
		pages = len(table) / rowMaxPerPage
	} else {
		pages = len(table)/rowMaxPerPage + 1
	}
	var moveToX, moveToY float64
	for p := 0; p < pages; p++ {
		if p > 0 {
			addNewPageWithCustomization(pdf, f)
			pdf.Rect(tableX, f.gridHeight-recHeight, tableWidth, recHeight, "F")
			moveToY = f.gridHeight
		} else {
			pdf.Rect(tableX, tableY-recHeight, tableWidth, recHeight, "F")
			moveToY = tableY
		}
		end := (p + 1) * rowMaxPerPage
		if p == pages-1 {
			end = len(table)
		}
		for _, rows := range table[p*rowMaxPerPage : end] {
			moveToX = tableX
			pdf.SetFont("times", "", f.contentFontSize)
			pdf.SetTextColor(51, 51, 51)
			for column, item := range rows {
				pdf.MoveTo(moveToX, moveToY)
				pdf.MultiCell(cellWidth[column], cellHeight, item,
					gofpdf.BorderNone, gofpdf.AlignRight, false)
				moveToX += cellWidth[column] + f.charSpacing
			}
			pdf.SetDrawColor(200, 200, 200)
			pdf.Line(tableX, moveToY+rowHeight, tableX+tableWidth, moveToY+rowHeight)
			moveToY += rowHeight
		}
	}

	return nil
}
