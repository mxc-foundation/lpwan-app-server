package report

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
)

// Server defines the download service Server API structure
type Server struct {
	financeReportCli pb.FinanceReportServiceClient
	auth             auth.Authenticator
	server           string
}

// NewServer creates a new download service server
func NewServer(mxpCli pb.FinanceReportServiceClient, auth auth.Authenticator, server string) *Server {
	return &Server{
		financeReportCli: mxpCli,
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
func (s *Server) MiningReportPDF(req *api.MiningReportRequest, srv api.ReportService_MiningReportPDFServer) error {
	cred, err := s.auth.GetCredentials(srv.Context(), auth.NewOptions().WithOrgID(req.OrganizationId))
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
	format := defaultPDFConfiguration(pdf)

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
			res, err := s.financeReportCli.GetMXCMiningReportByDate(srv.Context(), &pb.GetMXCMiningReportByDateRequest{
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
	// drawGrid(pdf, format)
	buff := bytes.Buffer{}
	buff.Reset()
	err = pdf.Output(&buff)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to output report content to buffer: %v", err)
	}

	return sendStream(buff, srv.Send)
}

func sendStream(data bytes.Buffer, send func(response *api.MiningReportResponse) error) (err error) {
	for {
		rb := make([]byte, 65535)
		n, err := io.ReadFull(&data, rb)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				if err = send(&api.MiningReportResponse{
					Data:   rb[:n],
					Finish: true,
				}); err != nil {
					return status.Errorf(codes.Internal, "server failed to send report data: %v", err)
				}
				break
			}
			return status.Errorf(codes.Internal, "failed to read from buffer: %v", err)
		}
		if err = send(&api.MiningReportResponse{
			Data:   rb,
			Finish: false,
		}); err != nil {
			return status.Errorf(codes.Internal, "server failed to send report data: %v", err)
		}
	}
	return nil
}

// MiningReportCSV formats mining data into csv with given filtering conditions then send to client in stream
func (s *Server) MiningReportCSV(req *api.MiningReportRequest, srv api.ReportService_MiningReportCSVServer) error {
	cred, err := s.auth.GetCredentials(srv.Context(), auth.NewOptions().WithOrgID(req.OrganizationId))
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

	buff := bytes.Buffer{}
	buff.Reset()
	wFile := csv.NewWriter(&buff)

	for _, item := range req.Currency {
		switch item {
		case "ETH_MXC":
			res, err := s.financeReportCli.GetMXCMiningReportByDate(srv.Context(), &pb.GetMXCMiningReportByDateRequest{
				OrganizationId: req.OrganizationId,
				Start:          req.Start,
				End:            req.End,
				FiatCurrency:   req.FiatCurrency,
				Decimals:       decimals,
			})
			if err != nil {
				return status.Errorf(codes.Internal, "failed to get MXC mining report : %v", err)
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
	if err := wFile.Write([]string{"*This information is provided to the best of our current knowledge & ability. " +
		"The MXC Foundation takes no legal responsibility for the accuracy or timeliness of this data. " +
		"On-chain data is used to compile this information"}); err != nil {
		return status.Errorf(codes.Internal, "Error occurs when writing disclaimer to buffer: %v", err)
	}
	wFile.Flush()

	return sendStream(buff, srv.Send)
}
