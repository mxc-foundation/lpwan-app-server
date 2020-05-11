package main

import (
	"github.com/mxc-foundation/lpwan-app-server/gw-client/cmd/cmd"
)

var host string

func main() {
	cmd.Execute(host)
}