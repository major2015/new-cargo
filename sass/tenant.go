package sass

import (
	"strings"

	cargo_db "github.com/major2015/new-cargo/models"
	. "github.com/volatiletech/sqlboiler/queries/qm"
)

// SignUpData defines struct for user sign up info
type SignUpData struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Tenant   string `form:"tenant"`
}

// ApplicationBootstrapData todo
type ApplicationBootstrapData struct {
	User      *cargo_db.User
	JWT       string
	WebDomain string
}

// IsEmailInUse returns true that specified email had used otherwise false
func IsEmailInUse(email string, db DB) bool {
	lowerEmail := strings.ToLower(email)
	m := Where("email = ?", lowerEmail)
	if cargo_db.Tenants(m).ExistsP(db) || cargo_db.Users(m).ExistsP {
		return true
	}
	return false
}
