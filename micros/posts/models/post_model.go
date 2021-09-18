package models

import (
	"github.com/GMcD/ts-serverless/constants"
	uuid "github.com/gofrs/uuid"
)

type PostModel struct {
	ObjectId         uuid.UUID                     `json:"objectId" bson:"objectId"`
	PostTypeId       int                           `json:"postTypeId" bson:"postTypeId"`
	Score            int64                         `json:"score" bson:"score"`
	Votes            map[string]string             `json:"votes" bson:"votes"`
	ViewCount        int64                         `json:"viewCount" bson:"viewCount"`
	Body             string                        `json:"body" bson:"body"`
	OwnerUserId      uuid.UUID                     `json:"ownerUserId" bson:"ownerUserId"`
	OwnerDisplayName string                        `json:"ownerDisplayName" bson:"ownerDisplayName"`
	OwnerAvatar      string                        `json:"ownerAvatar" bson:"ownerAvatar"`
	URLKey           string                        `json:"urlKey" bson:"urlKey"`
	Tags             []string                      `json:"tags" bson:"tags"`
	CommentCounter   int64                         `json:"commentCounter" bson:"commentCounter"`
	Image            string                        `json:"image" bson:"image"`
	ImageFullPath    string                        `json:"imageFullPath" bson:"imageFullPath"`
	Video            string                        `json:"video" bson:"video"`
	Thumbnail        string                        `json:"thumbnail" bson:"thumbnail"`
	Album            PostAlbumModel                `json:"album" bson:"album"`
	DisableComments  bool                          `json:"disableComments" bson:"disableComments"`
	DisableSharing   bool                          `json:"disableSharing" bson:"disableSharing"`
	Deleted          bool                          `json:"deleted" bson:"deleted"`
	DeletedDate      int64                         `json:"deletedDate" bson:"deletedDate"`
	CreatedDate      int64                         `json:"created_date" bson:"created_date"`
	LastUpdated      int64                         `json:"last_updated" bson:"last_updated"`
	AccessUserList   []string                      `json:"accessUserList" bson:"accessUserList"`
	Permission       constants.UserPermissionConst `json:"permission" bson:"permission"`
	Version          string                        `json:"version" bson:"version"`
}
