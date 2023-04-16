package server

import (
	"log"
	"net/http"
	"theztd/chuvicka/auth"

	"github.com/gin-gonic/gin"
)

func loginPage(ctx *gin.Context) {
	// zjistit jestli je prihlasenu
	message := "Welcome at the login page"
	tokenStr, err := ctx.Cookie("Authorization")
	if err != nil {
		log.Println("ERR: User without authorization cookie")
	} else {
		if auth.JWTValidate(tokenStr) {
			message = "You're already loged in, welcome back!"
		}

		ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
			// "appList": apps,
			"msg": message,
		})
	}

	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		// "appList": apps,
		"msg": message,
	})

	// pokud ano, rekneme mu kdo je
}

func login(ctx *gin.Context) {
	var usr auth.User

	login := ctx.PostForm("login")
	password := ctx.PostForm("password")
	usr, err := auth.Auth(login, password)
	if err != nil {
		log.Println("ERR: Invalid login or password", login)
		ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
			"msg":  "Incorrect login or password",
			"user": usr,
		})
		ctx.Abort()
	}
	log.Println("INFO: Loged in as", usr)
	usr.JWTGenerate()
	ctx.SetCookie("Authorization", usr.Token, 1000, "", "", false, true)
	ctx.Next()
	//ctx.Redirect(http.StatusTemporaryRedirect, "/ui")

}

func logout(ctx *gin.Context) {
	ctx.SetCookie("Authorization", "", 0, "", "", false, false)
	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		"msg": "Successfully logout",
	})
}

func delete(ctx *gin.Context) {

}
