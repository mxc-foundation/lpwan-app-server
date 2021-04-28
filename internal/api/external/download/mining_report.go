package download

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
)

// Server defines the download service Server API structure
type Server struct {
	financeReportCli pb.FinanceReportServiceClient
	auth             auth.Authenticator
	server           string
}

// NewServer creates a new download service server
func NewServer(mxpCli *mxpcli.Client, auth auth.Authenticator, server string) *Server {
	return &Server{
		financeReportCli: mxpCli.GetFianceReportClient(),
		auth:             auth,
		server:           server,
	}
}

// GetFiatCurrencyList returns a list of fiat currecy supported by supernode
func (s *Server) GetFiatCurrencyList(ctx context.Context, req *api.GetFiatCurrencyListRequest) (*api.GetFiatCurrencyListResponse, error) {
	resp, err := s.financeReportCli.GetSupportedFiatCurrencyList(ctx, &pb.GetSupportedFiatCurrencyListRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to get supported fiat currency list: %v", err)
	}

	response := &api.GetFiatCurrencyListResponse{
		FiatCurrencyList: []*api.FiatCurrency{},
	}
	for _, v := range resp.FiatCurrency {
		response.FiatCurrencyList = append(response.FiatCurrencyList, &api.FiatCurrency{
			Id:          v.Id,
			Description: v.Description,
		})
	}

	return response, nil
}

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
	titleFontSize      float64
	contentFontSize    float64
	disclaimerFontSize float64
	bannerHeight       float64
}

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
}

func addReportTitle(pdf *gofpdf.Fpdf, f pdfFormat, supernode, username string) {
	// draw banner
	pdf.SetFillColor(28, 20, 120)
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{f.pageWidth, 0},
		{f.pageHeight, f.bannerHeight},
		{0, f.bannerHeight},
	}, "F")

	// add title
	pdf.SetFont("arial", "B", f.titleFontSize)
	pdf.SetTextColor(255, 255, 255)
	pdf.Text(f.gridWidth, f.bannerHeight/2.0+f.titleFontSize/2.0, "Mining Income Report ")

	// add supernode name, username
	pdf.SetFont("arial", "", f.contentFontSize)
	pdf.SetTextColor(255, 255, 255)
	pdf.MoveTo(f.pageWidth-f.indentationRight-5*f.gridWidth, 0.5*f.gridHeight)
	pdf.MultiCell(5*f.gridWidth, f.contentFontSize+f.lineSpacing, fmt.Sprintf("Supernode: %s\n User: %s",
		supernode, username), gofpdf.BorderNone, gofpdf.AlignCenter, false)
}

