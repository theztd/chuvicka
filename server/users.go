package server

import (
	"log"
	"net/http"
	authModel "theztd/chuvicka/auth/model"
	"theztd/chuvicka/metrics"

	"github.com/gin-gonic/gin"
)

func usersList(ctx *gin.Context) {
	users := []authModel.UserRead{}
	metrics.DB.Model(&authModel.User{}).Find(&users)

	ctx.JSON(http.StatusOK, gin.H{"users": users})
	return
}

func usersGet(ctx *gin.Context) {
	id := ctx.Param("id")
	usr := authModel.User{}
	metrics.DB.First(&usr, id)

	ctx.JSON(http.StatusOK, gin.H{"user": usr})
}

func usersEdit(ctx *gin.Context) {
	id := ctx.Param("id")
	type User struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	xdb := metrics.DB.Debug()

	usr := authModel.User{}
	log.Println("DEBUG: [usersEdit]", usr)
	xdb.First(&usr, id)

	data := User{}
	err := ctx.Bind(&data)
	if err != nil {
		log.Println("ERR: [users] Unable to edit user", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "Unable to process your request"})
		return
	}

	xdb.Model(&usr).Updates(authModel.User{
		Email: data.Email,
		Login: data.Login,
	})

	if len(data.Password) > 5 {
		usr.HashPassword(data.Password)
	}
	//

	ctx.JSON(http.StatusOK, gin.H{"status": "User has been updated"})
}
