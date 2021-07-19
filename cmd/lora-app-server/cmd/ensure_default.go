package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	pb "github.com/mxc-foundation/lpwan-app-server/api/cmdserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
)

var ensureDefaultCmd = &cobra.Command{
	Use:   "ensure-default [inspect-ns|cleanup-ns <NETWORK_SERVER_ID>|organization]",
	Short: "connect to local command line service (:1000), inspect or manage internal services via command line",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Fatal("invalid argument")
		}

		conn, err := grpccli.Connect(grpccli.ConnectionOpts{
			Server:  "localhost:1000",
			CACert:  "",
			TLSCert: "",
			TLSKey:  "",
		})
		if err != nil {
			logrus.Fatalf("%v", err)
		}
		ctx := context.Background()
		client := pb.NewEnsureDefaultServiceClient(conn)
		switch args[0] {
		case "inspect-ns":
			resCli, err := client.InspectNetworkServerSettings(ctx, &pb.InspectNetworkServerSettingsRequest{})
			if err != nil {
				logrus.Fatal(err)
			}
			data := []byte("")
			for {
				res, err := resCli.Recv()
				if err != nil {
					logrus.Fatal(err)
				}
				data = append(data, res.Data...)
				if res.Finish {
					logrus.Println("Stream is over")
					// save fullData to file
					if err = ioutil.WriteFile("inspect_ns_report", data, os.ModePerm); err != nil {
						logrus.Fatalf("%v", err)
					}
					break
				}
			}
			return
		case "cleanup-ns":
			if len(args) != 2 {
				logrus.Fatal("invalid argument for cleanup-ns")
			}
			nsID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				logrus.Fatal("invalid network server id")
			}
			rescli, err := client.CorrectNetworkServerSettings(ctx, &pb.CorrectNetworkServerSettingsRequest{
				NetworkServerId: nsID})
			if err != nil {
				logrus.Fatal(err)
			}
			data := []byte("")
			for {
				res, err := rescli.Recv()
				if err != nil {
					logrus.Fatal(err)
				}
				data = append(data, res.Data...)
				if res.Finish {
					logrus.Println("Stream is over")
					// save fullData to file
					if err = ioutil.WriteFile("cleanup_ns_report", data, os.ModePerm); err != nil {
						logrus.Fatalf("%v", err)
					}
					break
				}
			}
			return
		case "organization":
		default:
			logrus.Fatal("invalid argument")
		}
	},
}
