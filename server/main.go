package server

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path"
	"theztd/chuvicka/auth"

	"github.com/gin-gonic/gin"
)

var BucketName string

//go:embed templates/*.tmpl
var templatesFS embed.FS

//go:embed assets
var staticFS embed.FS

func validate(ctx *gin.Context) {
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
	} else {
		log.Println("INFO: Invalid auth token", authToken)
		ctx.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"message": "You have to be authorized before accessing required page",
		})
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

	r.GET("/ui", validate, index)
	r.GET("/api/metrics", validate, metricList)
	r.POST("/api/metrics", validate, metricCreate)
	r.DELETE("/api/metrics/", validate, metricDelete)

	r.GET("/login", loginPage)
	r.POST("/login", login, index)
	r.GET("/logout", logout)

	// Admin part
	//r.GET("/admin", admin)
	//r.GET("/api/tables", bucketList)
	// r.POST("/api/tables", bucketCreate)

	r.Run()

}
