package models

import "gorm.io/gorm"

type Process struct {
	gorm.Model
	UserID  uint   `json:"user_id" gorm:"not null"`
	VideoID uint   `json:"video_id" gorm:"not null;unique"`
	Status  string `json:"status" gorm:"type:enum('loading', 'completed', 'failed');default:'loading'"`
	Message string `json:"message"`
}

type Anime struct {
    gorm.Model
	SeasonID    uint   `json:"season_id"`
	StudioID    uint   `json:"studio_id"`
	GenreID     uint   `json:"genre_id"`
	Title       string `json:"title"`
	EnTitle     string `json:"en_title"`
	Type        string `json:"type" gorm:"type:enum('TV', 'Movie', 'OVA');default:'TV'"`
	Status      string `json:"status" gorm:"type:enum('ongoing', 'completed');default:'ongoing'"`
	IsHidden    bool   `json:"is_hidden" gorm:"default:false"`
	Description string `json:"description"`
	CompletedAt *gorm.DeletedAt `json:"completed_at" gorm:"index"`
	Photo       string `json:"photo"`  // Add this field for the photo URL
}

type Episode struct {
	gorm.Model
	AnimeID   uint `json:"anime_id" gorm:"not null"`
	Episode   int  `json:"episode"`
	Video     string `json:"video"` 
	IsHidden  bool `json:"is_hidden" gorm:"default:false"`
	CreatedBy uint `json:"created_by"`
}

type Video struct {
	gorm.Model
	EpisodeID  uint   `json:"episode_id" gorm:"not null"`
	URL        string `json:"url"`
	Driver     string `json:"driver"`
	Resolution string `json:"resolution"`
}

type Genre struct {
	gorm.Model
	Name string `json:"name"`
}

type Studio struct {
	gorm.Model
	Name string `json:"name"`
}

type Season struct {
	gorm.Model
	Name string `json:"name"`
}

type AnimeGenre struct {
	gorm.Model
	AnimeID uint `json:"anime_id" gorm:"not null"`
	GenreID uint `json:"genre_id" gorm:"not null"`
}
