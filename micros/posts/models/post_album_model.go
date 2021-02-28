package models

import uuid "github.com/satori/go.uuid"

type PostAlbumModel struct {
	Count   int       `json:"count"`
	Cover   string    `json:"cover"`
	CoverId uuid.UUID `json:"coverId"`
	Photos  []string  `json:"photos"`
	Title   string    `json:"title"`
}
