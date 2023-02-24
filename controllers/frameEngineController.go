package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
)


func FrameImages (c *gin.Context) {
	var body struct {
		ImageUrl string
		FrameUrl string
		FramedImage string
		ImagePosition string
		FramePosition string
	}

	c.Bind(&body)
	
	inputImage := getInputImage(body.ImageUrl)
	frameImage := getFrameImage(body.FrameUrl)	

	// Resize the frame image to match the size of the input image
	frameImage.Thumbnail(1002,1024, vips.InterestingNone)
	inputImage.Thumbnail(760,820, vips.InterestingNone)
	// err = inputImage.ResizeWithVScale(float64(7),float64(7), vips.KernelLanczos3)
	


	// Combine the input image and the frame image using the composite function
	frameImage.Composite(inputImage, vips.BlendModeDestOver, 80,85)
	// err = frameImage.Composite(inputImage, vips.BlendModeOver, (frameImage.Width()-inputImage.Width())/2, (frameImage.Height()-inputImage.Height())/2)
	

	// Export image to jpeg
	out, metadata, err := frameImage.ExportJpeg(nil)

	outFile, err := os.Create(fmt.Sprintf("new_image-%v.jpeg", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	fmt.Println(metadata)

	outFile.Write(out)

	// c.Status(200)

	// Set response headers
	c.Header("Content-Type", "image/jpeg")
	// c.Header("Content-Disposition", "attachment; filename=new_image.jpeg")
	// c.Header("Content-Length", strconv.Itoa(len(out)))

	// Write image data to response body
	c.Data(http.StatusOK, "image/jpeg", out)

	c.Status(http.StatusOK)
}

func getInputImage(imageUrl string) *vips.ImageRef {
	resp, err := http.Get(imageUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	inputData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	inputImage, err := vips.NewImageFromBuffer(inputData)
	if err != nil {
		panic(err)
	}

	return inputImage
}

func getFrameImage(frameUrl string) *vips.ImageRef {
	resp, err := http.Get(frameUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	inputData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	frameImage, err := vips.NewImageFromBuffer(inputData)
	if err != nil {
		panic(err)
	}

	return frameImage
}