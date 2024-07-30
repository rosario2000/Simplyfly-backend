package partner

import (
	"context"
	"github.com/labstack/gommon/log"
	"net/http"
	"simplifly/common/model"
	"simplifly/common/mongo"
	"simplifly/common/request"
	"simplifly/internal/api"
)

func GetUserDetails(ctx context.Context, req *http.Request, api api.API) (_ interface{}, err error) {
	userName := req.Header.Get("username")
	password := req.Header.Get("password")
	return mongo.GetUserInfo(context.Background(), model.GetUserDetailsRequest{
		UserName: userName,
		Password: password,
	})
}

func UpsertUser(ctx context.Context, req *http.Request, api api.API) (_ interface{}, err error) {

	var userInfo model.UserInfo
	err = request.DecodeBody(req, &userInfo)
	if err != nil {
		log.Error("UpsertUser : Error in decoding request - ", err)
		return
	}
	if userInfo.Flights == nil {
		userInfo.Flights = []string{}
	}

	err = mongo.UpsertUser(ctx, userInfo)
	if err != nil {
		log.Error("UpsertUser : Error in inserting / updating data to userInfoDb - ", err)
		return
	}

	return

}

func GetAllUsers(ctx context.Context, req *http.Request, api api.API) (users []model.UserInfo, err error) {
	users, err = mongo.GetAllUsers()
	return
}

func BookFlight(ctx context.Context, req *http.Request, api api.API) (_ interface{}, err error) {
	var bookFlightRequest model.BookFlightReq
	err = request.DecodeBody(req, &bookFlightRequest)
	if err != nil {
		log.Error("BookFlight : Error in decoding request - ", err)
		return
	}

	userInfo, err := mongo.GetUserInfo(context.Background(), model.GetUserDetailsRequest{
		UserName: bookFlightRequest.UserName,
		Password: bookFlightRequest.Password,
	})
	if err != nil {
		return
	}
	if userInfo.Flights == nil {
		flightsArr := []string{bookFlightRequest.FlightId}
		userInfo.Flights = flightsArr
	} else {
		userInfo.Flights = append(userInfo.Flights, bookFlightRequest.FlightId)
	}

	err = mongo.UpsertUser(ctx, userInfo)

	return
}
