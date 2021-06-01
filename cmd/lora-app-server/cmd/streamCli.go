package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
)

var streamCliCmd = &cobra.Command{
	Use:   "stream-cli [JWT] [ORG_ID] [CRYPTO_CURRENCY] [FIAT_CURRENCY] [START] [END] pdf/csv",
	Short: "Call API with given URL and save file locally",
	Run: func(cmd *cobra.Command, args []string) {
		jwt := args[0]
		organizationIDStr := args[1]
		organizationID, err := strconv.ParseInt(organizationIDStr, 10, 64)
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		currency := args[2]
		fiatCurrency := args[3]
		start := args[4]
		startTime, err := time.Parse("2006-01-02T15:04:05Z07:00", start)
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		end := args[5]
		endTime, err := time.Parse("2006-01-02T15:04:05Z07:00", end)
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		choice := args[6]

		conn, err := grpccli.Connect(grpccli.ConnectionOpts{
			Server:  "localhost:8080",
			CACert:  "",
			TLSCert: "",
			TLSKey:  "",
		})
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		md := metadata.Pairs("authorization", jwt)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		request := &api.MiningReportRequest{
			OrganizationId: organizationID,
			Currency:       []string{currency},
			FiatCurrency:   fiatCurrency,
			Start:          timestamppb.New(startTime),
			End:            timestamppb.New(endTime),
			Decimals:       4,
			Jwt:            jwt,
		}

		if choice == "pdf" {
			exportPDF(ctx, api.NewReportServiceClient(conn), request)
		} else if choice == "csv" {
			exportCSV(ctx, api.NewReportServiceClient(conn), request)
		}
	},
}

func exportPDF(ctx context.Context, cli api.ReportServiceClient, request *api.MiningReportRequest) {
	fullData := []byte{}
	resCli, err := cli.MiningReportPDF(ctx, request)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	for {
		data, err := resCli.Recv()
		if err != nil {
			logrus.Fatalf("%v", err)
		}
		fullData = append(fullData, data.Data...)
		if data.Finish {
			log.Println("Stream is over")
			// save fullData to file
			if err = ioutil.WriteFile("report.pdf", fullData, os.ModePerm); err != nil {
				logrus.Fatalf("%v", err)
			}
			break
		}
	}
}

func exportCSV(ctx context.Context, cli api.ReportServiceClient, request *api.MiningReportRequest) {
	fullData := []byte{}
	resCli, err := cli.MiningReportCSV(ctx, request)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	for {
		data, err := resCli.Recv()
		if err != nil {
			logrus.Fatalf("%v", err)
		}
		fullData = append(fullData, data.Data...)
		if data.Finish {
			log.Println("Stream is over")
			// save fullData to file
			if err = ioutil.WriteFile("report.csv", fullData, os.ModePerm); err != nil {
				logrus.Fatalf("%v", err)
			}
			break
		}
	}
}