func addReportTable(pdf *gofpdf.Fpdf, f pdfFormat, table [][]string, cellWidth []float64) error {
	// insanity check
	if len(table[0]) != len(cellWidth) {
		return fmt.Errorf("length of cellWidth must be same as length of table columns")
	}

	tableWidth := f.pageWidth - f.indentationLeft/2 - f.indentationRight/2
	tableX := f.indentationLeft / 2
	tableY := f.bannerHeight + f.lineSpacing
	pdf.Rect(tableX, tableY, tableWidth, 3, "F")

	pdf.SetFont("times", "", f.contentFontSize)
	pdf.SetTextColor(51, 51, 51)

	var moveToX, moveToY float64
	cellHight := 1.2 * f.contentFontSize
	lineMax := int((f.pageHeight-tableY-f.indentationBottom)/(2*cellHight)) - 2
	for row, rowContent := range table {
		if row == 0 {
			moveToY = tableY
		} else {
			moveToY += 2 * cellHight
		}
		moveToX = tableX
		for column, item := range rowContent {
			if column == 0 {
				moveToX = tableX
			} else {
				moveToX += cellWidth[column-1] + f.contentFontSize/2
			}
			pdf.MoveTo(moveToX, moveToY)
			pdf.MultiCell(cellWidth[column], cellHight, item,
				gofpdf.BorderNone, gofpdf.AlignRight, false)
		}
		pdf.Line(tableX, moveToY+2*cellHight, tableX+tableWidth, moveToY+2*cellHight)
		if row != 0 && (row-lineMax*(row/lineMax)) == 0 {
			pdf.SetFont("times", "", f.disclaimerFontSize)
			pdf.MoveTo(tableX, moveToY+2*cellHight+f.lineSpacing)
			pdf.MultiCell(f.pageWidth-f.indentationLeft-f.indentationRight, f.disclaimerFontSize,
				"*This information is provided to the best of our current knowledge & ability. "+
					"The MXC Foundation takes no legal responsibility for the accuracy or timeliness of this data. "+
					"On-chain data is used to compile this information", gofpdf.BorderNone, gofpdf.AlignLeft, false)
			pdf.AddPage()
			pdf.SetTextColor(51, 51, 51)
			pdf.Text(f.pageWidth/2, f.pageHeight-2*f.disclaimerFontSize, fmt.Sprintf("%d", pdf.PageNo()))
			moveToY = 0
			pdf.Rect(tableX, 2*cellHight-f.lineSpacing, tableWidth, 3, "F")
			pdf.SetFont("times", "", f.contentFontSize)
		}
	}

	pdf.SetFont("times", "", f.disclaimerFontSize)
	pdf.SetTextColor(51, 51, 51)
	pdf.MoveTo(tableX, moveToY+2*cellHight+f.lineSpacing)
	pdf.MultiCell(f.pageWidth-f.indentationLeft-f.indentationRight, f.disclaimerFontSize,
		"*This information is provided to the best of our current knowledge & ability. "+
			"The MXC Foundation takes no legal responsibility for the accuracy or timeliness of this data. "+
			"On-chain data is used to compile this information", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	return nil
}

// MiningReportPDF formats mining data into pdf with given filtering conditions then send to client in stream
func (s *Server) MiningReportPDF(req *api.MiningReportRequest, server api.DownloadService_MiningReportPDFServer) error {
	cred, err := s.auth.GetCredentials(server.Context(), auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsOrgAdmin {
		return status.Errorf(codes.PermissionDenied, "permission denied")
	}

	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeA4, "")
	// configure format
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
	format.disclaimerFontSize = format.contentFontSize * 0.8
	format.bannerHeight = format.indentationUp

	pdf.AddPage()
	pdf.SetFont("times", "", format.disclaimerFontSize)
	pdf.SetTextColor(51, 51, 51)
	pdf.Text(format.pageWidth/2, format.pageHeight-2*format.disclaimerFontSize, fmt.Sprintf("%d", pdf.PageNo()))

	addReportTitle(pdf, format, s.server, cred.Username)

	for _, item := range req.Currency {
		switch item {
		case "ETH_MXC":
			tableContent := [][]string{{}}
			tableContent[0] = []string{
				"\nDate",
				"\nMXC Mined",
				//"MXC Settlement Price",
				"\nFiat Currency Mined",
				"\nOnline Seconds",
			}
			res, err := s.financeReportCli.GetMXCMiningReportByDate(server.Context(), &pb.GetMXCMiningReportByDateRequest{
				OrganizationId: req.OrganizationId,
				Start:          req.Start,
				End:            req.End,
				FiatCurrency:   req.FiatCurrency,
			})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to get MXC mining report : %v", err)
			}
			for _, v := range res.MiningRecordList {
				y, m, d := v.DateTime.AsTime().Date()
				dateStr := fmt.Sprintf("%04d-%02d-%02d", y, m, d)
				tableContent = append(tableContent, []string{
					dateStr,
					v.MXCMined,
					//v.MXCSettlementPrice,
					v.FiatCurrencyMined,
					fmt.Sprintf("%d", v.OnlineSeconds),
				})
			}
			cellWidth := []float64{
				5.5 * format.contentFontSize, // date
				14 * format.contentFontSize,  // mxc mined
				//10 * format.contentFontSize,  // mxc price
				20 * format.contentFontSize, // fiat mined
				11 * format.contentFontSize, //online seconds
			}
			if err = addReportTable(pdf, format, tableContent, cellWidth); err != nil {
				return status.Errorf(codes.Internal, "%v", err)
			}
		}
	}
	// output to file
	sy, sm, sd := req.Start.AsTime().Date()
	ey, em, ed := req.End.AsTime().Date()
	filename := fmt.Sprintf("mining_report_org_%d_%s_%s_%s", req.OrganizationId, req.FiatCurrency,
		fmt.Sprintf("%d-%d-%d", sy, sm, sd), fmt.Sprintf("%d-%d-%d", ey, em, ed))
	// drawGrid(pdf, format)
	err = pdf.OutputFileAndClose(filepath.Join("/tmp/mining-report", filename))

	return err
}

