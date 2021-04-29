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
	"strings"

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

// MiningReportPDF formats mining data into pdf with given filtering conditions then send to client in stream
func (s *Server) MiningReportPDF(req *api.MiningReportRequest, server api.DownloadService_MiningReportPDFServer) error {
	cred, err := s.auth.GetCredentials(server.Context(), auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsOrgAdmin {
		return status.Errorf(codes.PermissionDenied, "permission denied")
	}
	decimals := req.Decimals
	if decimals == 0 {
		decimals = 4
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
	format.charSpacing = format.contentFontSize / 2
	format.disclaimerFontSize = format.contentFontSize * 0.8
	format.bannerHeight = format.indentationUp
	format.tableWidth = format.pageWidth - format.indentationLeft/2 - format.indentationRight/2

	// new page
	addNewPageWithCustomization(pdf, format)
	// add banner for first page
	addReportBanner(pdf, format, s.server, cred.Username)

	for _, item := range req.Currency {
		switch item {
		case "ETH_MXC":
			tableContent := [][]string{{}}
			tableContent[0] = []string{
				"\nDate",
				"\nMXC Mined",
				"\nMXC Close Price",
				fmt.Sprintf("\n%s Mined", strings.ToUpper(req.FiatCurrency)),
				"\nOnline Seconds",
			}
			// configure cell width
			totalUnits := (format.tableWidth - float64(len(tableContent[0])-1)*format.charSpacing) / format.contentFontSize
			dateUnits := 5.5
			mxcMinedUnits := 14.0
			onlineSecondsUnits := 10.0
			mxcClosePriceUnits := (totalUnits - dateUnits - mxcMinedUnits - onlineSecondsUnits) / 2
			fiatMinedUnits := mxcClosePriceUnits
			cellWidth := []float64{
				dateUnits * format.contentFontSize,
				mxcMinedUnits * format.contentFontSize,
				mxcClosePriceUnits * format.contentFontSize,
				fiatMinedUnits * format.contentFontSize,
				onlineSecondsUnits * format.contentFontSize,
			}
			// get table content
			res, err := s.financeReportCli.GetMXCMiningReportByDate(server.Context(), &pb.GetMXCMiningReportByDateRequest{
				OrganizationId: req.OrganizationId,
				Start:          req.Start,
				End:            req.End,
				FiatCurrency:   req.FiatCurrency,
				Decimals:       decimals,
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
					v.MXCSettlementPrice,
					v.FiatCurrencyMined,
					fmt.Sprintf("%d", v.OnlineSeconds),
				})
			}
			// add table content
			if err = addReportTable(pdf, format, tableContent, cellWidth); err != nil {
				return status.Errorf(codes.Internal, "%v", err)
			}
		}
	}
	// output to file
	sy, sm, sd := req.Start.AsTime().Date()
	ey, em, ed := req.End.AsTime().Date()
	filename := fmt.Sprintf("mining_report_%s_org_%d_%s_%s_%s.pdf", s.server, req.OrganizationId, req.FiatCurrency,
		fmt.Sprintf("%04d-%02d-%02d", sy, sm, sd), fmt.Sprintf("%04d-%02d-%02d", ey, em, ed))
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
	decimals := req.Decimals
	if decimals == 0 {
		decimals = 4
	}

	// write report to csv file
	sy, sm, sd := req.Start.AsTime().Date()
	ey, em, ed := req.End.AsTime().Date()
	filename := fmt.Sprintf("mining_report_%s_org_%d_%s_%s_%s.csv",
		s.server, req.OrganizationId, req.FiatCurrency, fmt.Sprintf("%04d-%02d-%02d", sy, sm, sd),
		fmt.Sprintf("%04d-%02d-%02d", ey, em, ed))
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
				Decimals:       decimals,
			})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to get MXC mining report : %v", err)
			}

			err = w.Write([]string{
				"Date",
				"MXC Mined",
				"MXC Close Price",
				fmt.Sprintf("%s Mined", strings.ToUpper(req.FiatCurrency)),
				"Online Seconds"})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to write title to csv file: %v", err)
			}

			err = wFile.Write([]string{
				"Date",
				"MXC Mined",
				"MXC Close Price",
				fmt.Sprintf("%s Mined", strings.ToUpper(req.FiatCurrency)),
				"Online Seconds"})
			if err != nil {
				logrus.Debugf("Error occurs when writing title to buffer: %v", err)
			}

			for _, v := range res.MiningRecordList {
				y, m, d := v.DateTime.AsTime().Date()
				dateStr := fmt.Sprintf("%04d-%02d-%02d", y, m, d)
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
