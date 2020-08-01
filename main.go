package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"io/ioutil"
	"strconv"
	"strings"
	"github.com/gin-gonic/autotls"
	"log"
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
	if(v1[0:1]=="["&&v2[0:1]=="[") {
		v1 = strings.Trim(v1, "[")
		v1 = strings.Trim(v1, "]")
		v2 = strings.Trim(v2, "[")
		v2 = strings.Trim(v2, "]")
		vec1 = strings.Split(v1, ",")
		vec2 = strings.Split(v2, ",")
		for i := 0; i < 300; i++ {
			s1, _ := strconv.ParseFloat(vec1[i], 64)
			s2, _ := strconv.ParseFloat(vec2[i], 64)
			sum = sum + s1*s2
		}
		return sum
	}else{
		return 999999;
	}
}

func main(){




	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))
	//redis接続
	conn, err := redis.Dial("tcp", "ec2-18-182-35-25.ap-northeast-1.compute.amazonaws.com:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()



	router.POST("/pref", func(ctx *gin.Context) {
		userId := ctx.PostForm("userId")
		words := ctx.PostForm("words")
		//cloudfuncitonに投げた結果
		url := "https://us-central1-arctic-conduit-280420.cloudfunctions.net/GetVector?type=0&value="+words

		resp, _ := http.Get(url)
		defer resp.Body.Close()

		byteArray, _ := ioutil.ReadAll(resp.Body)
		if(string(byteArray)==""){
			fmt.Printf(string(byteArray))
			ctx.JSON(500, gin.H{
				"message": "failed in creating",
			})
		}else{
			_, err = conn.Do("SET", "user:"+userId+":vec", string(byteArray) )
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "failed in creating",
			})
		} else {
			ctx.JSON(200, gin.H{
				"message": "created",
			})
			}
		}
	})





	router.GET("/pref", func(ctx *gin.Context) {
		userid:= ctx.Query("userId")
		const num = 100

			//DB検索
		v1, err := redis.String(conn.Do("GET", "user:"+userid+":vec"))
			if err != nil {
				ctx.JSON(404, gin.H{
					"message": "invalid uesr id",
				})
			}
		var sims [40000]float64;//類似度を入れる配列

		fmt.Printf("connection done")

		var args [] interface{}

		for i := 0; i < 40000; i++ {
			args = append(args,"hotel:"+strconv.Itoa(i)+":vec")
		}



			v2, _:= redis.Strings(conn.Do("MGET", args...))

		fmt.Printf(v2[0])
		for i := 0; i < 40000; i++ {
			if(v2[i]==""){
				sims[i] = -999999.0;
			}else {

				dp := dot(v1, string(v2[i]))
				if (dp == 999999) { //データベース登録ミスは類似度0
					sims[i] = -999999.0
				} else {
					sims[i] = dot(v1, v2[i])
				}
			}
			}


		fmt.Printf("dot prod done")
		var large_indexes [num]int64;
		//初期化
		for i := 0; i < num; i++ {
			large_indexes [i]=0//まんまhotelid
		}



		//最大インデックス格納
		var reco [] interface{}//最大順に類似度を格納
		for i := 0; i < num; i++ {
			var max float64= -999999.0
			var j int64=0
			for j = 0; j < 4000; j++ {

				if(max<sims[j]){

					max=sims[j]
					large_indexes[i]=j//最大値の更新
				}
			}
			reco=append(reco,fmt.Sprintf("%f", sims[large_indexes[i]]))
			sims[large_indexes[i]]=-999999.0//最大値は消す。
		}

		fmt.Printf("max serach done")

			ctx.JSON(200, gin.H{
				"message": "succeed",
				"hotelids":large_indexes,
				"similarities":reco,
			})
	})


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
			dotprod := dot(v1, v2)
			if(dotprod==999999) {
				ctx.JSON(404,gin.H{
					"message": "dot product failed in calcurating",
				})
			}else{
				ctx.JSON(200, gin.H{
					"message":    "succeed",
					"user":       userid,
					"facility":   fid,
					"similarity": dotprod,
				})
			}
		}
	})




log.Fatal(autotls.Run(router, "ec2-18-182-35-25.ap-northeast-1.compute.amazonaws.com"))
}
