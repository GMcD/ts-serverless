package models

import (
	"github.com/GMcD/ts-serverless/constants"
	dto "github.com/GMcD/ts-serverless/micros/posts/dto"
	uuid "github.com/gofrs/uuid"
)

type CreateCollectivesPostModel struct {
	ObjectId         uuid.UUID                     `json:"objectId"`
	CollectiveId     uuid.UUID                     `json:"collectiveId"`
	PostTypeId       int                           `json:"postTypeId"`
	Score            int64                         `json:"score"`
	Votes            []dto.VoterProfile            `json:"votes"`
	ViewCount        int64                         `json:"viewCount"`
	Body             string                        `json:"body"`
	OwnerUserId      uuid.UUID                     `json:"ownerUserId"`
	OwnerDisplayName string                        `json:"ownerDisplayName"`
	OwnerAvatar      string                        `json:"ownerAvatar"`
	URLKey           string                        `json:"urlKey"`
	Tags             []string                      `json:"tags"`
	CommentCounter   int64                         `json:"commentCounter"`
	Image            string                        `json:"image"`
	ImageFullPath    string                        `json:"imageFullPath"`
	Video            string                        `json:"video"`
	Thumbnail        string                        `json:"thumbnail"`
	Album            PostAlbumModel                `json:"album"`
	DisableComments  bool                          `json:"disableComments"`
	DisableSharing   bool                          `json:"disableSharing"`
	Deleted          bool                          `json:"deleted"`
	DeletedDate      int64                         `json:"deletedDate"`
	CreatedDate      int64                         `json:"created_date"`
	LastUpdated      int64                         `json:"last_updated"`
	AccessUserList   []string                      `json:"accessUserList"`
	Permission       constants.UserPermissionConst `json:"permission"`
	Version          string                        `json:"version"`
}
