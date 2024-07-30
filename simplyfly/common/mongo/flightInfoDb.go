package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"simplifly/common/model"
)

func GetFlightCollection() *mongo.Collection {
	db := "simplyfly"
	return Instance().Database(db).Collection("flightInfoDb")
}

func UpsertFlight(ctx context.Context, flight model.Flight) (err error) {
	col := GetFlightCollection()
	filter := bson.M{"id": flight.Id}
	// Update document with upsert option
	update := bson.M{
		"$set": flight,
	}
	updateOptions := options.Update().SetUpsert(true)

	_, err = col.UpdateOne(context.TODO(), filter, update, updateOptions)
	if err != nil {
		return
	}
	return
}

func GetFlightDetails(ctx context.Context, flightId string) (flightDetails model.Flight, err error) {
	col := GetFlightCollection()
	filter := bson.M{"id": flightId}
	// Update document with upsert option

	cur := col.FindOne(context.TODO(), filter)
	err = cur.Err()
	if err != nil {
		//if err == mongo.ErrNoDocuments {
		//	return flightDetails, err
		//}
		return
	}

	if err = cur.Decode(&flightDetails); err != nil {
		return flightDetails, err
	}

	return flightDetails, nil
}

func GetAllFlights(fromCity string, toCity string) (flights []model.Flight, err error) {
	col := GetFlightCollection()

	var filter interface{}
	if len(fromCity) > 0 && len(toCity) > 0 {
		filter = bson.M{
			"from": fromCity,
			"to":   toCity,
		}
	}

	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return
	}

	if err = cursor.All(context.TODO(), &flights); err != nil {
		return
	}

	return
}

func DeleteFlight(flight model.Flight) (err error) {
	col := GetFlightCollection()
	filter := bson.M{"id": flight.Id}

	_, err = col.DeleteOne(context.TODO(), filter)
	if err != nil {
		return
	}

	return
}
