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
	return currentPermissions&required != 0
}
