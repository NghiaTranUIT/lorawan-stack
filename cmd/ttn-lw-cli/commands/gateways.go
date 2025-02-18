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

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.thethings.network/lorawan-stack/v3/cmd/internal/io"
	"go.thethings.network/lorawan-stack/v3/cmd/ttn-lw-cli/internal/api"
	"go.thethings.network/lorawan-stack/v3/cmd/ttn-lw-cli/internal/util"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	ttntypes "go.thethings.network/lorawan-stack/v3/pkg/types"
	"google.golang.org/grpc"
)

var (
	selectGatewayFlags     = util.FieldMaskFlags(&ttnpb.Gateway{})
	setGatewayFlags        = util.FieldFlags(&ttnpb.Gateway{})
	setGatewayAntennaFlags = util.FieldFlags(&ttnpb.GatewayAntenna{}, "antenna")
	selectAllGatewayFlags  = util.SelectAllFlagSet("gateway")

	gatewayFlattenPaths = []string{"lbs_lns_secret", "claim_authentication_code", "target_cups_key"}
)

func gatewayIDFlags() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.String("gateway-id", "", "")
	flagSet.String("gateway-eui", "", "")
	return flagSet
}

var (
	errNoGatewayID       = errors.DefineInvalidArgument("no_gateway_id", "no gateway ID set")
	errNoGatewayEUI      = errors.DefineInvalidArgument("no_gateway_eui", "no gateway EUI set")
	errInvalidGatewayEUI = errors.DefineInvalidArgument("invalid_gateway_eui", "invalid gateway EUI")
)

func getGatewayID(flagSet *pflag.FlagSet, args []string, requireID bool) (*ttnpb.GatewayIdentifiers, error) {
	gatewayID, _ := flagSet.GetString("gateway-id")
	gatewayEUIHex, _ := flagSet.GetString("gateway-eui")
	switch len(args) {
	case 0:
	case 1:
		gatewayID = args[0]
	case 2:
		gatewayID = args[0]
		gatewayEUIHex = args[1]
	default:
		logger.Warn("Multiple IDs found in arguments, considering the first")
		gatewayID = args[0]
		gatewayEUIHex = args[1]
	}
	if gatewayID == "" && requireID {
		return nil, errNoGatewayID.New()
	}
	ids := &ttnpb.GatewayIdentifiers{GatewayId: gatewayID}
	if gatewayEUIHex != "" {
		var gatewayEUI ttntypes.EUI64
		if err := gatewayEUI.UnmarshalText([]byte(gatewayEUIHex)); err != nil {
			return nil, errInvalidGatewayEUI.WithCause(err)
		}
		ids.Eui = &gatewayEUI
	}
	return ids, nil
}

func getGatewayEUI(flagSet *pflag.FlagSet, args []string, requireEUI bool) (*ttnpb.GatewayIdentifiers, error) {
	gatewayEUIHex, _ := flagSet.GetString("gateway-eui")
	switch len(args) {
	case 0:
	case 1:
		gatewayEUIHex = args[0]
	default:
		logger.Warn("Multiple EUIs found in arguments, considering the first")
		gatewayEUIHex = args[0]
	}
	if gatewayEUIHex == "" && requireEUI {
		return nil, errNoGatewayEUI.New()
	}
	ids := &ttnpb.GatewayIdentifiers{}
	if gatewayEUIHex != "" {
		var gatewayEUI ttntypes.EUI64
		if err := gatewayEUI.UnmarshalText([]byte(gatewayEUIHex)); err != nil {
			return nil, errInvalidGatewayEUI.WithCause(err)
		}
		ids.Eui = &gatewayEUI
	}
	return ids, nil
}

var searchGatewaysFlags = func() *pflag.FlagSet {
	flagSet := &pflag.FlagSet{}
	flagSet.AddFlagSet(searchFlags)
	// NOTE: These flags need to be named with underscores, not dashes!
	flagSet.String("eui_contains", "", "")
	return flagSet
}()

