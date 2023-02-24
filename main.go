package main

import (
	"go-framing-engine/controllers"
	"go-framing-engine/initializer"

	"github.com/gin-gonic/gin"
)


func init() {
	initializer.LoadEnv()
	initializer.ConnectToDB()
}



func main() {
	r := gin.Default()
	r.POST("/posts", controllers.PostCreate)
	r.POST("/frameImage", controllers.FrameImages)
	r.GET("/posts", controllers.PostIndex)
	r.GET("/posts/:id", controllers.PostShow)
	r.POST("/posts/:id", controllers.PostUpdate)
	r.DELETE("/posts/:id", controllers.PostDelete)
	r.Run() // listen and serve on 0.0.0.0:8080
}

