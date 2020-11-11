// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cli

import (
	"github.com/onosproject/onos-lib-go/pkg/cli"
	loglib "github.com/onosproject/onos-lib-go/pkg/logging/cli"
	"github.com/spf13/cobra"
)

const (
	configName     = "e2sub"
	defaultAddress = "onos-e2sub:5150"
)

// init initializes the command line
func init() {
	cli.InitConfig(configName)
}

// Init is a hook called after cobra initialization
func Init() {
	// noop for now
}

// GetCommand returns the root command for the RAN service
func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "e2sub {get} [args]",
		Short: "ONOS e2sub subsystem commands",
	}

	cli.AddConfigFlags(cmd, defaultAddress)
	cmd.AddCommand(cli.GetConfigCommand())
	cmd.AddCommand(getGetCommand())
	cmd.AddCommand(loglib.GetCommand())
	return cmd
}