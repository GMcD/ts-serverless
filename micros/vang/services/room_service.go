package service

import (
	"fmt"

	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/config"
	coreData "github.com/red-gold/telar-core/data"
	repo "github.com/red-gold/telar-core/data"
	"github.com/red-gold/telar-core/data/mongodb"
	mongoRepo "github.com/red-gold/telar-core/data/mongodb"
	"github.com/red-gold/telar-core/utils"
	dto "github.com/red-gold/ts-serverless/micros/vang/dto"
)

// RoomService handlers with injected dependencies
type RoomServiceImpl struct {
	RoomRepo repo.Repository
}

// NewRoomService initializes RoomService's dependencies and create new RoomService struct
func NewRoomService(db interface{}) (RoomService, error) {

	roomService := &RoomServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		roomService.RoomRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}

	return roomService, nil
}

// SaveRoom save the room
func (s RoomServiceImpl) SaveRoom(room *dto.Room) error {

	if room.ObjectId == uuid.Nil {
		var uuidErr error
		room.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if room.CreatedDate == 0 {
		room.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.RoomRepo.Save(vangRoomCollectionName, room)

	return result.Error
}

// FindOneRoom get one room
func (s RoomServiceImpl) FindOneRoom(filter interface{}) (*dto.Room, error) {

	result := <-s.RoomRepo.FindOne(vangRoomCollectionName, filter)
	if result.Error() != nil {
		if result.Error() == repo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Error()
	}

	var roomResult dto.Room
	errDecode := result.Decode(&roomResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.Room")
	}
	return &roomResult, nil
}

// FindRoomList get all rooms by filter
func (s RoomServiceImpl) FindRoomList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Room, error) {

	result := <-s.RoomRepo.Find(vangRoomCollectionName, filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var roomList []dto.Room
	for result.Next() {
		var room dto.Room
		errDecode := result.Decode(&room)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Room")
		}
		roomList = append(roomList, room)
	}

	return roomList, nil
}

// FindByOwnerUserId find by owner user id
func (s RoomServiceImpl) FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Room, error) {
	sortMap := make(map[string]int)
	sortMap["createdDate"] = -1
	filter := struct {
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		OwnerUserId: ownerUserId,
	}
	return s.FindRoomList(filter, 0, 0, sortMap)
}

// FindById find by room id
func (s RoomServiceImpl) FindById(objectId uuid.UUID) (*dto.Room, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}
	return s.FindOneRoom(filter)
}

// UpdateRoom update the room
func (s RoomServiceImpl) UpdateRoom(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error {

	result := <-s.RoomRepo.Update(vangRoomCollectionName, filter, data, opts...)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateRoom update the room
func (s RoomServiceImpl) UpdateRoomById(data *dto.Room) error {
	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: data.ObjectId,
	}

	updateOperator := coreData.UpdateOperator{
		Set: data,
	}
	err := s.UpdateRoom(filter, updateOperator)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRoom delete room by filter
func (s RoomServiceImpl) DeleteRoom(filter interface{}) error {

	result := <-s.RoomRepo.Delete(vangRoomCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteRoom delete room by ownerUserId and roomId
func (s RoomServiceImpl) DeleteRoomByOwner(ownerUserId uuid.UUID, roomId uuid.UUID) error {

	filter := struct {
		ObjectId    uuid.UUID `json:"objectId" bson:"objectId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		ObjectId:    roomId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteRoom(filter)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManyRoom delete many room by filter
func (s RoomServiceImpl) DeleteManyRoom(filter interface{}) error {

	result := <-s.RoomRepo.Delete(vangRoomCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateRoomIndex create index for room search.
func (s RoomServiceImpl) CreateRoomIndex(indexes map[string]interface{}) error {
	result := <-s.RoomRepo.CreateIndex(vangRoomCollectionName, indexes)
	return result
}

// GetRoomByRoomId get all room by room ID
func (s RoomServiceImpl) GetRoomByRoomId(roomId *uuid.UUID, sortBy string, page int64) ([]dto.Room, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems

	filter := make(map[string]interface{})

	if roomId != nil {
		filter["roomId"] = *roomId
	}

	result, err := s.FindRoomList(filter, limit, skip, sortMap)

	return result, err
}

// DeleteRoomByRoomId delete room by room id
func (s RoomServiceImpl) DeleteRoomByRoomId(ownerUserId uuid.UUID, roomId uuid.UUID) error {

	filter := struct {
		PostId      uuid.UUID `json:"roomId" bson:"roomId"`
		OwnerUserId uuid.UUID `json:"ownerUserId" bson:"ownerUserId"`
	}{
		PostId:      roomId,
		OwnerUserId: ownerUserId,
	}
	err := s.DeleteManyRoom(filter)
	if err != nil {
		return err
	}
	return nil
}

// FindOneRoomByMembers find one room by members
func (s RoomServiceImpl) FindOneRoomByMembers(userIds []string, roomType int8) (*dto.Room, error) {

	include := make(map[string]interface{})
	include["$in"] = userIds

	filter := make(map[string]interface{})
	filter["members"] = include
	filter["type"] = roomType

	return s.FindOneRoom(filter)
}
