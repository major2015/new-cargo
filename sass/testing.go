package sass

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-mail/mail"

	"github.com/gin-gonic/gin"
	cargo_db "github.com/major2015/new-cargo/models"
	"github.com/nathanstitt/webpacking"
	"github.com/onsi/ginkgo"
	"github.com/vattle/sqlboiler/boil"
)

// TestEmailDelivery defines
type TestEmailDelivery struct {
	To       string
	Subject  string
	Contents string
}

// SendEmail bind a function
func (f *TestEmailDelivery) SendEmail(config Configuration, m *mail.Message) error {
	to := m.GetHeader("To")
	if len(to) > 0 {
		f.To = to[0]
	}
	subj := m.GetHeader("Subject")
	if len(subj) > 0 {
		f.Subject = subj[0]
	}
	buf := new(bytes.Buffer)
	_, err := m.WriteTo(buf)
	if err == nil {
		f.Contents = buf.String()
	}
	return err
}

// LastEmailDelivery var
var LastEmailDelivery *TestEmailDelivery

func testingContextMiddleware(config Configuration, tx DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbTx", tx)
		c.Set("config", config)
		c.Next()
	}
}

// TestEnv defines
type TestEnv struct {
	Router *gin.Engine
	DB     DB
	Config Configuration
	Tenant *cargo_db.Tenant
}

// RequestOptions defines
type RequestOptions struct {
	Body         *string
	SessionCooie string
	ContentType  string
	User         *cargo_db.User
}

func contentType(method string, options *RequestOptions) string {
	if options != nil && options.ContentType != "" {
		return options.ContentType
	}
	if method == "POST" {
		return "application/x-www-form-urlencoded"
	}
	return "application/json"
}

// MakeRequest bind a function
func (env *TestEnv) MakeRequest(
	method string, path string, options *RequestOptions,
) *httptest.ResponseRecorder {
	var body io.Reader
	if options != nil && options.Body != nil {
		body = strings.NewReader(*options.Body)
	}
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", contentType(method, options))
	if options != nil {
		if options.User != nil {
			req.Header.Set(
				"Cookie",
				TestingCookieForUser(options.User, env.Config),
			)
		}
	}
	resp := httptest.NewRecorder()
	env.Router.ServeHTTP(resp, req)
	return resp
}

// TestingCookieForUser defines
func TestingCookieForUser(u *cargo_db.User, config Configuration) string {
	r := gin.Default()
	InitSession("test", r, config)
	r.GET("/", func(c *gin.Context) {

		c.String(200, "")
	})
	return ""
}

// TestFlags defines
type TestFlags struct {
	DebugDB    bool
	WithRoutes func(
		*gin.Engine,
		Configuration,
		*webpacking.WebPacking,
	)
}

// Test defines
func Test(description string, flags *TestFlags, testFunc func(*TestEnv)) {
	ginkgo.It(description, func() {
		RunSpec(flags, testFunc)
	})
}

var testingDBConn *sql.DB = nil

// RunSpec defines
func RunSpec(flags *TestFlags, testFunc func(*TestEnv)) {
	boil.DebugMode = flags != nil && flags.DebugDB

	// LastEmailDelivery = &TestEmailDelivery{}

}
