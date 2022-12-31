package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// views/
func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", index)
	log.Println(r.Run())

}
