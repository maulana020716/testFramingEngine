package controllers

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

type ImageData struct {
	ImageUrl      string   `json:"imageUrl"`
	FrameUrl      string   `json:"frameUrl"`
	ImagePosition []int    `json:"imagePosition"`
	FramePosition []int    `json:"framePosition"`
}

func FrameImages (cg *gin.Context) {
	var data ImageData
		if err := cg.ShouldBindJSON(&data); err != nil {
			cg.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Load image from URL
		resp, err := http.Get(data.ImageUrl)
		if err != nil {
			cg.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			cg.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Load frame image
		f, err := os.Open(data.FrameUrl)
		if err != nil {
			cg.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()

		frame, err := png.Decode(f)
		if err != nil {
			cg.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Resize frame image to match the size of the image
		frame = resize.Resize(uint(data.ImagePosition[2]), uint(data.ImagePosition[3]), frame, resize.Lanczos3)

		// Crop image to match the position and size
		croppedImg, err := cutter.Crop(img, cutter.Config{
			Width:   data.ImagePosition[2],
			Height:  data.ImagePosition[3],
			Mode:    cutter.TopLeft, 
			Anchor:  image.Point{data.ImagePosition[0], data.ImagePosition[1]},
			Options: cutter.Ratio,
		})
		if err != nil {
			cg.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create a new image and draw the frame and cropped image onto it
		rgba := image.NewRGBA(image.Rect(0, 0, data.FramePosition[2], data.FramePosition[3]))
		draw.Draw(rgba, rgba.Bounds(), frame, image.Point{data.ImagePosition[0], data.ImagePosition[1]}, draw.Src)
		draw.Draw(rgba, rgba.Bounds(), croppedImg, image.Point{0, 0}, draw.Src)

		// Encode the image as JPEG and write it to the response writer
		if err := jpeg.Encode(cg.Writer, rgba, &jpeg.Options{Quality: 80}); err != nil {
			cg.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cg.Header("Content-Type", "image/jpeg")
		cg.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "framed_image.jpeg"))
		cg.File("framed_image.jpeg")
}