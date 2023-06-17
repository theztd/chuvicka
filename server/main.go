package server

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path"
	"theztd/chuvicka/auth"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	VERSION    string = "0.1.0-alpha01"
	BucketName string
	UI         string = "false"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

//go:embed assets
var staticFS embed.FS

func validate(ctx *gin.Context) {
	headerToken := ctx.Request.Header.Get("X-Auth-Token")
	log.Println(headerToken)

	/*


		TODO: Implement token auth, token is already in model



	*/
	// if headerToken == "develop"

	if auth.ValidetaApiToken(headerToken) {
		log.Println("DEBUG: Authorized by X-Auth-Token")
		ctx.Next()
		return
	}

	authToken, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"message": "You have to be authorized before accessing required page",
		})
		ctx.Abort()
		return
	}

	if auth.JWTValidate(authToken) {
		log.Println("DEBUG: Authorized user")
		ctx.Next()
		return
	} else {
		log.Println("INFO: Invalid auth token", authToken)
		ctx.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"message": "You have to be authorized before accessing required page",
		})
		/*
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization required. Use Login form or X-Auth-Token",
			})
		*/
		ctx.Abort()
		return
	}
}

func Run() {
	r := gin.Default()

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.tmpl"))
	r.SetHTMLTemplate(tmpl)

	// r.LoadHTMLGlob("server/templates/*.tmpl")
	// r.Static("/assets", "./server/assets")
	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.FileFromFS(path.Join("/", c.Request.URL.Path), http.FS(staticFS))
	})

	// CORS allow all
	// r.Use(cors.Default())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour,
	}))

	r.GET("/_healthz/ready.json", healthStatus)
	if UI == "true" {
		r.GET("/ui/", validate, index)
	}

	r.GET("/api/metrics/", validate, metricList)
	r.POST("/api/metrics/", validate, metricCreate)
	r.DELETE("/api/metrics/:id", validate, metricDelete)

	r.GET("/login", loginPage)
	r.POST("/login", login, index)
	r.GET("/logout", logout)

	// Admin part
	//r.GET("/admin", admin)
	//r.GET("/api/tables", bucketList)
	// r.POST("/api/users", validate, usersCreate)
	// r.DELETE("/api/users/:id", validate, usersDelete)

	r.GET("/api/users/", validate, usersList)
	r.GET("/api/users/:id", validate, usersGet)
	r.POST("/api/users/:id", validate, usersEdit)

	r.Run()

}
