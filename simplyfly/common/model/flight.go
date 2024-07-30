package model

import "time"

type Flight struct {
	Id        string    `json:"id" bson:"id"`
	From      string    `json:"from" bson:"from"`
	To        string    `json:"to" bson:"to"`
	Arrival   time.Time `json:"arrival" bson:"arrival"`
	Departure time.Time `json:"departure" bson:"departure"`
	Gate      int       `json:"gate" bson:"gate"`
	Cancelled bool      `json:"cancelled" bson:"cancelled"`
}

type UpdateFlightRequest struct {
	Id        string `json:"id"`
	Delay     int    `json:"delay"`
	Gate      int    `json:"gate"`
	Cancelled bool   `json:"cancelled"`
}
