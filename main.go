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
	r.POST("/frameImage", controllers.FrameImages)
	r.Run() // listen and serve on 0.0.0.0:8080
}

