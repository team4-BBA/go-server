package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	)

type Facilities struct {
	id int `json:"id"`
	similarity float32 `json:"similarity"`
}

func getRuneAt(s string, i int) rune {
	rs := []rune(s)
	return rs[i]
}

type Response struct {
	message   string    `json:"message"`
	similarity float32 `json:"similarity"`
}
func dot(v1 string , v2 string){
	//どっちもベクトル
	if (string(getRuneAt(v1, 0))=="[")&&(string(getRuneAt(v2, 0))=="["){
	//コンマごとに区切って配列化
	var  sum float32 =0.0
		//要素をfloatに変更
		// 同じインデックスの要素を掛けてsumにたす
	return sum
}
func main(){
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))
	//redis接続
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()


	router.GET("/pref", func(ctx *gin.Context) {
		userid:= ctx.Query("userId")
		fid:= ctx.Param("facilityId")

		//DB検索
		v1, err1 := redis.String(conn.Do("GET", "user:"+userid+":vec"))
		v2, err2 := redis.String(conn.Do("GET", "user:"+fid+":vec"))
		fmt.Printf(string(v1))
		if err1 != nil {
			ctx.JSON(404, gin.H{
				"message": "user  not found",
			})
		}else if err2 != nil {
				ctx.JSON(404,gin.H{
					"message": "hotel  not found",
				})
		}else {

			dotprod:=dot(v1,v2)
			ctx.JSON(200, gin.H{
				"message": "succeed",
				"user": userid,
				"facility": fid,
				"similarity": dotprod,
			})
		}
	})

	router.Run()
}