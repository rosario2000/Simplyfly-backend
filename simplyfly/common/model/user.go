package model

type UserInfo struct {
	Name     string   `json:"name" bson:"name"`
	Email    string   `json:"email" bson:"email"`
	Contact  string   `json:"contact" bson:"contact"`
	UserName string   `json:"userName" bson:"userName"`
	Password string   `json:"password" bson:"password"`
	Flights  []string `json:"flights" bson:"flights"`
	Active   bool     `json:"active" bson:"active"`
}

type GetUserDetailsRequest struct {
	UserName string `json:"userName" bson:"userName"`
	Password string `json:"password" bson:"password"`
}

type BookFlightReq struct {
	UserName string `json:"userName" `
	Password string `json:"password" `
	FlightId string `json:"flightId"`
}
