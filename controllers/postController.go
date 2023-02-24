package controllers

import (
	"fmt"
	"go-framing-engine/initializer"
	"go-framing-engine/models"

	"github.com/gin-gonic/gin"
)

func PostCreate (c *gin.Context) {
	//get data of req body
	var body struct {
		Title string
		Body string
	}

	c.Bind(&body)

	//create post
	post := models.Post{Title: body.Title, Body: body.Body}
	result := initializer.DB.Create(&post) // pass pointer of data to Create

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostIndex (c *gin.Context) {
	var posts []models.Post
	initializer.DB.Find(&posts)

	c.JSON(200, gin.H{
		"post": posts,
	})
}

func PostShow (c *gin.Context) {
	id := c.Param("id")

	var posts []models.Post
	initializer.DB.First(&posts,id)

	c.JSON(200, gin.H{
		"post": posts,
	})
}

func PostUpdate (c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Title string
		Body string
	}

	c.Bind(&body)

	var posts models.Post
	initializer.DB.First(&posts,id)

	posts.Title = body.Title
	posts.Body = body.Body
	initializer.DB.Save(&posts)

	c.JSON(200, gin.H{
		"post": posts,
	})
}

func PostDelete (c *gin.Context) {
	id := c.Param("id")

	initializer.DB.Delete(&models.Post{},id)

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Post with id %v deleted",id),
	})
}