package models

import "gorm.io/gorm"

type Episode struct {
    gorm.Model
    AnimeID  uint   `json:"anime_id"`
    Episode  int    `json:"episode"`
    Video    string `json:"video"`
}