// MiningReportCSV formats mining data into csv with given filtering conditions then send to client in stream
func (s *Server) MiningReportCSV(req *api.MiningReportRequest, server api.DownloadService_MiningReportCSVServer) error {
	cred, err := s.auth.GetCredentials(server.Context(), auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsOrgAdmin {
		return status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// write report to csv file
	sy, sm, sd := req.Start.AsTime().Date()
	ey, em, ed := req.End.AsTime().Date()
	filename := fmt.Sprintf("mining_report_org_%d_%s_%s_%s.csv",
		req.OrganizationId, req.FiatCurrency, fmt.Sprintf("%d-%d-%d", sy, sm, sd), fmt.Sprintf("%d-%d-%d", ey, em, ed))
	buff := bytes.Buffer{}
	buff.Reset()

	buffFile := bytes.Buffer{}
	buffFile.Reset()

	w := csv.NewWriter(&buff)
	wFile := csv.NewWriter(&buffFile)

	for _, item := range req.Currency {
		switch item {
		case "ETH_MXC":
			res, err := s.financeReportCli.GetMXCMiningReportByDate(server.Context(), &pb.GetMXCMiningReportByDateRequest{
				OrganizationId: req.OrganizationId,
				Start:          req.Start,
				End:            req.End,
				FiatCurrency:   req.FiatCurrency,
			})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to get MXC mining report : %v", err)
			}

			err = w.Write([]string{
				"Date",
				"MXC Mined",
				"MXC Settlement Price",
				"Fiat Currency Mined",
				"Online Seconds"})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to write title to csv file: %v", err)
			}

			err = wFile.Write([]string{
				"Date",
				"MXC Mined",
				"MXC Settlement Price",
				"Fiat Currency Mined",
				"Online Seconds"})
			if err != nil {
				logrus.Debugf("Error occurs when writing title to buffer: %v", err)
			}

			for _, v := range res.MiningRecordList {
				y, m, d := v.DateTime.AsTime().Date()
				dateStr := fmt.Sprintf("%d-%d-%d", y, m, d)
				err = w.Write([]string{
					dateStr,
					v.MXCMined,
					v.MXCSettlementPrice,
					v.FiatCurrencyMined,
					fmt.Sprintf("%d", v.OnlineSeconds)})
				if err != nil {
					return status.Errorf(codes.Internal, "failed to write value to csv file: %v", err)
				}

				err = wFile.Write([]string{
					dateStr,
					v.MXCMined,
					v.MXCSettlementPrice,
					v.FiatCurrencyMined,
					fmt.Sprintf("%d", v.OnlineSeconds)})
				if err != nil {
					logrus.Debugf("Error occurs when writing value to buffer: %v", err)
				}
			}
		}
	}

	if err := w.Write([]string{"*This information is provided to the best of our current knowledge & ability. " +
		"The MXC Foundation takes no legal responsibility for the accuracy or timeliness of this data. " +
		"On-chain data is used to compile this information"}); err != nil {
		return status.Errorf(codes.Internal, "failed to write disclaimer to file: %v", err)
	}
	if err := wFile.Write([]string{"*This information is provided to the best of our current knowledge & ability. " +
		"The MXC Foundation takes no legal responsibility for the accuracy or timeliness of this data. " +
		"On-chain data is used to compile this information"}); err != nil {
		return status.Errorf(codes.Internal, "Error occurs when writing disclaimer to buffer: %v", err)
	}

	w.Flush()
	wFile.Flush()
	messageByte, err := buff.ReadBytes('\n')
	for err == nil {
		_ = server.Send(&api.MiningReportResponse{FileChunk: messageByte})
		messageByte, err = buff.ReadBytes('\n')
	}
	if err == io.EOF {
		// just for debugging, save file locally
		if err := ioutil.WriteFile(filepath.Join("/tmp/mining-report", filename), buffFile.Bytes(), os.ModePerm); err != nil {
			logrus.Debugf("Error occurs when writing file %s: %v", filename, err)
		}

		return nil
	} else {
		return status.Errorf(codes.Internal, "failed to download csv file: %v", err)
	}
}
