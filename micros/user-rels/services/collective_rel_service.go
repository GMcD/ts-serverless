package service

import (
	"fmt"

	"github.com/GMcD/ts-serverless/micros/user-rels/dto"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/utils"

	repo "github.com/red-gold/telar-core/data"
	"github.com/red-gold/telar-core/data/mongodb"
	mongoRepo "github.com/red-gold/telar-core/data/mongodb"
)

type CollectiveRelServiceImpl struct {
	CollectiveRelRepo repo.Repository
}

func NewCollectiveRelService(db interface{}) (CollectiveRelService, error) {

	collectiveRelService := &CollectiveRelServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		collectiveRelService.CollectiveRelRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return collectiveRelService, nil
}

//FollowCollective create relation between two users
func (s CollectiveRelServiceImpl) FollowCollective(leftUser dto.UserRelMeta, collective dto.CollectiveRelMeta, tags []string) error {

	newCollectiveRel := &dto.CollectiveRel{
		Left:         leftUser,
		LeftId:       leftUser.UserId,
		Collective:   collective,
		CollectiveId: collective.CollectiveId,
		Rel:          []string{leftUser.UserId.String(), collective.CollectiveId.String()},
		Tags:         tags,
	}
	err := s.SaveCollectiveRel(newCollectiveRel)
	return err
}

// UnfollowCollective delete relation between a collective, and a user (inherits left-syntax from UnFollowUser)
func (s CollectiveRelServiceImpl) UnfollowCollective(leftId uuid.UUID, collectiveId uuid.UUID) error {

	filter := struct {
		LeftId       uuid.UUID `json:"leftId" bson:"leftId"`
		CollectiveId uuid.UUID `json:"collectiveId" bson:"collectiveId"`
	}{
		LeftId:       leftId,
		CollectiveId: collectiveId,
	}
	err := s.DeleteCollectiveRel(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteColectiveRel delete collectiveRel by filter
func (s CollectiveRelServiceImpl) DeleteCollectiveRel(filter interface{}) error {

	result := <-s.CollectiveRelRepo.Delete(collectiveRelCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// SaveUserRel save the userRel
func (s CollectiveRelServiceImpl) SaveCollectiveRel(collectiveRel *dto.CollectiveRel) error {

	if collectiveRel.ObjectId == uuid.Nil {
		var uuidErr error
		collectiveRel.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if collectiveRel.CreatedDate == 0 {
		collectiveRel.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.CollectiveRelRepo.Save(collectiveRelCollectionName, collectiveRel)

	return result.Error
}

// GetCollectiveFollowing Get Collective following by collectiveId
func (s CollectiveRelServiceImpl) GetCollectiveFollowing(userId uuid.UUID) ([]dto.CollectiveRel, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		LeftId uuid.UUID `json:"leftId" bson:"leftId"`
	}{
		LeftId: userId,
	}
	return s.FindCollectiveRelsIncludeProfile(filter, 0, 0, sortMap)
}

func (s CollectiveRelServiceImpl) FindCollectiveRelsIncludeProfile(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.CollectiveRel, error) {
	var pipeline []interface{}

	matchOperator := make(map[string]interface{})
	matchOperator["$match"] = filter

	sortOperator := make(map[string]interface{})
	sortOperator["$sort"] = sort

	pipeline = append(pipeline, matchOperator, sortOperator)

	if skip > 0 {
		skipOperator := make(map[string]interface{})
		skipOperator["$skip"] = skip
		pipeline = append(pipeline, skipOperator)
	}

	if limit > 0 {
		limitOperator := make(map[string]interface{})
		limitOperator["$limit"] = limit
		pipeline = append(pipeline, limitOperator)
	}

	// Add left user pipeline
	lookupLeftUser := make(map[string]map[string]string)
	lookupLeftUser["$lookup"] = map[string]string{
		"localField":   "leftId",
		"from":         "userProfile",
		"foreignField": "objectId",
		"as":           "leftUser",
	}

	unwindLeftUser := make(map[string]interface{})
	unwindLeftUser["$unwind"] = "$leftUser"
	pipeline = append(pipeline, lookupLeftUser, unwindLeftUser)

	// Add right user pipeline
	lookupCollective := make(map[string]map[string]string)
	lookupCollective["$lookup"] = map[string]string{
		"localField":   "collectiveId",
		"from":         "collectiveProfile",
		"foreignField": "objectId",
		"as":           "collective",
	}

	unwindCollective := make(map[string]interface{})
	unwindCollective["$unwind"] = "$collective"
	pipeline = append(pipeline, lookupCollective, unwindCollective)
	log.Info("pipeline %v", pipeline)

	projectOperator := make(map[string]interface{})
	project := make(map[string]interface{})

	// Add project operator
	project["objectId"] = 1
	project["created_date"] = 1
	project["leftId"] = 1
	project["collectiveId"] = 1
	project["rel"] = 1
	project["tags"] = 1

	// left user
	project["left.userId"] = "$leftId"
	project["left.fullName"] = "$leftUser.fullName"
	project["left.instagramId"] = "$leftUser.instagramId"
	project["left.twitterId"] = "$leftUser.twitterId"
	project["left.linkedInId"] = "$leftUser.linkedInId"
	project["left.facebookId"] = "$leftUser.facebookId"
	project["left.socialName"] = "$leftUser.socialName"
	project["left.created_date"] = "$leftUser.created_date"
	project["left.banner"] = "$leftUser.banner"
	project["left.avatar"] = "$leftUser.avatar"
	// Right user
	project["collective.collectiveId"] = "$collectiveId"
	project["collective.title"] = "$collectiveUser.title"
	project["collective.tagline"] = "$collectiveUser.tagline"
	project["collective.created_date"] = "$collective.created_date"
	project["collective.banner"] = "$collective.banner"
	project["collective.avatar"] = "$collective.avatar"

	projectOperator["$project"] = project

	pipeline = append(pipeline, projectOperator)

	result := <-s.CollectiveRelRepo.Aggregate(collectiveRelCollectionName, pipeline)

	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var postList []dto.CollectiveRel
	for result.Next() {
		var post dto.CollectiveRel
		errDecode := result.Decode(&post)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.CollectiveRel")
		}
		postList = append(postList, post)
	}

	return postList, nil
}
