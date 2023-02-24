package models

import "gorm.io/gorm"

type FramedImages struct {
	gorm.Model
	ImageUrl string
    FrameUrl string
    FramedImage string
    ImagePosition string
    FramePosition string
}