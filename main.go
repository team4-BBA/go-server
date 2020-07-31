package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	)

func main() {


	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))


	router.GET("/room/:id", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		username := session.Get("username")
		ishost:=session.Get("ishost")
		id:= ctx.Param("id")
		ctx.HTML(200, "index.html", gin.H{"page": 1, "error": -1,"name":username,"ishost":ishost,"roomid":id})
	})
	router.POST("/play/:id", func(ctx *gin.Context) {
		team := ctx.PostForm("player")
		name := ctx.PostForm("name")
		enemyteam := ctx.PostForm("enemyteam")
		enemyname := ctx.PostForm("enemyname")
		ishost := ctx.PostForm("ishost")
		id:= "you"
		enemyid:="enemy"
		ctx.HTML(200, "index.html", gin.H{"page": 2, "ishost":ishost,"team":team,"name":name,"id":id,"enemyteam":enemyteam,"enemyname":enemyname,"enemyid":enemyid})
	})
	router.Run()
}