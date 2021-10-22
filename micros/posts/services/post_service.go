package service

import (
	"fmt"
	"strings"

	dto "github.com/GMcD/ts-serverless/micros/posts/dto"
	"github.com/GMcD/ts-serverless/micros/posts/models"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/config"
	coreData "github.com/red-gold/telar-core/data"
	mongoRepo "github.com/red-gold/telar-core/data/mongodb"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostService handlers with injected dependencies
type PostServiceImpl struct {
	PostRepo coreData.Repository
}

// NewPostService initializes PostService's dependencies and create new PostService struct
func NewPostService(db interface{}) (PostService, error) {

	postService := &PostServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongoRepo.MongoDatabase)
		postService.PostRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return postService, nil
}

// SavePost save the post
func (s PostServiceImpl) SavePost(post *dto.Post) error {

	if post.ObjectId == uuid.Nil {
		var uuidErr error
		post.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if post.CreatedDate == 0 {
		post.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.PostRepo.Save(postCollectionName, post)

	return result.Error
}

// FindOnePost get one post
func (s PostServiceImpl) FindOnePost(filter interface{}) (*dto.Post, error) {

	result := <-s.PostRepo.FindOne(postCollectionName, filter)
	if result.Error() != nil {
		return nil, result.Error()
	}

	var postResult dto.Post
	errDecode := result.Decode(&postResult)
	if errDecode != nil {
		// return nil, fmt.Errorf("Error decoding on dto.Post")
		log.Info("Error decoding on dto.Post : %v", postResult)
	}
	return &postResult, nil
}

// FindPostList get all posts by filter
func (s PostServiceImpl) FindPostList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Post, error) {

	result := <-s.PostRepo.Find(postCollectionName, filter, limit, skip, sort)

	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var postList []dto.Post
	for result.Next() {
		var post dto.Post
		errDecode := result.Decode(&post)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Post")
		}
		postList = append(postList, post)
	}

	return postList, nil
}

// FindPostsIncludeProfile get all posts by filter including user profile entity
func (s PostServiceImpl) FindPostsIncludeProfile(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Post, error) {
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

	lookupOperator := make(map[string]interface{})
	lookupOperator["$lookup"] = map[string]string{
		"localField":   "ownerUserId",
		"from":         "userProfile",
		"foreignField": "objectId",
		"as":           "userinfo",
	}

	unwindOperator := make(map[string]interface{})
	unwindOperator["$unwind"] = "$userinfo"

	projectOperator := make(map[string]interface{})
	project := make(map[string]interface{})

	project["objectId"] = 1
	project["collectiveId"] = 1
	project["postTypeId"] = 1
	project["score"] = 1
	project["votes"] = 1
	project["viewCount"] = 1
	project["body"] = 1
	project["ownerUserId"] = 1
	project["ownerDisplayName"] = "$userinfo.fullName"
	project["ownerAvatar"] = "$userinfo.avatar"
	project["tags"] = 1
	project["commentCounter"] = 1
	project["image"] = 1
	project["imageFullPath"] = 1
	project["video"] = 1
	project["thumbnail"] = 1
	project["album"] = 1
	project["disableComments"] = 1
	project["disableSharing"] = 1
	project["deleted"] = 1
	project["deletedDate"] = 1
	project["created_date"] = 1
	project["last_updated"] = 1
	project["accessUserList"] = 1
	project["permission"] = 1
	project["version"] = 1

	projectOperator["$project"] = project

	pipeline = append(pipeline, lookupOperator, unwindOperator, projectOperator)

	log.Info("FindPostsIncludeProfile pipeline : %s", pipeline)

	result := <-s.PostRepo.Aggregate(postCollectionName, pipeline)

	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var postList []dto.Post
	for result.Next() {
		var post dto.Post
		errDecode := result.Decode(&post)
		if errDecode != nil {
			// return nil, fmt.Errorf("Error decoding on dto.Post")
			log.Info("Error decoding on dto.Post : %v", post)
		}
		postList = append(postList, post)
	}

	return postList, nil
}

