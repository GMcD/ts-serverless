package service

import (
	dto "github.com/GMcD/ts-serverless/micros/user-rels/dto"
	uuid "github.com/gofrs/uuid"
)

type CollectiveRelService interface {
	FollowCollective(leftUser dto.UserRelMeta, collective dto.CollectiveRelMeta, tags []string) error
	UnfollowCollective(leftId uuid.UUID, collectiveId uuid.UUID) error
	DeleteCollectiveRel(filter interface{}) error
	SaveCollectiveRel(collectiveRel *dto.CollectiveRel) error
	GetCollectiveFollowing(userId uuid.UUID) ([]dto.CollectiveRel, error)
	FindCollectiveById(objectId uuid.UUID) (*dto.CollectiveRel, error)
	FindCollectiveRelsIncludeProfile(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.CollectiveRel, error)
	UpdateCollectiveRel(filter interface{}, data interface{}) error
	DeleteManyCollectiveRel(filter interface{}) error
	CreateCollectiveRelIndex(indexes map[string]interface{}) error
	GetFollowers(collectiveId uuid.UUID) ([]dto.CollectiveRel, error)
}
