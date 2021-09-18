package models

import (
	uuid "github.com/gofrs/uuid"
)

// func (s *CommentModel) SetField(value string) error {
// 	if len(value) > 1500 {
// 		return error.New("Maximun length of Comment Text exceeded")
// 	}
// 	s.Text =< value
// 	return nil
// }

type CommentModel struct {
	ObjectId         uuid.UUID `json:"objectId"`
	Score            int64     `json:"score"`
	OwnerUserId      uuid.UUID `json:"ownerUserId"`
	OwnerDisplayName string    `json:"ownerDisplayName"`
	OwnerAvatar      string    `json:"ownerAvatar"`
	PostId           uuid.UUID `json:"postId"`
	Text             string    `json:"text, validate:"max=1500"`
	Deleted          bool      `json:"deleted"`
	DeletedDate      int64     `json:"deletedDate"`
	CreatedDate      int64     `json:"created_date"`
	LastUpdated      int64     `json:"last_updated"`
}
