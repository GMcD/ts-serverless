package models

type CollectivesModel struct {
	CollectiveIds string `json:"collectiveIds" bson:"collectiveIds"`
	Name          string `json:"Name" bson:"Name", validate:"max=50"`
	Avatar        string `json:"avatar" bson:"avatar", validate:"max=500"`
	Banner        string `json:"banner", bson:"banner", validate:"max=500"`
	CreatedDate   int64  `json:"created_date" bson:"created_date"`
	LastUpdated   int64  `json:"last_updated" bson:"last_updated"`
	VoteCount     int64  `json:"voteCount" bson:"voteCount"`
	ShareCount    int64  `json:"shareCount" bson:"shareCount"`
	FollowerCount int64  `json:"followerCount" bson:"followerCount"`
	PostCount     int64  `json:"postCount" bson:"postCount"`
}
