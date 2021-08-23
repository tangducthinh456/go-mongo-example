package mongoservice

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"thinhtd4/config"
	"thinhtd4/customlog"
	"thinhtd4/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConnectMGDB struct {
	collection *mongo.Database
	M          *sync.Mutex
}

var connectMG ConnectMGDB

const (
	COLLECTION_USER     = "user"
	COLLECTION_QUESTION = "question"
)

func Init() {
	customlog.Info("initialize mongodb connector")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoURL := "mongodb://" + config.MongoConfig().Host + ":" + config.MongoConfig().Port
	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		customlog.Err(err)
		panic(err)
	}

	connectMG.collection = new(mongo.Database)
	connectMG.M = new(sync.Mutex)

	connectMG.collection = client.Database(config.MongoConfig().Database)

}

func MongoDB() ConnectMGDB {
	return connectMG
}

func (e ConnectMGDB) FindUser(ctx context.Context, username string) (reUser model.UserInfo, er error) {
	e.M.Lock()
	defer e.M.Unlock()

	customlog.Info("Get user " + fmt.Sprintf("%v", username) + " from mongo collection " + COLLECTION_USER)
	filter := bson.D{
		{"username", username},
	}

	er = connectMG.collection.Collection(COLLECTION_USER).FindOne(ctx, filter).Decode(&reUser)
	if er != nil {
		customlog.Err(er)
		return
	}

	return
}

func (e ConnectMGDB) InsertUser(ctx context.Context, user *model.UserInfo) (er error) {
	e.M.Lock()
	defer e.M.Unlock()

	customlog.Info("Insert user " + fmt.Sprintf("%v", user.Username) + " to mongo collection " + COLLECTION_USER)
	reUserId, er := connectMG.collection.Collection(COLLECTION_USER).InsertOne(ctx, user)
	if er != nil {
		customlog.Err(er)
		return
	}

	user.UserID = reUserId.InsertedID.(primitive.ObjectID)
	return
}

func (e ConnectMGDB) UpdateUserQuestionDone(ctx context.Context, userID string, questionID string, answer []int) (er error) {
	e.M.Lock()
	defer e.M.Unlock()

	customlog.Info("Update question done user " + fmt.Sprintf("%v", userID) + " from mongo collection " + COLLECTION_USER)

	objectId, er := primitive.ObjectIDFromHex(userID)
	if er != nil {
		customlog.Err(er)
		return
	}

	filter := bson.M{
		"_id": objectId,
	}

	user := model.UserInfo{}

	er = connectMG.collection.Collection(COLLECTION_USER).FindOne(ctx, filter).Decode(&user)
	if er != nil {
		customlog.Err(er)
		return
	}
	if len(user.ListAnswers) == 0 {
		user.ListAnswers = make(map[string][]int)
	}

	if questionID == "" {
		er = errors.New("Update question err : object can not be null or key can not be empty")
		return
	}
	
	user.ListAnswers[questionID] = answer

	_, er = connectMG.collection.Collection(COLLECTION_USER).UpdateOne(ctx, filter, bson.D{
		{
			"$set", bson.D{
				{
					"listanswers", user.ListAnswers,
				},
			},
		},
	})
	if er != nil {
		customlog.Err(er)
		return
	}

	return
}

func (e ConnectMGDB) UpdateUserHaveDoneSurvey(ctx context.Context, userID string, isDone bool) (er error) {
	e.M.Lock()
	defer e.M.Unlock()

	customlog.Info("Update user " + fmt.Sprintf("%v", userID) + " have done survey " + COLLECTION_USER)

	objectId, er := primitive.ObjectIDFromHex(userID)
	if er != nil {
		customlog.Err(er)
		return
	}

	filter := bson.M{
		"_id": objectId,
	}

	_, er = connectMG.collection.Collection(COLLECTION_USER).UpdateOne(ctx, filter, bson.D{
		{
			"$set", bson.D{
				{
					"havedonesurvey", isDone,
				},
			},
		},
	})

	if er != nil {
		customlog.Err(er)
		return
	}

	return

}

func (e ConnectMGDB) InsertQuestion(ctx context.Context, question *model.QuestionWithAnswer) (er error) {
	e.M.Lock()
	defer e.M.Unlock()

	customlog.Info("Insert question to mongo collection " + COLLECTION_QUESTION)
	reQuestion, er := connectMG.collection.Collection(COLLECTION_QUESTION).InsertOne(ctx, question)
	if er != nil {
		customlog.Err(er)
		return
	}
	question.QuestionID = reQuestion.InsertedID.(primitive.ObjectID)
	return
}

func (e ConnectMGDB) GetQuestionsList(ctx context.Context) (list []model.QuestionWithAnswer, er error) {
	e.M.Lock()
	defer e.M.Unlock()

	customlog.Info("Get question list from mongo collection " + COLLECTION_QUESTION)
	filter := bson.D{
		{},
	}
	findOptions := options.Find()

	cur, er := connectMG.collection.Collection(COLLECTION_QUESTION).Find(ctx, filter, findOptions)
	if er != nil {
		customlog.Err(er)
		return
	}

	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		var elem model.QuestionWithAnswer
		er = cur.Decode(&elem)
		if er != nil {
			customlog.Err(er)
			return
		}
		// elem.QuestionID = cur.
		list = append(list, elem)
	}
	if er = cur.Err(); er != nil {
		customlog.Err(er)
		return
	}
	return
}
