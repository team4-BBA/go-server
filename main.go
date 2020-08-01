package main

import (

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
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


func dot(v1 string , v2 string)(float64){
	var vec1 []string
	var vec2 []string
	var sum float64=0.0
	//string array to float  array
	v1=strings.Trim(v1,"[")
	v1=strings.Trim(v1,"]")
	v2=strings.Trim(v2,"[")
	v2=strings.Trim(v2,"]")
	vec1 = strings.Split(v1, ",")
	vec2 = strings.Split(v2, ",")
	for i := 0; i < 300; i++ {
		s1,_:=strconv.ParseFloat(vec1[i], 64)
		s2,_:=strconv.ParseFloat(vec2[i], 64)
		sum =sum+ s1*s2
	}
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


	router.GET("/similarity", func(ctx *gin.Context) {
		userid:= ctx.Query("userId")
		fid:= ctx.Query("facilityId")

		//DB検索
		v1, err1 := redis.String(conn.Do("GET", "user:"+userid+":vec"))
		v2, err2 := redis.String(conn.Do("GET", "hotel:"+fid+":vec"))

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