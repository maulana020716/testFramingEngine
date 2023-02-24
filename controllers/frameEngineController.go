package controllers

import (
	"encoding/base64"
	"go-framing-engine/initializer"
	"go-framing-engine/models"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
)


func FrameImages (c *gin.Context) {
	var body struct {
		ImageUrl string
		FrameUrl string
		FramedImage string
		ImageSize []int
		FrameSize []int
		ImagePosition []int
	}

	c.Bind(&body)
	
	inputImage := getInputImage(body.ImageUrl)
	frameImage := getFrameImage(body.FrameUrl)	

	// Resize the frame image to match the size of the input image
	frameImage.Thumbnail(body.FrameSize[0],body.FrameSize[1], vips.InterestingNone)
	inputImage.Thumbnail(body.ImageSize[0],body.ImageSize[0], vips.InterestingNone)

	// Combine the input image and the frame image using the composite function
	frameImage.Composite(inputImage, vips.BlendModeDestOver,body.ImagePosition[0],body.ImagePosition[0])	

	// Export image to jpeg
	out, _, err := frameImage.ExportJpeg(nil)
	if err != nil {
		panic(err)
	}
	
	// if we want to write the file to local
	// _, err = os.Create(fmt.Sprintf("new_image-%v.jpeg", time.Now().Unix()))
	// if err != nil {
	// 	panic(err)
	// }
	// outFile.Write(out)

	// Encode image to base64 then insert to db
	encodedImage := base64.StdEncoding.EncodeToString(out)
	insertFramedImage := models.FramedImages{ImageUrl: body.ImageUrl, FrameUrl: body.FrameUrl, FramedImage: encodedImage, ImageSize: strconv.Itoa(body.ImageSize[0])+","+strconv.Itoa(body.ImageSize[1]),FrameSize: strconv.Itoa(body.FrameSize[0])+","+strconv.Itoa(body.FrameSize[1]),ImagePosition: strconv.Itoa(body.ImagePosition[0])+","+strconv.Itoa(body.ImagePosition[1])}
	result := initializer.DB.Create(&insertFramedImage) // pass pointer of data to Create
	if result.Error != nil {
		c.Status(400)
		return
	}

	// Set response headers
	c.Header("Content-Type", "image/jpeg")

	// Write image data to response body
	c.Data(http.StatusOK, "image/jpeg", out)

	// Set status
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