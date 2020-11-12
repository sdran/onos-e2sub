// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cli

import "github.com/spf13/cobra"

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list {endpoints} [args]",
		Short: "List E2T resources",
	}
	cmd.AddCommand(getListEndPointsCommand())
	return cmd
}

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {endpoint} [args]",
		Short: "Get E2T resources",
	}
	cmd.AddCommand(getGetEndPointCommand())
	return cmd
}

func getAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add {endpoint} [args]",
		Short: "Add E2T resources",
	}
	cmd.AddCommand(getAddEndPointCommand())
	return cmd
}

func getRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove {endpoint} [args]",
		Short: "Remove E2T resources",
	}
	cmd.AddCommand(getRemoveEndPointCommand())
	return cmd
}
