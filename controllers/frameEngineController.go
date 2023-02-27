package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"go-framing-engine/initializer"
	"go-framing-engine/models"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type RequestBody struct {
	ImageUrl    string
	FrameUrl    string
	Settings    Settings
}

type Settings struct {
	ImageWidth  int
	ImageHeight int
	FrameWidth  int
	FrameHeight int
	ImagePosX   int
	ImagePosY   int
}


func FrameImages(c *gin.Context) {
	// Get the request body
	var reqBody RequestBody
	if err := c.BindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "image_processing_requests",
		Balancer: &kafka.LeastBytes{},
	})

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(reqBody.ImageUrl),
		Value: jsonBytes,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Close the Kafka producer
	err = writer.Close()
	if err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	}

	c.Status(http.StatusOK)

	
}

func ImageProcessingConsumer() {
	// Create a Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "image_processing_requests",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	// Process messages from the Kafka topic
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Unmarshal the message value to a request body
		var reqBody RequestBody
		err = json.Unmarshal(msg.Value, &reqBody)
		if err != nil {
			log.Printf("Error unmarshaling message value: %v", err)
			continue
		}

		// Perform image processing using the request body
		// Check if the combination of image and frame already made before
		var framedImage models.FramedImages
		if err := initializer.DB.Where(&models.FramedImages{ImageUrl: reqBody.ImageUrl, FrameUrl: reqBody.FrameUrl}).First(&framedImage).Error; err != nil {
			// no matching record found, create a new one
			// Encode image from url to vips image
			inputImage, err := getInputImage(reqBody.ImageUrl)
			if err != nil {
				log.Printf("Error getting input image: %v", err)
				continue
			}
			frameImage, err := getFrameImage(reqBody.FrameUrl)
			if err != nil {
				log.Printf("Error getting frame image: %v", err)
				continue
			}

			// Resize the frame image to match the size of the input image
			frameImage.Thumbnail(reqBody.Settings.FrameWidth,reqBody.Settings.FrameHeight, vips.InterestingNone)
			inputImage.Thumbnail(reqBody.Settings.ImageWidth,reqBody.Settings.ImageHeight, vips.InterestingNone)

			// Combine the input image and the frame image using the composite function
			frameImage.Composite(inputImage, vips.BlendModeDestOver,reqBody.Settings.ImagePosX,reqBody.Settings.ImagePosY)

			// Export image to jpeg	
			out, _, err := frameImage.ExportJpeg(nil)
			if err != nil {
				log.Printf("Error error exporting framed image to jpeg: %v", err)
				continue
			}			
	

			// Insert the framed image into the database
			if err := insertFramedImage(reqBody, out); err != nil {
				log.Printf("Error inserting data to database: %v", err)
				continue
			}

			// Set response headers,status and Write image data to response body
			log.Printf("Processed message for image URL %s", reqBody.ImageUrl)
			// c.Header("Content-Type", "image/jpeg")
			// c.Data(http.StatusOK, "image/jpeg", out)
			// c.Status(http.StatusOK)
		} else {
			// matching record found, send the image
			// Decode the base64 to image
			_, err := base64.StdEncoding.DecodeString(framedImage.FramedImage)
			if err != nil {
				log.Printf("Error decoding base64 to image: %v", err)
				continue
			}

			// Set response headers,status and Write image data to response body
			log.Printf("Processed message for image URL %s", reqBody.ImageUrl)
		}	
	}
}

func insertFramedImage(reqBody RequestBody, out []byte) error {
	// Encode image to base64
	encodedImage := base64.StdEncoding.EncodeToString(out)

	// Create a new FramedImages record
	insertFramedImage := models.FramedImages{
		ImageUrl:      reqBody.ImageUrl,
		FrameUrl:      reqBody.FrameUrl,
		FramedImage:   encodedImage,
		ImageSize:     strconv.Itoa(reqBody.Settings.ImageWidth) + "," + strconv.Itoa(reqBody.Settings.ImageHeight),
		FrameSize:     strconv.Itoa(reqBody.Settings.FrameWidth) + "," + strconv.Itoa(reqBody.Settings.FrameHeight),
		ImagePosition: strconv.Itoa(reqBody.Settings.ImagePosX) + "," + strconv.Itoa(reqBody.Settings.ImagePosY),
	}

	// Insert the record into the database
	result := initializer.DB.Create(&insertFramedImage) // pass pointer of data to Create
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getInputImage(imageUrl string) (*vips.ImageRef, error) {
    resp, err := http.Get(imageUrl)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    inputData, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    inputImage, err := vips.NewImageFromBuffer(inputData)
    if err != nil {
        return nil, err
    }

    return inputImage, nil
}

func getFrameImage(frameUrl string) (*vips.ImageRef, error) {
	resp, err := http.Get(frameUrl)
	if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

	inputData, err := io.ReadAll(resp.Body)
	if err != nil {
        return nil, err
    }

	frameImage, err := vips.NewImageFromBuffer(inputData)
	if err != nil {
        return nil, err
    }

	return frameImage, nil
}