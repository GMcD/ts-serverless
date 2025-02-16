package service

import (
	dto "github.com/GMcD/ts-serverless/micros/circles/dto"
	uuid "github.com/gofrs/uuid"
)

type CircleService interface {
	SaveCircle(circle *dto.Circle) error
	FindOneCircle(filter interface{}) (*dto.Circle, error)
	FindCircleList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Circle, error)
	QueryCircle(search string, ownerUserId *uuid.UUID, sortBy string, page int64) ([]dto.Circle, error)
	FindById(objectId uuid.UUID) (*dto.Circle, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Circle, error)
	UpdateCircle(filter interface{}, data interface{}) error
	UpdateCircleById(data *dto.Circle) error
	DeleteCircle(filter interface{}) error
	DeleteCircleByOwner(ownerUserId uuid.UUID, circleId uuid.UUID) error
	DeleteManyCircle(filter interface{}) error
	CreateCircleIndex(indexes map[string]interface{}) error
}
