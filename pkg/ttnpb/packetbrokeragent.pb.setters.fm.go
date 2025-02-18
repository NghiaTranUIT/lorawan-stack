// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import fmt "fmt"

func (dst *PacketBrokerGateway) SetFields(src *PacketBrokerGateway, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "ids":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerGateway_GatewayIdentifiers
				if (src == nil || src.Ids == nil) && dst.Ids == nil {
					continue
				}
				if src != nil {
					newSrc = src.Ids
				}
				if dst.Ids != nil {
					newDst = dst.Ids
				} else {
					newDst = &PacketBrokerGateway_GatewayIdentifiers{}
					dst.Ids = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Ids = src.Ids
				} else {
					dst.Ids = nil
				}
			}
		case "contact_info":
			if len(subs) > 0 {
				return fmt.Errorf("'contact_info' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ContactInfo = src.ContactInfo
			} else {
				dst.ContactInfo = nil
			}
		case "antennas":
			if len(subs) > 0 {
				return fmt.Errorf("'antennas' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Antennas = src.Antennas
			} else {
				dst.Antennas = nil
			}
		case "status_public":
			if len(subs) > 0 {
				return fmt.Errorf("'status_public' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.StatusPublic = src.StatusPublic
			} else {
				var zero bool
				dst.StatusPublic = zero
			}
		case "location_public":
			if len(subs) > 0 {
				return fmt.Errorf("'location_public' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.LocationPublic = src.LocationPublic
			} else {
				var zero bool
				dst.LocationPublic = zero
			}
		case "frequency_plan_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'frequency_plan_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FrequencyPlanIds = src.FrequencyPlanIds
			} else {
				dst.FrequencyPlanIds = nil
			}
		case "update_location_from_status":
			if len(subs) > 0 {
				return fmt.Errorf("'update_location_from_status' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpdateLocationFromStatus = src.UpdateLocationFromStatus
			} else {
				var zero bool
				dst.UpdateLocationFromStatus = zero
			}
		case "online":
			if len(subs) > 0 {
				return fmt.Errorf("'online' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Online = src.Online
			} else {
				var zero bool
				dst.Online = zero
			}
		case "rx_rate":
			if len(subs) > 0 {
				return fmt.Errorf("'rx_rate' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RxRate = src.RxRate
			} else {
				dst.RxRate = nil
			}
		case "tx_rate":
			if len(subs) > 0 {
				return fmt.Errorf("'tx_rate' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TxRate = src.TxRate
			} else {
				dst.TxRate = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *UpdatePacketBrokerGatewayRequest) SetFields(src *UpdatePacketBrokerGatewayRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "gateway":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerGateway
				if (src == nil || src.Gateway == nil) && dst.Gateway == nil {
					continue
				}
				if src != nil {
					newSrc = src.Gateway
				}
				if dst.Gateway != nil {
					newDst = dst.Gateway
				} else {
					newDst = &PacketBrokerGateway{}
					dst.Gateway = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Gateway = src.Gateway
				} else {
					dst.Gateway = nil
				}
			}
		case "field_mask":
			if len(subs) > 0 {
				return fmt.Errorf("'field_mask' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FieldMask = src.FieldMask
			} else {
				dst.FieldMask = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *UpdatePacketBrokerGatewayResponse) SetFields(src *UpdatePacketBrokerGatewayResponse, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "online_ttl":
			if len(subs) > 0 {
				return fmt.Errorf("'online_ttl' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.OnlineTtl = src.OnlineTtl
			} else {
				dst.OnlineTtl = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerNetworkIdentifier) SetFields(src *PacketBrokerNetworkIdentifier, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "net_id":
			if len(subs) > 0 {
				return fmt.Errorf("'net_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.NetId = src.NetId
			} else {
				var zero uint32
				dst.NetId = zero
			}
		case "tenant_id":
			if len(subs) > 0 {
				return fmt.Errorf("'tenant_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TenantId = src.TenantId
			} else {
				var zero string
				dst.TenantId = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerDevAddrBlock) SetFields(src *PacketBrokerDevAddrBlock, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "dev_addr_prefix":
			if len(subs) > 0 {
				var newDst, newSrc *DevAddrPrefix
				if (src == nil || src.DevAddrPrefix == nil) && dst.DevAddrPrefix == nil {
					continue
				}
				if src != nil {
					newSrc = src.DevAddrPrefix
				}
				if dst.DevAddrPrefix != nil {
					newDst = dst.DevAddrPrefix
				} else {
					newDst = &DevAddrPrefix{}
					dst.DevAddrPrefix = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.DevAddrPrefix = src.DevAddrPrefix
				} else {
					dst.DevAddrPrefix = nil
				}
			}
		case "home_network_cluster_id":
			if len(subs) > 0 {
				return fmt.Errorf("'home_network_cluster_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.HomeNetworkClusterId = src.HomeNetworkClusterId
			} else {
				var zero string
				dst.HomeNetworkClusterId = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerNetwork) SetFields(src *PacketBrokerNetwork, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "id":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerNetworkIdentifier
				if (src == nil || src.Id == nil) && dst.Id == nil {
					continue
				}
				if src != nil {
					newSrc = src.Id
				}
				if dst.Id != nil {
					newDst = dst.Id
				} else {
					newDst = &PacketBrokerNetworkIdentifier{}
					dst.Id = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Id = src.Id
				} else {
					dst.Id = nil
				}
			}
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}
		case "dev_addr_blocks":
			if len(subs) > 0 {
				return fmt.Errorf("'dev_addr_blocks' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DevAddrBlocks = src.DevAddrBlocks
			} else {
				dst.DevAddrBlocks = nil
			}
		case "contact_info":
			if len(subs) > 0 {
				return fmt.Errorf("'contact_info' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ContactInfo = src.ContactInfo
			} else {
				dst.ContactInfo = nil
			}
		case "listed":
			if len(subs) > 0 {
				return fmt.Errorf("'listed' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Listed = src.Listed
			} else {
				var zero bool
				dst.Listed = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerNetworks) SetFields(src *PacketBrokerNetworks, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "networks":
			if len(subs) > 0 {
				return fmt.Errorf("'networks' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Networks = src.Networks
			} else {
				dst.Networks = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerInfo) SetFields(src *PacketBrokerInfo, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "registration":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerNetwork
				if (src == nil || src.Registration == nil) && dst.Registration == nil {
					continue
				}
				if src != nil {
					newSrc = src.Registration
				}
				if dst.Registration != nil {
					newDst = dst.Registration
				} else {
					newDst = &PacketBrokerNetwork{}
					dst.Registration = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Registration = src.Registration
				} else {
					dst.Registration = nil
				}
			}
		case "forwarder_enabled":
			if len(subs) > 0 {
				return fmt.Errorf("'forwarder_enabled' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ForwarderEnabled = src.ForwarderEnabled
			} else {
				var zero bool
				dst.ForwarderEnabled = zero
			}
		case "home_network_enabled":
			if len(subs) > 0 {
				return fmt.Errorf("'home_network_enabled' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.HomeNetworkEnabled = src.HomeNetworkEnabled
			} else {
				var zero bool
				dst.HomeNetworkEnabled = zero
			}
		case "register_enabled":
			if len(subs) > 0 {
				return fmt.Errorf("'register_enabled' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.RegisterEnabled = src.RegisterEnabled
			} else {
				var zero bool
				dst.RegisterEnabled = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerRegisterRequest) SetFields(src *PacketBrokerRegisterRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "listed":
			if len(subs) > 0 {
				return fmt.Errorf("'listed' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Listed = src.Listed
			} else {
				dst.Listed = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerRoutingPolicyUplink) SetFields(src *PacketBrokerRoutingPolicyUplink, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "join_request":
			if len(subs) > 0 {
				return fmt.Errorf("'join_request' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.JoinRequest = src.JoinRequest
			} else {
				var zero bool
				dst.JoinRequest = zero
			}
		case "mac_data":
			if len(subs) > 0 {
				return fmt.Errorf("'mac_data' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MacData = src.MacData
			} else {
				var zero bool
				dst.MacData = zero
			}
		case "application_data":
			if len(subs) > 0 {
				return fmt.Errorf("'application_data' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ApplicationData = src.ApplicationData
			} else {
				var zero bool
				dst.ApplicationData = zero
			}
		case "signal_quality":
			if len(subs) > 0 {
				return fmt.Errorf("'signal_quality' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.SignalQuality = src.SignalQuality
			} else {
				var zero bool
				dst.SignalQuality = zero
			}
		case "localization":
			if len(subs) > 0 {
				return fmt.Errorf("'localization' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Localization = src.Localization
			} else {
				var zero bool
				dst.Localization = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerRoutingPolicyDownlink) SetFields(src *PacketBrokerRoutingPolicyDownlink, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "join_accept":
			if len(subs) > 0 {
				return fmt.Errorf("'join_accept' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.JoinAccept = src.JoinAccept
			} else {
				var zero bool
				dst.JoinAccept = zero
			}
		case "mac_data":
			if len(subs) > 0 {
				return fmt.Errorf("'mac_data' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.MacData = src.MacData
			} else {
				var zero bool
				dst.MacData = zero
			}
		case "application_data":
			if len(subs) > 0 {
				return fmt.Errorf("'application_data' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ApplicationData = src.ApplicationData
			} else {
				var zero bool
				dst.ApplicationData = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerDefaultRoutingPolicy) SetFields(src *PacketBrokerDefaultRoutingPolicy, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "updated_at":
			if len(subs) > 0 {
				return fmt.Errorf("'updated_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpdatedAt = src.UpdatedAt
			} else {
				dst.UpdatedAt = nil
			}
		case "uplink":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerRoutingPolicyUplink
				if (src == nil || src.Uplink == nil) && dst.Uplink == nil {
					continue
				}
				if src != nil {
					newSrc = src.Uplink
				}
				if dst.Uplink != nil {
					newDst = dst.Uplink
				} else {
					newDst = &PacketBrokerRoutingPolicyUplink{}
					dst.Uplink = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Uplink = src.Uplink
				} else {
					dst.Uplink = nil
				}
			}
		case "downlink":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerRoutingPolicyDownlink
				if (src == nil || src.Downlink == nil) && dst.Downlink == nil {
					continue
				}
				if src != nil {
					newSrc = src.Downlink
				}
				if dst.Downlink != nil {
					newDst = dst.Downlink
				} else {
					newDst = &PacketBrokerRoutingPolicyDownlink{}
					dst.Downlink = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Downlink = src.Downlink
				} else {
					dst.Downlink = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerRoutingPolicy) SetFields(src *PacketBrokerRoutingPolicy, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "forwarder_id":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerNetworkIdentifier
				if (src == nil || src.ForwarderId == nil) && dst.ForwarderId == nil {
					continue
				}
				if src != nil {
					newSrc = src.ForwarderId
				}
				if dst.ForwarderId != nil {
					newDst = dst.ForwarderId
				} else {
					newDst = &PacketBrokerNetworkIdentifier{}
					dst.ForwarderId = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ForwarderId = src.ForwarderId
				} else {
					dst.ForwarderId = nil
				}
			}
		case "home_network_id":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerNetworkIdentifier
				if (src == nil || src.HomeNetworkId == nil) && dst.HomeNetworkId == nil {
					continue
				}
				if src != nil {
					newSrc = src.HomeNetworkId
				}
				if dst.HomeNetworkId != nil {
					newDst = dst.HomeNetworkId
				} else {
					newDst = &PacketBrokerNetworkIdentifier{}
					dst.HomeNetworkId = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.HomeNetworkId = src.HomeNetworkId
				} else {
					dst.HomeNetworkId = nil
				}
			}
		case "updated_at":
			if len(subs) > 0 {
				return fmt.Errorf("'updated_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpdatedAt = src.UpdatedAt
			} else {
				dst.UpdatedAt = nil
			}
		case "uplink":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerRoutingPolicyUplink
				if (src == nil || src.Uplink == nil) && dst.Uplink == nil {
					continue
				}
				if src != nil {
					newSrc = src.Uplink
				}
				if dst.Uplink != nil {
					newDst = dst.Uplink
				} else {
					newDst = &PacketBrokerRoutingPolicyUplink{}
					dst.Uplink = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Uplink = src.Uplink
				} else {
					dst.Uplink = nil
				}
			}
		case "downlink":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerRoutingPolicyDownlink
				if (src == nil || src.Downlink == nil) && dst.Downlink == nil {
					continue
				}
				if src != nil {
					newSrc = src.Downlink
				}
				if dst.Downlink != nil {
					newDst = dst.Downlink
				} else {
					newDst = &PacketBrokerRoutingPolicyDownlink{}
					dst.Downlink = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Downlink = src.Downlink
				} else {
					dst.Downlink = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SetPacketBrokerDefaultRoutingPolicyRequest) SetFields(src *SetPacketBrokerDefaultRoutingPolicyRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "uplink":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerRoutingPolicyUplink
				if (src == nil || src.Uplink == nil) && dst.Uplink == nil {
					continue
				}
				if src != nil {
					newSrc = src.Uplink
				}
				if dst.Uplink != nil {
					newDst = dst.Uplink
				} else {
					newDst = &PacketBrokerRoutingPolicyUplink{}
					dst.Uplink = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Uplink = src.Uplink
				} else {
					dst.Uplink = nil
				}
			}
		case "downlink":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerRoutingPolicyDownlink
				if (src == nil || src.Downlink == nil) && dst.Downlink == nil {
					continue
				}
				if src != nil {
					newSrc = src.Downlink
				}
				if dst.Downlink != nil {
					newDst = dst.Downlink
				} else {
					newDst = &PacketBrokerRoutingPolicyDownlink{}
					dst.Downlink = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Downlink = src.Downlink
				} else {
					dst.Downlink = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ListHomeNetworkRoutingPoliciesRequest) SetFields(src *ListHomeNetworkRoutingPoliciesRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "limit":
			if len(subs) > 0 {
				return fmt.Errorf("'limit' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Limit = src.Limit
			} else {
				var zero uint32
				dst.Limit = zero
			}
		case "page":
			if len(subs) > 0 {
				return fmt.Errorf("'page' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Page = src.Page
			} else {
				var zero uint32
				dst.Page = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerRoutingPolicies) SetFields(src *PacketBrokerRoutingPolicies, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "policies":
			if len(subs) > 0 {
				return fmt.Errorf("'policies' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Policies = src.Policies
			} else {
				dst.Policies = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SetPacketBrokerRoutingPolicyRequest) SetFields(src *SetPacketBrokerRoutingPolicyRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "home_network_id":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerNetworkIdentifier
				if (src == nil || src.HomeNetworkId == nil) && dst.HomeNetworkId == nil {
					continue
				}
				if src != nil {
					newSrc = src.HomeNetworkId
				}
				if dst.HomeNetworkId != nil {
					newDst = dst.HomeNetworkId
				} else {
					newDst = &PacketBrokerNetworkIdentifier{}
					dst.HomeNetworkId = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.HomeNetworkId = src.HomeNetworkId
				} else {
					dst.HomeNetworkId = nil
				}
			}
		case "uplink":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerRoutingPolicyUplink
				if (src == nil || src.Uplink == nil) && dst.Uplink == nil {
					continue
				}
				if src != nil {
					newSrc = src.Uplink
				}
				if dst.Uplink != nil {
					newDst = dst.Uplink
				} else {
					newDst = &PacketBrokerRoutingPolicyUplink{}
					dst.Uplink = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Uplink = src.Uplink
				} else {
					dst.Uplink = nil
				}
			}
		case "downlink":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerRoutingPolicyDownlink
				if (src == nil || src.Downlink == nil) && dst.Downlink == nil {
					continue
				}
				if src != nil {
					newSrc = src.Downlink
				}
				if dst.Downlink != nil {
					newDst = dst.Downlink
				} else {
					newDst = &PacketBrokerRoutingPolicyDownlink{}
					dst.Downlink = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Downlink = src.Downlink
				} else {
					dst.Downlink = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerGatewayVisibility) SetFields(src *PacketBrokerGatewayVisibility, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "location":
			if len(subs) > 0 {
				return fmt.Errorf("'location' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Location = src.Location
			} else {
				var zero bool
				dst.Location = zero
			}
		case "antenna_placement":
			if len(subs) > 0 {
				return fmt.Errorf("'antenna_placement' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AntennaPlacement = src.AntennaPlacement
			} else {
				var zero bool
				dst.AntennaPlacement = zero
			}
		case "antenna_count":
			if len(subs) > 0 {
				return fmt.Errorf("'antenna_count' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AntennaCount = src.AntennaCount
			} else {
				var zero bool
				dst.AntennaCount = zero
			}
		case "fine_timestamps":
			if len(subs) > 0 {
				return fmt.Errorf("'fine_timestamps' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FineTimestamps = src.FineTimestamps
			} else {
				var zero bool
				dst.FineTimestamps = zero
			}
		case "contact_info":
			if len(subs) > 0 {
				return fmt.Errorf("'contact_info' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ContactInfo = src.ContactInfo
			} else {
				var zero bool
				dst.ContactInfo = zero
			}
		case "status":
			if len(subs) > 0 {
				return fmt.Errorf("'status' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Status = src.Status
			} else {
				var zero bool
				dst.Status = zero
			}
		case "frequency_plan":
			if len(subs) > 0 {
				return fmt.Errorf("'frequency_plan' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FrequencyPlan = src.FrequencyPlan
			} else {
				var zero bool
				dst.FrequencyPlan = zero
			}
		case "packet_rates":
			if len(subs) > 0 {
				return fmt.Errorf("'packet_rates' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.PacketRates = src.PacketRates
			} else {
				var zero bool
				dst.PacketRates = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerDefaultGatewayVisibility) SetFields(src *PacketBrokerDefaultGatewayVisibility, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "updated_at":
			if len(subs) > 0 {
				return fmt.Errorf("'updated_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpdatedAt = src.UpdatedAt
			} else {
				dst.UpdatedAt = nil
			}
		case "visibility":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerGatewayVisibility
				if (src == nil || src.Visibility == nil) && dst.Visibility == nil {
					continue
				}
				if src != nil {
					newSrc = src.Visibility
				}
				if dst.Visibility != nil {
					newDst = dst.Visibility
				} else {
					newDst = &PacketBrokerGatewayVisibility{}
					dst.Visibility = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Visibility = src.Visibility
				} else {
					dst.Visibility = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SetPacketBrokerDefaultGatewayVisibilityRequest) SetFields(src *SetPacketBrokerDefaultGatewayVisibilityRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "visibility":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerGatewayVisibility
				if (src == nil || src.Visibility == nil) && dst.Visibility == nil {
					continue
				}
				if src != nil {
					newSrc = src.Visibility
				}
				if dst.Visibility != nil {
					newDst = dst.Visibility
				} else {
					newDst = &PacketBrokerGatewayVisibility{}
					dst.Visibility = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Visibility = src.Visibility
				} else {
					dst.Visibility = nil
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ListPacketBrokerNetworksRequest) SetFields(src *ListPacketBrokerNetworksRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "limit":
			if len(subs) > 0 {
				return fmt.Errorf("'limit' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Limit = src.Limit
			} else {
				var zero uint32
				dst.Limit = zero
			}
		case "page":
			if len(subs) > 0 {
				return fmt.Errorf("'page' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Page = src.Page
			} else {
				var zero uint32
				dst.Page = zero
			}
		case "with_routing_policy":
			if len(subs) > 0 {
				return fmt.Errorf("'with_routing_policy' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.WithRoutingPolicy = src.WithRoutingPolicy
			} else {
				var zero bool
				dst.WithRoutingPolicy = zero
			}
		case "tenant_id_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'tenant_id_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TenantIdContains = src.TenantIdContains
			} else {
				var zero string
				dst.TenantIdContains = zero
			}
		case "name_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'name_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.NameContains = src.NameContains
			} else {
				var zero string
				dst.NameContains = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ListPacketBrokerHomeNetworksRequest) SetFields(src *ListPacketBrokerHomeNetworksRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "limit":
			if len(subs) > 0 {
				return fmt.Errorf("'limit' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Limit = src.Limit
			} else {
				var zero uint32
				dst.Limit = zero
			}
		case "page":
			if len(subs) > 0 {
				return fmt.Errorf("'page' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Page = src.Page
			} else {
				var zero uint32
				dst.Page = zero
			}
		case "tenant_id_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'tenant_id_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.TenantIdContains = src.TenantIdContains
			} else {
				var zero string
				dst.TenantIdContains = zero
			}
		case "name_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'name_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.NameContains = src.NameContains
			} else {
				var zero string
				dst.NameContains = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ListForwarderRoutingPoliciesRequest) SetFields(src *ListForwarderRoutingPoliciesRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "home_network_id":
			if len(subs) > 0 {
				var newDst, newSrc *PacketBrokerNetworkIdentifier
				if (src == nil || src.HomeNetworkId == nil) && dst.HomeNetworkId == nil {
					continue
				}
				if src != nil {
					newSrc = src.HomeNetworkId
				}
				if dst.HomeNetworkId != nil {
					newDst = dst.HomeNetworkId
				} else {
					newDst = &PacketBrokerNetworkIdentifier{}
					dst.HomeNetworkId = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.HomeNetworkId = src.HomeNetworkId
				} else {
					dst.HomeNetworkId = nil
				}
			}
		case "limit":
			if len(subs) > 0 {
				return fmt.Errorf("'limit' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Limit = src.Limit
			} else {
				var zero uint32
				dst.Limit = zero
			}
		case "page":
			if len(subs) > 0 {
				return fmt.Errorf("'page' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Page = src.Page
			} else {
				var zero uint32
				dst.Page = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *PacketBrokerGateway_GatewayIdentifiers) SetFields(src *PacketBrokerGateway_GatewayIdentifiers, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "gateway_id":
			if len(subs) > 0 {
				return fmt.Errorf("'gateway_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.GatewayId = src.GatewayId
			} else {
				var zero string
				dst.GatewayId = zero
			}
		case "eui":
			if len(subs) > 0 {
				return fmt.Errorf("'eui' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Eui = src.Eui
			} else {
				dst.Eui = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
