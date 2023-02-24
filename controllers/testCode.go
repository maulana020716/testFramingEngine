package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
)

func TestCode (c *gin.Context) {
	// Load the original image
	inputImage, err := vips.NewImageFromFile("image.jpg")
	if err != nil {
		panic(err)
	}

	// Load the frame image
	frameImage, err := vips.NewImageFromFile("frame_new.png")
	if err != nil {
		panic(err)
	}
	

	// Resize the frame image to match the size of the input image
	err = frameImage.ResizeWithVScale(float64(7),float64(7), vips.KernelLanczos3)
	// err = inputImage.ResizeWithVScale(float64(7),float64(7), vips.KernelLanczos3)
	if err != nil {
		panic(err)

	}


	// Combine the input image and the frame image using the composite function
	err = inputImage.Composite(frameImage, vips.BlendModeOver, (inputImage.Width()-frameImage.Width())/2, (inputImage.Height()-frameImage.Height())/2)
	// err = frameImage.Composite(inputImage, vips.BlendModeOver, (frameImage.Width()-inputImage.Width())/2, (frameImage.Height()-inputImage.Height())/2)
	if err != nil {
		panic(err)
	}


	out, metadata, err := inputImage.ExportJpeg(nil)

	outFile, err := os.Create(fmt.Sprintf("new_image-%v.jpeg", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	fmt.Println(metadata)

	outFile.Write(out)

	// Save the output image
	if err != nil {
		panic(err)

	}
}