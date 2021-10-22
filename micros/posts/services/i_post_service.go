package service

import (
	dto "github.com/GMcD/ts-serverless/micros/posts/dto"
	"github.com/GMcD/ts-serverless/micros/posts/models"
	uuid "github.com/gofrs/uuid"
	repo "github.com/red-gold/telar-core/data"
)

type PostService interface {
	SavePost(post *dto.Post) error
	FindOnePost(filter interface{}) (*dto.Post, error)
	FindPostList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Post, error)
	FindPostsIncludeProfile(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Post, error)
	QueryPost(search string, ownerUserIds []uuid.UUID, postTypeId int, sortBy string, page int64) ([]dto.Post, error)
	QueryPostIncludeUser(search string, ownerUserIds []uuid.UUID, collectiveId uuid.UUID, postTypeId int, sortBy string, page int64) ([]dto.Post, error)
	FindById(objectId uuid.UUID) (*dto.Post, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Post, error)
	FindByURLKey(urlKey string) (*dto.Post, error)
	UpdatePost(filter interface{}, data interface{}, opts ...*repo.UpdateOptions) error
	UpdateManyPost(filter interface{}, data interface{}, opts ...*repo.UpdateOptions) error
	UpdatePostById(data *models.PostUpdateModel) error
	DeletePost(filter interface{}) error
	DeletePostByOwner(ownerUserId uuid.UUID, postId uuid.UUID) error
	DeleteManyPost(filter interface{}) error
	CreatePostIndex(indexes map[string]interface{}) error
	DisableCommnet(OwnerUserId uuid.UUID, objectId uuid.UUID, value bool) error
	DisableSharing(OwnerUserId uuid.UUID, objectId uuid.UUID, value bool) error
	IncrementScoreCount(objectId uuid.UUID, ownerUserId uuid.UUID, displayName string, avatar string) error
	DecrementScoreCount(objectId uuid.UUID, ownerUserId uuid.UUID, displayName string, avatar string) error
	Increment(objectId uuid.UUID, field string, value int) error
	IncrementCommentCount(objectId uuid.UUID) error
	DecrementCommentCount(objectId uuid.UUID) error
	UpdatePostProfile(ownerUserId uuid.UUID, ownerDisplayName string, ownerAvatar string) error
	UpdatePostURLKey(postId uuid.UUID, urlKey string) error
}
