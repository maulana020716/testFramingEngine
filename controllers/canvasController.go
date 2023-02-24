package controllers

import (
	"go-framing-engine/initializer"
	"go-framing-engine/models"
	"image/png"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tdewolff/canvas"
)

func CanvasFunc (cg *gin.Context) {
	//get data of req body
	var body struct {
		ImageUrl string
		FrameUrl string
		FramedImage string
		ImagePosition string
		FramePosition string
	}

	cg.Bind(&body)

	// Create new canvas of dimension 100x100 mm
    c := canvas.New(1002, 1024)

    // Create a canvas context used to keep drawing state
    ctx := canvas.NewContext(c)

	// Load the image data
	lenna, err := os.Open("frame_new.png")
	if err != nil {
		panic(err)
	}

	// Decode the PNG image data to an image.Image
	img, err := png.Decode(lenna)
	if err != nil {
		panic(err)
	}

	// Draw the image at coordinates (10,10) with a resolution of 30 px/mm
	ctx.DrawImage(80, 80, img, canvas.DPMM(100.0))
	
	//create post
	images := models.FramedImages{ImageUrl: body.ImageUrl, FrameUrl: body.FrameUrl, FramedImage: body.FramedImage, ImagePosition: body.ImagePosition, FramePosition: body.FramePosition}
	result := initializer.DB.Create(&images) // pass pointer of data to Create

	if result.Error != nil {
		cg.Status(400)
		return
	}

	cg.JSON(200, gin.H{
		"post": images,
	})
}