var (
	gatewaysCommand = &cobra.Command{
		Use:     "gateways",
		Aliases: []string{"gateway", "gtw", "g"},
		Short:   "Gateway commands",
	}
	gatewaysListFrequencyPlans = &cobra.Command{
		Use:               "list-frequency-plans",
		Aliases:           []string{"get-frequency-plans", "frequency-plans", "fps"},
		Short:             "List available frequency plans for gateways",
		PersistentPreRunE: preRun(),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseFrequency, _ := cmd.Flags().GetUint32("base-frequency")
			gs, err := api.Dial(ctx, config.GatewayServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewConfigurationClient(gs).ListFrequencyPlans(ctx, &ttnpb.ListFrequencyPlansRequest{
				BaseFrequency: baseFrequency,
			})
			if err != nil {
				return err
			}
			return io.Write(os.Stdout, config.OutputFormat, res.FrequencyPlans)
		},
	}
	gatewaysListCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List gateways",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := util.SelectFieldMask(cmd.Flags(), selectGatewayFlags)
			paths = ttnpb.AllowedFields(paths, ttnpb.RPCFieldMaskPaths["/ttn.lorawan.v3.GatewayRegistry/List"].Allowed)

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			limit, page, opt, getTotal := withPagination(cmd.Flags())
			res, err := ttnpb.NewGatewayRegistryClient(is).List(ctx, &ttnpb.ListGatewaysRequest{
				Collaborator: getCollaborator(cmd.Flags()),
				FieldMask:    &pbtypes.FieldMask{Paths: paths},
				Limit:        limit,
				Page:         page,
				Order:        getOrder(cmd.Flags()),
				Deleted:      getDeleted(cmd.Flags()),
			}, opt)
			if err != nil {
				return err
			}
			getTotal()

			return io.Write(os.Stdout, config.OutputFormat, res.Gateways)
		},
	}
	gatewaysSearchCommand = &cobra.Command{
		Use:   "search",
		Short: "Search for gateways",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := util.SelectFieldMask(cmd.Flags(), selectGatewayFlags)
			paths = ttnpb.AllowedFields(paths, ttnpb.RPCFieldMaskPaths["/ttn.lorawan.v3.EntityRegistrySearch/SearchGateways"].Allowed)

			req := &ttnpb.SearchGatewaysRequest{}
			if err := util.SetFields(req, searchGatewaysFlags); err != nil {
				return err
			}
			var (
				opt      grpc.CallOption
				getTotal func() uint64
			)
			req.Limit, req.Page, opt, getTotal = withPagination(cmd.Flags())
			req.FieldMask = &pbtypes.FieldMask{Paths: paths}
			req.Deleted = getDeleted(cmd.Flags())

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewEntityRegistrySearchClient(is).SearchGateways(ctx, req, opt)
			if err != nil {
				return err
			}
			getTotal()

			return io.Write(os.Stdout, config.OutputFormat, res.Gateways)
		},
	}
	gatewaysGetCommand = &cobra.Command{
		Use:     "get [gateway-id]",
		Aliases: []string{"info"},
		Short:   "Get a gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			gtwID, err := getGatewayID(cmd.Flags(), args, false)
			if err != nil {
				return err
			}
			paths := util.SelectFieldMask(cmd.Flags(), selectGatewayFlags)
			paths = ttnpb.AllowedFields(paths, ttnpb.RPCFieldMaskPaths["/ttn.lorawan.v3.GatewayRegistry/Get"].Allowed)

			paths = append(paths, ttnpb.FlattenPaths(paths, gatewayFlattenPaths)...)

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}

			cli := ttnpb.NewGatewayRegistryClient(is)

			if gtwID.GatewayId == "" && gtwID.Eui != nil {
				gtwID, err = cli.GetIdentifiersForEUI(ctx, &ttnpb.GetGatewayIdentifiersForEUIRequest{
					Eui: gtwID.Eui,
				})
				if err != nil {
					return err
				}
			}

			res, err := cli.Get(ctx, &ttnpb.GetGatewayRequest{
				GatewayIds: gtwID,
				FieldMask:  &pbtypes.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	gatewaysCreateCommand = &cobra.Command{
		Use:     "create [gateway-id]",
		Aliases: []string{"add", "register"},
		Short:   "Create a gateway",
		RunE: asBulk(func(cmd *cobra.Command, args []string) (err error) {
			gtwID, err := getGatewayID(cmd.Flags(), args, false)
			if err != nil {
				return err
			}

			collaborator := getCollaborator(cmd.Flags())
			if collaborator == nil {
				return errNoCollaborator.New()
			}
			var gateway ttnpb.Gateway
			if inputDecoder != nil {
				_, err := inputDecoder.Decode(&gateway)
				if err != nil {
					return err
				}
			}

			if setDefaults, _ := cmd.Flags().GetBool("defaults"); setDefaults {
				gateway.GatewayServerAddress = getHost(config.GatewayServerGRPCAddress)
				gateway.AutoUpdate = true
				gateway.EnforceDutyCycle = true
				gateway.StatusPublic = true
				gateway.LocationPublic = true
			}

			if err = util.SetFields(&gateway, setGatewayFlags); err != nil {
				return err
			}

			gateway.Attributes = mergeAttributes(gateway.Attributes, cmd.Flags())

			if gateway.Ids == nil {
				gateway.Ids = &ttnpb.GatewayIdentifiers{}
			}
			if gtwID != nil {
				if gtwID.GatewayId != "" {
					gateway.Ids.GatewayId = gtwID.GatewayId
				}
				if gtwID.Eui != nil {
					gateway.Ids.Eui = gtwID.Eui
				}
			}
			if gateway.Ids.GatewayId == "" {
				return errNoGatewayID.New()
			}

			var antenna *ttnpb.GatewayAntenna
			if err = util.SetFields(antenna, setGatewayAntennaFlags, "antenna"); err != nil {
				return err
			}
			if antenna != nil {
				gateway.Antennas = []*ttnpb.GatewayAntenna{antenna}
			}
			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			res, err := ttnpb.NewGatewayRegistryClient(is).Create(ctx, &ttnpb.CreateGatewayRequest{
				Gateway:      &gateway,
				Collaborator: collaborator,
			})
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		}),
	}
	errAntennaIndex    = errors.DefineInvalidArgument("antenna_index", "index of antenna to update out of bounds")
	gatewaysSetCommand = &cobra.Command{
		Use:     "set [gateway-id]",
		Aliases: []string{"update"},
		Short:   "Set properties of a gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			gtwID, err := getGatewayID(cmd.Flags(), args, true)
			if err != nil {
				return err
			}
			paths := util.UpdateFieldMask(cmd.Flags(), setGatewayFlags, attributesFlags())
			antennaPaths := util.UpdateFieldMask(cmd.Flags(), setGatewayAntennaFlags)
			paths = append(paths, ttnpb.FlattenPaths(paths, gatewayFlattenPaths)...)

			if gtwID.Eui != nil {
				paths = append(paths, "ids.eui")
			}
			antennaAdd, _ := cmd.Flags().GetBool("antenna.add")
			antennaRemove, _ := cmd.Flags().GetBool("antenna.remove")
			if len(paths)+len(antennaPaths) == 0 && !antennaRemove {
				logger.Warn("No fields selected, won't update anything")
				return nil
			}

			var gateway ttnpb.Gateway
			if err = util.SetFields(&gateway, setGatewayFlags); err != nil {
				return err
			}
			gateway.Attributes = mergeAttributes(gateway.Attributes, cmd.Flags())
			gateway.Ids = gtwID

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}

			if len(antennaPaths) > 0 || antennaAdd || antennaRemove {
				res, err := ttnpb.NewGatewayRegistryClient(is).Get(ctx, &ttnpb.GetGatewayRequest{
					GatewayIds: gateway.GetIds(),
					FieldMask:  &pbtypes.FieldMask{Paths: []string{"antennas"}},
				})
				if err != nil {
					return err
				}
				antennaIndex, _ := cmd.Flags().GetInt("antenna.index")
				if antennaAdd || len(res.Antennas) == 0 {
					res.Antennas = append(res.Antennas, &ttnpb.GatewayAntenna{})
					antennaIndex = len(res.Antennas) - 1
				} else if antennaIndex > len(res.Antennas) {
					return errAntennaIndex.New()
				}
				if antennaRemove {
					gateway.Antennas = append(res.Antennas[:antennaIndex], res.Antennas[antennaIndex+1:]...)
				} else { // create or update
					if err = util.SetFields(&res.Antennas[antennaIndex], setGatewayAntennaFlags, "antenna"); err != nil {
						return err
					}
					gateway.Antennas = res.Antennas
				}
				paths = append(paths, "antennas")
			}

			res, err := ttnpb.NewGatewayRegistryClient(is).Update(ctx, &ttnpb.UpdateGatewayRequest{
				Gateway:   &gateway,
				FieldMask: &pbtypes.FieldMask{Paths: paths},
			})
			if err != nil {
				return err
			}

			res.SetFields(&gateway, "ids")
			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	gatewaysDeleteCommand = &cobra.Command{
		Use:     "delete [gateway-id]",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "Delete a gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			gtwID, err := getGatewayID(cmd.Flags(), args, true)
			if err != nil {
				return err
			}

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewGatewayRegistryClient(is).Delete(ctx, gtwID)
			if err != nil {
				return err
			}

			return nil
		},
	}
	gatewaysRestoreCommand = &cobra.Command{
		Use:   "restore [gateway-id]",
		Short: "Restore a gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			gtwID, err := getGatewayID(cmd.Flags(), args, true)
			if err != nil {
				return err
			}

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewGatewayRegistryClient(is).Restore(ctx, gtwID)
			if err != nil {
				return err
			}

			return nil
		},
	}
	gatewaysConnectionStats = &cobra.Command{
		Use:     "get-connection-stats [gateway-id]",
		Aliases: []string{"connection-stats", "cnx-stats", "stats"},
		Short:   "Get connection stats for a gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			gtwID, err := getGatewayID(cmd.Flags(), args, true)
			if err != nil {
				return err
			}

			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}

			gateway, err := ttnpb.NewGatewayRegistryClient(is).Get(ctx, &ttnpb.GetGatewayRequest{
				GatewayIds: gtwID,
				FieldMask:  &pbtypes.FieldMask{Paths: []string{"gateway_server_address"}},
			})
			if err != nil {
				return err
			}

			if gsMismatch := compareServerAddressGateway(gateway, config); gsMismatch {
				return errAddressMismatchGateway.New()
			}

			gs, err := api.Dial(ctx, config.GatewayServerGRPCAddress)
			if err != nil {
				return err
			}

			res, err := ttnpb.NewGsClient(gs).GetGatewayConnectionStats(ctx, gtwID)
			if err != nil {
				return err
			}

			return io.Write(os.Stdout, config.OutputFormat, res)
		},
	}
	gatewaysContactInfoCommand = contactInfoCommands("gateway", func(cmd *cobra.Command, args []string) (*ttnpb.EntityIdentifiers, error) {
		gtwID, err := getGatewayID(cmd.Flags(), args, true)
		if err != nil {
			return nil, err
		}
		return gtwID.GetEntityIdentifiers(), nil
	})
	gatewaysPurgeCommand = &cobra.Command{
		Use:     "purge [gateway-id]",
		Aliases: []string{"permanent-delete", "hard-delete"},
		Short:   "Purge a gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			gtwID, err := getGatewayID(cmd.Flags(), args, true)
			if err != nil {
				return err
			}

			force, err := cmd.Flags().GetBool("force")
			if err != nil {
				return err
			}
			if !confirmChoice(gatewayPurgeWarning, force) {
				return errNoConfirmation.New()
			}
			is, err := api.Dial(ctx, config.IdentityServerGRPCAddress)
			if err != nil {
				return err
			}
			_, err = ttnpb.NewGatewayRegistryClient(is).Purge(ctx, gtwID)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	gatewaysListFrequencyPlans.Flags().Uint32("base-frequency", 0, "Base frequency in MHz for hardware support (433, 470, 868 or 915)")
	gatewaysCommand.AddCommand(gatewaysListFrequencyPlans)
	gatewaysListCommand.Flags().AddFlagSet(collaboratorFlags())
	gatewaysListCommand.Flags().AddFlagSet(deletedFlags)
	gatewaysListCommand.Flags().AddFlagSet(selectGatewayFlags)
	gatewaysListCommand.Flags().AddFlagSet(paginationFlags())
	gatewaysListCommand.Flags().AddFlagSet(orderFlags())
	gatewaysListCommand.Flags().AddFlagSet(selectAllGatewayFlags)
	gatewaysCommand.AddCommand(gatewaysListCommand)
	gatewaysSearchCommand.Flags().AddFlagSet(searchGatewaysFlags)
	gatewaysSearchCommand.Flags().AddFlagSet(deletedFlags)
	gatewaysSearchCommand.Flags().AddFlagSet(selectGatewayFlags)
	gatewaysSearchCommand.Flags().AddFlagSet(selectAllGatewayFlags)
	gatewaysCommand.AddCommand(gatewaysSearchCommand)
	gatewaysGetCommand.Flags().AddFlagSet(gatewayIDFlags())
	gatewaysGetCommand.Flags().AddFlagSet(selectGatewayFlags)
	gatewaysGetCommand.Flags().AddFlagSet(selectAllGatewayFlags)
	gatewaysCommand.AddCommand(gatewaysGetCommand)
	gatewaysCreateCommand.Flags().AddFlagSet(gatewayIDFlags())
	gatewaysCreateCommand.Flags().AddFlagSet(collaboratorFlags())
	gatewaysCreateCommand.Flags().AddFlagSet(setGatewayFlags)
	gatewaysCreateCommand.Flags().AddFlagSet(setGatewayAntennaFlags)
	gatewaysCreateCommand.Flags().AddFlagSet(attributesFlags())
	gatewaysCreateCommand.Flags().Bool("defaults", true, "configure gateway with defaults")
	gatewaysCommand.AddCommand(gatewaysCreateCommand)
	gatewaysSetCommand.Flags().AddFlagSet(gatewayIDFlags())
	gatewaysSetCommand.Flags().AddFlagSet(setGatewayFlags)
	gatewaysSetCommand.Flags().Int("antenna.index", 0, "index of the antenna to update or remove")
	gatewaysSetCommand.Flags().Bool("antenna.add", false, "add an extra antenna")
	gatewaysSetCommand.Flags().Bool("antenna.remove", false, "remove an antenna")
	gatewaysSetCommand.Flags().AddFlagSet(setGatewayAntennaFlags)
	gatewaysSetCommand.Flags().AddFlagSet(attributesFlags())
	gatewaysCommand.AddCommand(gatewaysSetCommand)
	gatewaysDeleteCommand.Flags().AddFlagSet(gatewayIDFlags())
	gatewaysCommand.AddCommand(gatewaysDeleteCommand)
	gatewaysRestoreCommand.Flags().AddFlagSet(gatewayIDFlags())
	gatewaysCommand.AddCommand(gatewaysRestoreCommand)
	gatewaysConnectionStats.Flags().AddFlagSet(gatewayIDFlags())
	gatewaysCommand.AddCommand(gatewaysConnectionStats)
	gatewaysContactInfoCommand.PersistentFlags().AddFlagSet(gatewayIDFlags())
	gatewaysCommand.AddCommand(gatewaysContactInfoCommand)
	gatewaysPurgeCommand.Flags().AddFlagSet(gatewayIDFlags())
	gatewaysPurgeCommand.Flags().AddFlagSet(forceFlags())
	gatewaysCommand.AddCommand(gatewaysPurgeCommand)
	Root.AddCommand(gatewaysCommand)
}

var errAddressMismatchGateway = errors.DefineAborted("gateway_server_address_mismatch", "gateway server address mismatch")

func compareServerAddressGateway(gateway *ttnpb.Gateway, config *Config) (gsMismatch bool) {
	gsHost := getHost(config.GatewayServerGRPCAddress)
	if host := getHost(gateway.GatewayServerAddress); host != "" && host != gsHost {
		gsMismatch = true
		logger.WithFields(log.Fields(
			"configured", gsHost,
			"registered", host,
		)).Warn("Registered Gateway Server address does not match CLI configuration")
	}
	return
}
