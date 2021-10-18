package dto

import uuid "github.com/gofrs/uuid"

type CollectiveRel struct {
	ObjectId     uuid.UUID         `json:"objectId" bson:"objectId"`
	CreatedDate  int64             `json:"created_date" bson:"created_date"`
	Left         UserRelMeta       `json:"left" bson:"left"`
	LeftId       uuid.UUID         `json:"leftId" bson:"leftId"`
	Collective   CollectiveRelMeta `json:"collective" bson:"collective"`
	CollectiveId uuid.UUID         `json:"collectiveId" bson:"collectiveId"`
	Rel          []string          `json:"rel" bson:"rel"`
	Tags         []string          `json:"tags" bson:"tags"`
}
