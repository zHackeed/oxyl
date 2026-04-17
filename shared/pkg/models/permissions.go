package models

import (
	"errors"
	"log/slog"
)

// todo: Rethink this, might want to have a more fine grained permission system.

var (
	ErrPermissionDenied = errors.New("user does not have the required permissions to perform this action")
)

const (
	CompanyPermissionView             CompanyPermission = 1  // see company + agents, implicit for all members - 00000001
	CompanyPermissionManageAgents     CompanyPermission = 2  // add, remove, update agents - 00000010
	CompanyPermissionManageCompany    CompanyPermission = 4  // update display name, limit nodes, enable/disable - 00000100
	CompanyPermissionManageMembers    CompanyPermission = 8  // add, remove, update member permissions - 00001000
	CompanyPermissionManageThresholds CompanyPermission = 16 // manage notification thresholds - 00010000
	CompanyPermissionManageWebhooks   CompanyPermission = 32 // manage webhooks - 00100000

	// Addition of the non overlapping bits. https://www.geeksforgeeks.org/go-language/go-operators/.

	CompanyPermissionMember CompanyPermission = CompanyPermissionView |
		CompanyPermissionManageAgents // 00000011

	CompanyPermissionAdmin CompanyPermission = CompanyPermissionMember |
		CompanyPermissionManageCompany |
		CompanyPermissionManageMembers |
		CompanyPermissionManageThresholds |
		CompanyPermissionManageWebhooks // 00111111

	CompanyPermissionOwner CompanyPermission = 255
)

// CompanyPermission represents a set of permissions encoded as a bitmask.
// Permissions are combined using bitwise OR and checked using bitwise AND.
type CompanyPermission int

func HasPermission(currentPermissions CompanyPermission, required CompanyPermission) bool {
	// we must have all the required permissions to have the permission. Match all the 1s of the required permissions.
	slog.Info("checking permissions", "current", currentPermissions, "required", required, "has", currentPermissions&required, "result", currentPermissions&required != 0)
	return currentPermissions&required == required // all required bits must be set to have the permission that we are looking for
}

func (p CompanyPermission) StringifiedPermissions() []string {
	var parts []string
	if p&CompanyPermissionManageWebhooks != 0 {
		parts = append(parts, "manage_webhooks")
	}
	if p&CompanyPermissionManageThresholds != 0 {
		parts = append(parts, "manage_thresholds")
	}
	if p&CompanyPermissionManageMembers != 0 {
		parts = append(parts, "manage_members")
	}
	if p&CompanyPermissionManageCompany != 0 {
		parts = append(parts, "manage_company")
	}
	if p&CompanyPermissionManageAgents != 0 {
		parts = append(parts, "manage_agents")
	}
	if p&CompanyPermissionView != 0 {
		parts = append(parts, "view")
	}

	if p&CompanyPermissionOwner != 0 {
		parts = append(parts, "owner")
	}

	return parts
}
