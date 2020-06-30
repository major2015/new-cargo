package sass

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	cargo_db "github.com/major2015/new-cargo/models"
)

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
