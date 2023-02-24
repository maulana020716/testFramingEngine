package models

import "gorm.io/gorm"

type FramedImages struct {
	gorm.Model
	ImageUrl string
    FrameUrl string
    FramedImage string
    ImageSize string
	FrameSize string
	ImagePosition string
}