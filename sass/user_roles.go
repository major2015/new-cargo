package sass

import (
	"fmt"

	cargo_db "github.com/major2015/new-cargo/models"
	"github.com/vattle/sqlboiler/boil"
	qm "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// defines User role enums
const (
	UserGuestRoleID   = 1
	UserMemberRoleID  = 2
	UserManagerRoleID = 3
	UserOwnerRoleID   = 4
)

// UserIsGuest returns true if specifying user RoleID
// is Equals to UserGuestRoleID, otherwise false.
func UserIsGuest(u *cargo_db.User) bool {
	return u.RoleID == UserGuestRoleID
}

// UserIsMember returns true if specifying user RoleID
// is Equals to UserMemberRoleID, otherwise false.
func UserIsMember(u *cargo_db.User) bool {
	return u.RoleID == UserMemberRoleID
}

// UserIsManager returns true if specifying user RoleID
// is Equals to UserManagerRoleID, otherwise false.
func UserIsManager(u *cargo_db.User) bool {
	return u.RoleID == UserManagerRoleID
}

// UserIsOwner returns true if specifying user RoleID
// is Equals to UserOwnerRoleID, otherwise false.
func UserIsOwner(u *cargo_db.User) bool {
	return u.RoleID == UserManagerRoleID
}

// UserIsAdmin returns true if specifying user RoleId
// is equals to UserOwnerRoleID and user TenantId
// is equals to config#administrator_uuid
func UserIsAdmin(u *cargo_db.User, config Configuration) bool {
	return u.RoleID == UserOwnerRoleID &&
		u.TenantID == config.String("administrator_uuid")
}

// UserRoleName returns the role name of specifying user.
func UserRoleName(u *cargo_db.User) string {
	switch u.RoleID {
	case UserOwnerRoleID:
		return "owner"
	case UserManagerRoleID:
		return "manager"
	case UserMemberRoleID:
		return "member"
	case UserGuestRoleID:
		return "guest"
	default:
		return "invalid"
	}
}

// UserAllowedRoleNames returns allowed role names array
// of specifying user.
func UserAllowedRoleNames(u *cargo_db.User) []string {
	switch u.RoleID {
	case UserOwnerRoleID:
		return []string{"owner", "manager", "member", "guest"}
	case UserManagerRoleID:
		return []string{"manager", "member", "guest"}
	case UserMemberRoleID:
		return []string{"member", "guest"}
	case UserGuestRoleID:
		return []string{"guest"}
	default:
		return []string{}
	}
}

// Define hook to prevent deleing last owner or guest account.
func ensureOwnerAndGuest(exec boil.Executor, u *cargo_db.User) error {
	count := cargo_db.Users(qm.Where("tenant_id = ? and role_id = ?",
		u.TenantID, u.RoleID)).CountP(exec)
	if (UserIsOwner(u) || UserIsGuest(u)) && (count < 2) {
		return fmt.Errorf(
			"all accounts must have at least 1 user with role %s present",
			UserRoleName(u))
	}
	return nil
}
