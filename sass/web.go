package sass

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// InitSession defines
func InitSession(cookieName string, r *gin.Engine, config Configuration) {
	secret := []byte(config.String("session_secret"))
	SessionsKeyValue = secret
	store := cookie.NewStore(secret)
	r.Use(sessions.Sessions(cookieName, store))
}
