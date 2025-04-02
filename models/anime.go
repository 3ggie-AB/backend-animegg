package models

import "gorm.io/gorm"

type Anime struct {
    gorm.Model
    Title       string   `json:"title"`
    Photo       string   `json:"photo"`
    Tags        string   `json:"tags"`
    Description string   `json:"description"`
    Episodes    []Episode `gorm:"foreignKey:AnimeID"`
}