// QueryPost get all posts by query
func (s PostServiceImpl) QueryPost(search string, ownerUserIds []uuid.UUID, postTypeId int, sortBy string, page int64) ([]dto.Post, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems

	filter := make(map[string]interface{})
	if search != "" {
		// filter["$text"] = coreData.SearchOperator{Search: search}
		terms := strings.Split(search, " ")
		regexp := strings.Join(terms, "|")
		filter["$or"] = bson.A{
			bson.D{{"body", primitive.Regex{Pattern: regexp, Options: "i"}}},
			bson.D{{"ownerDisplayName", primitive.Regex{Pattern: regexp, Options: "i"}}},
		}
	}
	if len(ownerUserIds) > 0 {
		inFilter := make(map[string]interface{})
		inFilter["$in"] = ownerUserIds
		filter["ownerUserId"] = inFilter
	}
	if postTypeId > 0 {
		filter["postTypeId"] = postTypeId
	}
	log.Info("FindPostList filter : %s", filter)
	result, err := s.FindPostList(filter, limit, skip, sortMap)

	return result, err
}

// QueryPostIncludeUser get all posts by query including user entity
func (s PostServiceImpl) QueryPostIncludeUser(search string, ownerUserIds []uuid.UUID, collectiveIds []uuid.UUID, postTypeId int, sortBy string, page int64) ([]dto.Post, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems

	filter := make(map[string]interface{})
	if search != "" {
		//filter["$text"] = coreData.SearchOperator{Search: search}
		terms := strings.Split(search, " ")
		regexp := strings.Join(terms, "|")
		filter["$or"] = bson.A{
			bson.D{{"body", primitive.Regex{Pattern: regexp, Options: "i"}}},
			bson.D{{"ownerDisplayName", primitive.Regex{Pattern: regexp, Options: "i"}}},
		}
	}
	if len(ownerUserIds) > 0 {
		inFilter := make(map[string]interface{})
		inFilter["$in"] = ownerUserIds
		filter["ownerUserId"] = inFilter
	}
	if postTypeId > 0 {
		filter["postTypeId"] = postTypeId
	}
	if len(collectiveIds) > 0 {
		inFilter := make(map[string]interface{})
		inFilter["$in"] = collectiveIds
		filter["collectiveId"] = inFilter
	}
	log.Info("FindPostsIncludeProfile filter : %s", filter)
	result, err := s.FindPostsIncludeProfile(filter, limit, skip, sortMap)

	return result, err
}

