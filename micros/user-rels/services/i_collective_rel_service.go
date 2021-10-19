package service

import (
	dto "github.com/GMcD/ts-serverless/micros/user-rels/dto"
	uuid "github.com/gofrs/uuid"
)

type CollectiveRelService interface {
	CollectiveFollowHandle(leftUser dto.UserRelMeta, collective dto.CollectiveRelMeta) error
	FollowCollective(leftUser dto.UserRelMeta, collective dto.CollectiveRelMeta, tags []string) error
	SaveCollectiveRel(collectiveRel *dto.CollectiveRel) error
	FindCollectiveById(objectId uuid.UUID) (*dto.CollectiveRel, error)
	UpdateCollectiveRel(filter interface{}, data interface{}) error
	DeleteCollectiveRel(filter interface{}) error
	DeleteManyCollectiveRel(filter interface{}) error
	CreateCollectiveRelIndex(indexes map[string]interface{}) error
	GetFollowers(collectiveId uuid.UUID) ([]dto.CollectiveRel, error)
	GetFollowing(userId uuid.UUID) ([]dto.CollectiveRel, error)
	UnfollowCollective(leftId uuid.UUID, collectiveId uuid.UUID) error
}
