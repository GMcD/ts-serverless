package dto

import uuid "github.com/gofrs/uuid"

type CollectiveRelMeta struct {
	CollectiveId uuid.UUID `json:"collectiveId" bson:"collectiveId"`
	CreatedDate  int64     `json:"created_date" bson:"created_date"`
	Title        string    `json:"title" bson:"title"`
	Banner       string    `json:"banner" bson:"banner"`
	Avatar       string    `json:"avatar" bson:"avatar"`
	Tagline      string    `json:"tagline" bson:"tagline"`
}
