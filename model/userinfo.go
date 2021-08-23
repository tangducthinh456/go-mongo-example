package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	UserID          primitive.ObjectID `bson:"_id,omitempty"`
	Username        string
	EncryptPassword string
	HaveDoneSurvey  bool
	ListAnswers     map[string][]int
}

type UserDto struct {
	// UserID   string `json:"user_id"`
	Username string `form:"username"`
	Password string `form:"password"`
}

