package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// import "time"

type Question struct {
	// QuestionID primitive.ObjectID `bson:"_id,omitempty" json:"question_id"`
	Question   string   `json:"question"`
	ListChoose []string `json:"list_choose"`
}

type QuestionWithAnswer struct {
	QuestionID primitive.ObjectID `bson:"_id,omitempty" json:"question_id"`
	Question
	Answers []int `json:"answers"`
}

type QuestionDto struct {
	QuestionID primitive.ObjectID `json:"question_id"`
	Question
}

type AnswerUpdate struct {
	QuestionID string `json:"question_id"`
	Answer     []int  `json:"answer"`
}
