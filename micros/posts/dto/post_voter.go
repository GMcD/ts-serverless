package dto

import uuid "github.com/gofrs/uuid"

type VoterProfile struct {
	ObjectId      uuid.UUID `json:"objectId" bson:"objectId"`
	FullName      string    `json:"fullName" bson:"fullName"`
	DisplayName   string    `json:"displayName" bson:"displayName"`
	SocialName    string    `json:"socialName" bson:"socialName"`
	Avatar        string    `json:"avatar" bson:"avatar"`
	FollowCount   int64     `json:"followCount" bson:"followCount"`
	FollowerCount int64     `json:"followerCount" bson:"followerCount"`
}
