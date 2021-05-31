package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
)

var streamCliCmd = &cobra.Command{
	Use:   "stream-cli [JWT] [ORG_ID] [CRYPTO_CURRENCY] [FIAT_CURRENCY] [START] [END] pdf/csv",
	Short: "Call API with given URL and save file locally",
	Run: func(cmd *cobra.Command, args []string) {
		jwt := args[0]
		organizationIdStr := args[1]
		organizationId, err := strconv.ParseInt(organizationIdStr, 10, 64)
		if err != nil {
			logrus.Fatal("%v", err)
		}

		currency := args[2]
		fiatCurrency := args[3]
		start := args[4]
		startTime, err := time.Parse("2006-01-02T15:04:05Z07:00", start)
		if err != nil {
			logrus.Fatal("%v", err)
		}

		end := args[5]
		endTime, err := time.Parse("2006-01-02T15:04:05Z07:00", end)
		if err != nil {
			logrus.Fatal("%v", err)
		}

		choice := args[6]

		conn, err := grpccli.Connect(grpccli.ConnectionOpts{
			Server:  "localhost:8080",
			CACert:  "",
			TLSCert: "",
			TLSKey:  "",
		})
		if err != nil {
			logrus.Fatal("%v", err)
		}

		md := metadata.Pairs("authorization", jwt)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		request := &api.MiningReportRequest{
			OrganizationId: organizationId,
			Currency:       []string{currency},
			FiatCurrency:   fiatCurrency,
			Start:          timestamppb.New(startTime),
			End:            timestamppb.New(endTime),
			Decimals:       4,
			Jwt:            jwt,
		}

		if choice == "pdf" {
			exportPDF(conn, ctx, request)
		} else if choice == "csv" {
			exportCSV(conn, ctx, request)
		}
	},
}

func exportPDF(conn *grpc.ClientConn, ctx context.Context, request *api.MiningReportRequest) {
	fullData := []byte{}
	cli := api.NewReportServiceClient(conn)
	resCli, err := cli.MiningReportPDF(ctx, request)
	if err != nil {
		logrus.Fatal("%v", err)
	}
	for {
		data, err := resCli.Recv()
		if err != nil {
			logrus.Fatal("%v", err)
		}
		fullData = append(fullData, data.Data...)
		if data.Finish == true {
			log.Println("Stream is over")
			// save fullData to file
			if err = ioutil.WriteFile("report.pdf", fullData, os.ModePerm); err != nil {
				logrus.Fatal("%v", err)
			}
			break
		}
	}
}

func exportCSV(conn *grpc.ClientConn, ctx context.Context, request *api.MiningReportRequest) {
	fullData := []byte{}
	cli := api.NewReportServiceClient(conn)
	resCli, err := cli.MiningReportCSV(ctx, request)
	if err != nil {
		logrus.Fatal("%v", err)
	}
	for {
		data, err := resCli.Recv()
		if err != nil {
			logrus.Fatal("%v", err)
		}
		fullData = append(fullData, data.Data...)
		if data.Finish == true {
			log.Println("Stream is over")
			// save fullData to file
			if err = ioutil.WriteFile("report.csv", fullData, os.ModePerm); err != nil {
				logrus.Fatal("%v", err)
			}
			break
		}
	}
}