// FindByOwnerUserId find by owner user id
func (s PostServiceImpl) FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Post, error) {
	sortMap := make(map[string]int)
	sortMap["created_date"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindPostList(filter, 0, 0, sortMap)
}

// FindById find by post id
func (s PostServiceImpl) FindById(objectId uuid.UUID) (*dto.Post, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOnePost(filter)
}

// FindByURLKey find by URL key
func (s PostServiceImpl) FindByURLKey(urlKey string) (*dto.Post, error) {

	filter := struct {
		URLKey string `json:"urlKey" bson:"urlKey"`
	}{
		URLKey: urlKey,
	}
	return s.FindOnePost(filter)
}

// UpdatePost update the post
func (s PostServiceImpl) UpdatePost(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error {

	result := <-s.PostRepo.Update(postCollectionName, filter, data, opts...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateManyPost update the post
func (s PostServiceImpl) UpdateManyPost(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error {

	result := <-s.PostRepo.UpdateMany(postCollectionName, filter, data, opts...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdatePost update the post
func (s PostServiceImpl) UpdatePostById(data *models.PostUpdateModel) error {
	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    data.ObjectId,
		OwnerUserId: data.OwnerUserId,
	}

	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdatePost(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// DeletePost delete post by filter
func (s PostServiceImpl) DeletePost(filter interface{}) error {

	result := <-s.PostRepo.Delete(postCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeletePost delete post by ownerUserId and postId
func (s PostServiceImpl) DeletePostByOwner(ownerUserId uuid.UUID, postId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    postId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeletePost(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyPost delete many post by filter
func (s PostServiceImpl) DeleteManyPost(filter interface{}) error {

	result := <-s.PostRepo.Delete(postCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreatePostIndex create index for post search.
func (s PostServiceImpl) CreatePostIndex(indexes map[string]interface{}) error {
	result := <-s.PostRepo.CreateIndex(postCollectionName, indexes)
	return result
}

// IncrementScoreCount increment score of post
func (s PostServiceImpl) IncrementScoreCount(objectId uuid.UUID, ownerUserId uuid.UUID, displayName string, avatar string) error {
	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}

	log.Info("IncrementScoreCount %v - (%v, %v, %v)", objectId, ownerUserId, displayName, avatar)

	voter := dto.VoterProfile{
		ObjectId:    ownerUserId,
		DisplayName: displayName,
		Avatar:      avatar,
	}

	updateOperator := bson.M{
		"$addToSet": bson.M{"votes": voter},
	}

	options := &coreData.UpdateOptions{}
	options.SetUpsert(true)
	return s.UpdatePost(filter, updateOperator, options)
}

// DecrementScoreCount decrement score of post
func (s PostServiceImpl) DecrementScoreCount(objectId uuid.UUID, ownerUserId uuid.UUID, displayName string, avatar string) error {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}

	voter := dto.VoterProfile{
		ObjectId:    ownerUserId,
		DisplayName: displayName,
		Avatar:      avatar,
	}

	updateOperator := bson.M{
		"$pull": bson.M{"votes": voter},
	}

	options := &coreData.UpdateOptions{}
	options.SetUpsert(true)
	return s.UpdatePost(filter, updateOperator, options)
}

// DisableCommnet
func (s PostServiceImpl) DisableCommnet(OwnerUserId uuid.UUID, objectId uuid.UUID, value bool) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    objectId,
		OwnerUserId: OwnerUserId,
	}

	data := make(map[string]interface{})
	data["disableComments"] = value

	incOperator := coreData.UpdateOperator{
		Set: data,
	}
	return s.UpdatePost(filter, incOperator)
}

// DisableSharing
func (s PostServiceImpl) DisableSharing(OwnerUserId uuid.UUID, objectId uuid.UUID, value bool) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    objectId,
		OwnerUserId: OwnerUserId,
	}

	data := make(map[string]interface{})
	data["disableSharing"] = value

	incOperator := coreData.UpdateOperator{
		Set: data,
	}
	return s.UpdatePost(filter, incOperator)
}

// Increment increment a post field
func (s PostServiceImpl) Increment(objectId uuid.UUID, field string, value int) error {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}

	data := make(map[string]interface{})
	data[field] = value

	incOperator := coreData.IncrementOperator{
		Inc: data,
	}
	return s.UpdatePost(filter, incOperator)
}

// Decrement decrement a post field
func (s PostServiceImpl) Decrement(objectId uuid.UUID, field string, value int) error {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}

	data := make(map[string]interface{})
	data[field] = value

	incOperator := coreData.IncrementOperator{
		Inc: data,
	}
	return s.UpdatePost(filter, incOperator)
}

// IncrementCommentCount increment comment count of post by +1
func (s PostServiceImpl) IncrementCommentCount(objectId uuid.UUID) error {
	return s.Increment(objectId, "commentCounter", 1)
}

// DecrementCommentCount increment comment count of post by -1
func (s PostServiceImpl) DecrementCommentCount(objectId uuid.UUID) error {
	return s.Increment(objectId, "commentCounter", -1)
}

// UpdatePostProfile update the post
func (s PostServiceImpl) UpdatePostProfile(ownerUserId uuid.UUID, ownerDisplayName string, ownerAvatar string) error {
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}

	data := struct {
		OwnerDisplayName string `json:"ownerDisplayName" bson:"ownerDisplayName"`
		OwnerAvatar      string `json:"ownerAvatar" bson:"ownerAvatar"`
	}{
		OwnerDisplayName: ownerDisplayName,
		OwnerAvatar:      ownerAvatar,
	}

	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateManyPost(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePostURLKey update the post URL key
func (s PostServiceImpl) UpdatePostURLKey(postId uuid.UUID, urlKey string) error {
	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: postId,
	}

	data := struct {
		URLKey string `json:"urlKey" bson:"urlKey"`
	}{
		URLKey: urlKey,
	}

	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateManyPost(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}
