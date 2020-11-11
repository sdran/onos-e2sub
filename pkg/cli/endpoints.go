// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cli

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	regapi "github.com/onosproject/onos-e2sub/api/e2/endpoint/v1beta1"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	registrationHeaders = "ID\tIP\tPort\n"
	endPointFormat      = "%s\t%s\t%d\n"
)

func displayHeaders(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, registrationHeaders)
}

func displayEndPoint(writer io.Writer, ep regapi.TerminationEndpoint) {
	_, _ = fmt.Fprintf(writer, endPointFormat, ep.ID, ep.ID, ep.Port)
}

func getListEndPointsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoints",
		Short: "Get endpoints",
		RunE:  runListEndpointsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListEndpointsCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		displayHeaders(writer)
		_ = writer.Flush()
	}

	request := regapi.ListTerminationsRequest{}

	client := regapi.NewE2RegistryServiceClient(conn)

	response, err := client.ListTerminations(context.Background(), &request)
	if err != nil {
		return err
	}

	for _, ep := range response.Endpoints {
		displayEndPoint(writer, ep)
	}

	_ = writer.Flush()

	return nil
}
