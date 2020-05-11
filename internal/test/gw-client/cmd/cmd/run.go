package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	hostname     string
	serverAddr   string
	serialNumber string
	gatewayMac   string
	model        string
	osVersion    string

	rootCAPath     string
	clientCertPath string
	clientKeyPath  string
)

// Execute prepares and starts the service
func Execute(host string) {
	hostname = host
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	RunE: run,
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&serverAddr, "server", "", "", "address of provisioning server")
	rootCmd.PersistentFlags().StringVarP(&serialNumber, "sn", "", "", "serial number of gateway")
	rootCmd.PersistentFlags().StringVarP(&gatewayMac, "mac", "", "", "mac address of gateway")
	rootCmd.PersistentFlags().StringVarP(&model, "model", "", "", "model of gateway")
	rootCmd.PersistentFlags().StringVarP(&osVersion, "os-version", "", "", "os version of gateway")


	rootCmd.PersistentFlags().StringVarP(&rootCAPath, "root-ca-path", "", "", "rootCA")
	rootCmd.PersistentFlags().StringVarP(&clientCertPath, "client-tls-certificate-path", "", "", "client TLS certificate")
	rootCmd.PersistentFlags().StringVarP(&clientKeyPath, "client-tls-key-path", "", "", "client TLS key")

	if serverAddr == "" {
		serverAddr = hostname
	}

	rootCmd.AddCommand(sendHeartbeatCmd)
}

func run(cmd *cobra.Command, args []string) error {
	println("Use help option to know more.")
	return nil
}
