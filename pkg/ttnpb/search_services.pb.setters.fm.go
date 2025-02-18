// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import fmt "fmt"

func (dst *SearchApplicationsRequest) SetFields(src *SearchApplicationsRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "id_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'id_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.IdContains = src.IdContains
			} else {
				var zero string
				dst.IdContains = zero
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
		case "description_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'description_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DescriptionContains = src.DescriptionContains
			} else {
				var zero string
				dst.DescriptionContains = zero
			}
		case "attributes_contain":
			if len(subs) > 0 {
				return fmt.Errorf("'attributes_contain' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AttributesContain = src.AttributesContain
			} else {
				dst.AttributesContain = nil
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
		case "order":
			if len(subs) > 0 {
				return fmt.Errorf("'order' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Order = src.Order
			} else {
				var zero string
				dst.Order = zero
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
		case "deleted":
			if len(subs) > 0 {
				return fmt.Errorf("'deleted' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Deleted = src.Deleted
			} else {
				var zero bool
				dst.Deleted = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SearchClientsRequest) SetFields(src *SearchClientsRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "id_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'id_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.IdContains = src.IdContains
			} else {
				var zero string
				dst.IdContains = zero
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
		case "description_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'description_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DescriptionContains = src.DescriptionContains
			} else {
				var zero string
				dst.DescriptionContains = zero
			}
		case "attributes_contain":
			if len(subs) > 0 {
				return fmt.Errorf("'attributes_contain' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AttributesContain = src.AttributesContain
			} else {
				dst.AttributesContain = nil
			}
		case "state":
			if len(subs) > 0 {
				return fmt.Errorf("'state' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.State = src.State
			} else {
				dst.State = nil
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
		case "order":
			if len(subs) > 0 {
				return fmt.Errorf("'order' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Order = src.Order
			} else {
				var zero string
				dst.Order = zero
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
		case "deleted":
			if len(subs) > 0 {
				return fmt.Errorf("'deleted' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Deleted = src.Deleted
			} else {
				var zero bool
				dst.Deleted = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SearchGatewaysRequest) SetFields(src *SearchGatewaysRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "id_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'id_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.IdContains = src.IdContains
			} else {
				var zero string
				dst.IdContains = zero
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
		case "description_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'description_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DescriptionContains = src.DescriptionContains
			} else {
				var zero string
				dst.DescriptionContains = zero
			}
		case "attributes_contain":
			if len(subs) > 0 {
				return fmt.Errorf("'attributes_contain' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AttributesContain = src.AttributesContain
			} else {
				dst.AttributesContain = nil
			}
		case "eui_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'eui_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.EuiContains = src.EuiContains
			} else {
				var zero string
				dst.EuiContains = zero
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
		case "order":
			if len(subs) > 0 {
				return fmt.Errorf("'order' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Order = src.Order
			} else {
				var zero string
				dst.Order = zero
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
		case "deleted":
			if len(subs) > 0 {
				return fmt.Errorf("'deleted' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Deleted = src.Deleted
			} else {
				var zero bool
				dst.Deleted = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SearchOrganizationsRequest) SetFields(src *SearchOrganizationsRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "id_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'id_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.IdContains = src.IdContains
			} else {
				var zero string
				dst.IdContains = zero
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
		case "description_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'description_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DescriptionContains = src.DescriptionContains
			} else {
				var zero string
				dst.DescriptionContains = zero
			}
		case "attributes_contain":
			if len(subs) > 0 {
				return fmt.Errorf("'attributes_contain' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AttributesContain = src.AttributesContain
			} else {
				dst.AttributesContain = nil
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
		case "order":
			if len(subs) > 0 {
				return fmt.Errorf("'order' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Order = src.Order
			} else {
				var zero string
				dst.Order = zero
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
		case "deleted":
			if len(subs) > 0 {
				return fmt.Errorf("'deleted' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Deleted = src.Deleted
			} else {
				var zero bool
				dst.Deleted = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SearchUsersRequest) SetFields(src *SearchUsersRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "id_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'id_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.IdContains = src.IdContains
			} else {
				var zero string
				dst.IdContains = zero
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
		case "description_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'description_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DescriptionContains = src.DescriptionContains
			} else {
				var zero string
				dst.DescriptionContains = zero
			}
		case "attributes_contain":
			if len(subs) > 0 {
				return fmt.Errorf("'attributes_contain' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AttributesContain = src.AttributesContain
			} else {
				dst.AttributesContain = nil
			}
		case "state":
			if len(subs) > 0 {
				return fmt.Errorf("'state' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.State = src.State
			} else {
				dst.State = nil
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
		case "order":
			if len(subs) > 0 {
				return fmt.Errorf("'order' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Order = src.Order
			} else {
				var zero string
				dst.Order = zero
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
		case "deleted":
			if len(subs) > 0 {
				return fmt.Errorf("'deleted' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Deleted = src.Deleted
			} else {
				var zero bool
				dst.Deleted = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SearchEndDevicesRequest) SetFields(src *SearchEndDevicesRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "application_ids":
			if len(subs) > 0 {
				var newDst, newSrc *ApplicationIdentifiers
				if (src == nil || src.ApplicationIds == nil) && dst.ApplicationIds == nil {
					continue
				}
				if src != nil {
					newSrc = src.ApplicationIds
				}
				if dst.ApplicationIds != nil {
					newDst = dst.ApplicationIds
				} else {
					newDst = &ApplicationIdentifiers{}
					dst.ApplicationIds = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.ApplicationIds = src.ApplicationIds
				} else {
					dst.ApplicationIds = nil
				}
			}
		case "id_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'id_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.IdContains = src.IdContains
			} else {
				var zero string
				dst.IdContains = zero
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
		case "description_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'description_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DescriptionContains = src.DescriptionContains
			} else {
				var zero string
				dst.DescriptionContains = zero
			}
		case "attributes_contain":
			if len(subs) > 0 {
				return fmt.Errorf("'attributes_contain' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.AttributesContain = src.AttributesContain
			} else {
				dst.AttributesContain = nil
			}
		case "dev_eui_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'dev_eui_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DevEuiContains = src.DevEuiContains
			} else {
				var zero string
				dst.DevEuiContains = zero
			}
		case "join_eui_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'join_eui_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.JoinEuiContains = src.JoinEuiContains
			} else {
				var zero string
				dst.JoinEuiContains = zero
			}
		case "dev_addr_contains":
			if len(subs) > 0 {
				return fmt.Errorf("'dev_addr_contains' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DevAddrContains = src.DevAddrContains
			} else {
				var zero string
				dst.DevAddrContains = zero
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
		case "order":
			if len(subs) > 0 {
				return fmt.Errorf("'order' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Order = src.Order
			} else {
				var zero string
				dst.Order = zero
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
