package main

import (
	"go-framing-engine/initializer"
	"go-framing-engine/models"
)

func init(){
	initializer.LoadEnv()
	initializer.ConnectToDB()
}

func main(){
	initializer.DB.AutoMigrate(&models.Post{})
	initializer.DB.AutoMigrate(&models.FramedImages{})
}