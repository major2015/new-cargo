package sass

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	cargo_db "github.com/major2015/new-cargo/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	qm "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
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
	m := qm.Where("email = ?", lowerEmail)
	if cargo_db.Tenants(m).ExistsP(db) ||
		cargo_db.Users(m).ExistsP(db) {
		return true
	}
	return false
}

// CreateTenant creates and returns a specified cargo_db.Tenant
// if not exists, otherwise returns errors
func CreateTenant(data *SignUpData, db DB) (*cargo_db.Tenant, error) {
	email := strings.ToLower(data.Email)
	if IsEmailInUse(email, db) {
		return nil, errors.New("email " + email + " is in use")
	}
	tenant := cargo_db.Tenant{
		Name:       data.Tenant,
		Email:      email,
		Identifier: slug.Make(data.Tenant),
	}
	var err error
	var admin *cargo_db.User

	if err = tenant.Insert(db, boil.Infer()); err != nil {
		return nil, err
	}

	admin = &cargo_db.User{
		Name:   data.Name,
		Email:  data.Email,
		RoleID: UserOwnerRoleID,
	}
	SetUserPassword(admin, data.Password)
	if err = tenant.AddUsers(db, true, admin); err != nil {
		return nil, err
	}

	err = tenant.AddUsers(db, true, &cargo_db.User{
		Name:   "Anonymous",
		RoleID: UserGuestRoleID,
	})
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// TenantSignupHandler returns func.
func TenantSignupHandler(afterSignUp string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var form SignUpData
		if err := c.ShouldBind(&form); err != nil {
			RenderErrorPage("Failed to read signup data, please retry", c, &err)
			return
		}
		tx := GetDB(c)
		tenant, err := CreateTenant(&form, tx)
		if err != nil {
			RenderHomePage(&form, &err, c)
			return
		}
		admin := tenant.R.Users[0]
		c.Redirect(http.StatusFound, afterSignUp)
		RenderApplication(admin, c)
	}
}

// GetConfig returns Configuration of *gin.Context.
func GetConfig(c *gin.Context) Configuration {
	config, ok := c.MustGet("config").(Configuration)
	if ok {
		return config
	}
	panic("config isn't the correct type")
}

// GetDB returns DB from *gin.Context
func GetDB(c *gin.Context) DB {
	tx, ok := c.MustGet("dbTx").(DB)
	if ok {
		return tx
	}
	panic("config isn't the correct type")
}

// RenderHomePage render a home page by *gin.
func RenderHomePage(signup *SignUpData, err *error, c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"signup": signup,
		"error":  err,
	})
}

// RenderApplication render application page by *gin.
func RenderApplication(user *cargo_db.User, c *gin.Context) {
	config := GetConfig(c)
	c.HTML(http.StatusOK, "application.html", gin.H{
		"bootstrapData": BootstrapData(user, config),
	})
}

// BootstrapData teturns template.JS
func BootstrapData(user *cargo_db.User, config Configuration) template.JS {
	type BootstrapDataT map[string]interface{}
	bootstrapData, err := json.Marshal(
		BootstrapDataT{
			"serverUrl": config.String("server_url"),
			"user":      user,
			"graphsql": BootstrapDataT{
				"token":    JWTForUser(user, config),
				"endpoint": config.String("server_url"),
			},
		})
	if err != nil {
		panic(err)
	}
	return template.JS(string(bootstrapData))
}

// RenderErrorPage render a error page by *gin.
func RenderErrorPage(message string, c *gin.Context, err *error) {
	if err != nil {
		log.Printf("Error occured: %s", *err)
	}
	c.HTML(http.StatusNotFound, "not-found.html", gin.H{
		"message": message,
	})
}

// SetUserPassword generates a securt password for user with bcrypt modules
func SetUserPassword(u *cargo_db.User, password string) {
	hashenPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.PasswordDigest = string(hashenPassword)
}
