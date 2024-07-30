package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"simplifly/common/model"
)

func getUserCollection() *mongo.Collection {
	db := "simplyfly"
	return Instance().Database(db).Collection("userInfoDb")
}

func UpsertUser(ctx context.Context, user model.UserInfo) (err error) {
	col := getUserCollection()
	filter := bson.M{"userName": user.UserName, "password": user.Password}
	// Update document with upsert option
	update := bson.M{
		"$set": user,
	}
	updateOptions := options.Update().SetUpsert(true)

	_, err = col.UpdateOne(context.TODO(), filter, update, updateOptions)
	if err != nil {
		return
	}
	return
}

func GetUserInfo(ctx context.Context, getUserDetailsReq model.GetUserDetailsRequest) (userInfo model.UserInfo, err error) {
	col := getUserCollection()
	filter := bson.M{"userName": getUserDetailsReq.UserName}

	err = col.FindOne(context.Background(), filter).Decode(&userInfo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.New("User Not found")
			return
		}
		return
	}

	if userInfo.Password != getUserDetailsReq.Password {
		err = errors.New("Invalid Password")
	}

	// User found in the database
	return

}

func GetAllUsers() (users []model.UserInfo, err error) {
	col := getUserCollection()

	cursor, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		return
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		return
	}

	return
}
