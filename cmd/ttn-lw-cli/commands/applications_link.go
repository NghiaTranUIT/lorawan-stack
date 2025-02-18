// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"os"
	"strings"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.thethings.network/lorawan-stack/v3/cmd/internal/io"
	"go.thethings.network/lorawan-stack/v3/cmd/ttn-lw-cli/internal/api"
	"go.thethings.network/lorawan-stack/v3/cmd/ttn-lw-cli/internal/util"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	selectApplicationLinkFlags = util.FieldMaskFlags(&ttnpb.ApplicationLink{})
	setApplicationLinkFlags    = util.FieldFlags(&ttnpb.ApplicationLink{})

	selectAllApplicationLinkFlags = util.SelectAllFlagSet("application link")
)

func deprecatedApplicationLinkFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.String("api-key", "", "")
	flagSet.Lookup("api-key").Hidden = true
	flagSet.String("network-server-address", "", "")
	flagSet.Lookup("network-server-address").Hidden = true
	return flagSet
}

var (
	applicationsLinkCommand = &cobra.Command{
		Use:   "link",
		Short: "Application link commands",
	}
	applicationsLinkGetCommand = &cobra.Command{
		Use:     "get [application-id]",
		Aliases: []string{"info"},
		Short:   "Get the properties of an application link",
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID.New()
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectApplicationLinkFlags)
			if len(paths) == 0 {
				logger.Warn("No fields selected, will select everything")
				selectApplicationLinkFlags.VisitAll(func(flag *pflag.Flag) {
					paths = append(paths, strings.Replace(flag.Name, "-", "_", -1))
				})
			}
			paths = ttnpb.AllowedFields(paths, ttnpb.RPCFieldMaskPaths["/ttn.lorawan.v3.As/GetLink"].Allowed)

			as, err := api.Dial(ctx, config.ApplicationServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewAsClient(as).GetLink(ctx, &ttnpb.GetApplicationLinkRequest{
				ApplicationIds: appID,
				FieldMask:      &pbtypes.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	applicationsLinkSetCommand = &cobra.Command{
		Use:     "set [application-id]",
		Aliases: []string{"update"},
		Short:   "Set the properties of an application link",
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID.New()
			}
			paths := util.UpdateFieldMask(cmd.Flags(), setApplicationLinkFlags)

			link := &ttnpb.ApplicationLink{}
			if err := util.SetFields(link, setApplicationLinkFlags); err != nil {
				return err
			}
			newPaths, err := parsePayloadFormatterParameterFlags("default-formatters", link.DefaultFormatters, cmd.Flags())
			if err != nil {
				return err
			}
			paths = append(paths, newPaths...)
			as, err := api.Dial(ctx, config.ApplicationServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewAsClient(as).SetLink(ctx, &ttnpb.SetApplicationLinkRequest{
				ApplicationIds: appID,
				Link:           link,
				FieldMask:      &pbtypes.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	applicationsLinkDeleteCommand = &cobra.Command{
		Use:     "delete [application-id]",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete an application link",
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := getApplicationID(cmd.Flags(), args)
			if appID == nil {
				return errNoApplicationID.New()
			}

			as, err := api.Dial(ctx, config.ApplicationServerGRPCAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewAsClient(as).DeleteLink(ctx, appID)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	applicationsLinkGetCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsLinkGetCommand.Flags().AddFlagSet(selectApplicationLinkFlags)
	applicationsLinkGetCommand.Flags().AddFlagSet(selectAllApplicationLinkFlags)
	applicationsLinkCommand.AddCommand(applicationsLinkGetCommand)
	applicationsLinkSetCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsLinkSetCommand.Flags().AddFlagSet(setApplicationLinkFlags)
	applicationsLinkSetCommand.Flags().AddFlagSet(payloadFormatterParameterFlags("default-formatters"))
	applicationsLinkSetCommand.Flags().AddFlagSet(deprecatedApplicationLinkFlags())
	applicationsLinkCommand.AddCommand(applicationsLinkSetCommand)
	applicationsLinkDeleteCommand.Flags().AddFlagSet(applicationIDFlags())
	applicationsLinkCommand.AddCommand(applicationsLinkDeleteCommand)
	applicationsCommand.AddCommand(applicationsLinkCommand)
}
