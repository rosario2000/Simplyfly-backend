package partner

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"simplifly/common/constants"
	"simplifly/common/kafka/kafkaProducer"
	"simplifly/common/model"
	"simplifly/common/mongo"
	"simplifly/common/request"
	"simplifly/internal/api"
	"time"
)

func RegisterFlight(ctx context.Context, req *http.Request, api api.API) (res interface{}, err error) {
	var flight model.Flight
	err = request.DecodeBody(req, &flight)
	if err != nil {
		log.Error("RegisterFlight : Error in decoding request - ", err)
		return
	}

	err = mongo.UpsertFlight(ctx, flight)
	if err != nil {
		log.Error("RegisterFlight : Error in inserting / updating data to userInfoDb - ", err)
		return
	}

	return model.Response{
		Status:     constants.Success,
		Message:    "Flight registered successfully",
		HttpStatus: http.StatusOK,
	}, nil

}

func UpdateFlight(ctx context.Context, req *http.Request, api api.API) (res interface{}, err error) {
	var updateFlightRequest model.UpdateFlightRequest
	err = request.DecodeBody(req, &updateFlightRequest)
	if err != nil {
		log.Error("UpdateFlight : Error in decoding request - ", err)
		return
	}

	flightDetails, err := GetFlightDetails(ctx, updateFlightRequest.Id, api)
	if err != nil {
		log.Error("UpdateFlight : Error in retrieving flight details  - ", err)
		return
	}

	if updateFlightRequest.Cancelled == true {
		flightDetails.Cancelled = true
		err = mongo.UpsertFlight(ctx, flightDetails)
		if err != nil {
			return
		}

		// send update to all users
		sendUpdatesErr := SendUpdates(updateFlightRequest, flightDetails, constants.Cancellation)
		if sendUpdatesErr != nil {
			log.Error("UpdateFlight : Error in sending updates - ", err)
			return
		}

		return model.Response{
			Status:     constants.Success,
			Message:    "Flight Updated successfully",
			HttpStatus: http.StatusOK,
		}, nil
	}

	if updateFlightRequest.Delay > 0 {
		delayDuration := time.Minute * time.Duration(updateFlightRequest.Delay)
		flightDetails.Arrival = flightDetails.Arrival.Add(delayDuration)
		flightDetails.Departure = flightDetails.Departure.Add(delayDuration)
		err = mongo.UpsertFlight(ctx, flightDetails)
		if err != nil {
			return
		}

		// send update to all users with this flight
		sendUpdatesErr := SendUpdates(updateFlightRequest, flightDetails, constants.Delay)
		if sendUpdatesErr != nil {
			log.Error("UpdateFlight : Error in sending updates - ", sendUpdatesErr)
		}
	}

	if updateFlightRequest.Gate > 0 && updateFlightRequest.Gate != flightDetails.Gate {
		flightDetails.Gate = updateFlightRequest.Gate
		err = mongo.UpsertFlight(ctx, flightDetails)
		// send update to all users with this flight
		sendUpdatesErr := SendUpdates(updateFlightRequest, flightDetails, constants.GateChange)
		if sendUpdatesErr != nil {
			log.Error("UpdateFlight : Error in sending updates - ", sendUpdatesErr)
		}
	}

	return model.Response{
		Status:     constants.Success,
		Message:    "Flight Updated successfully",
		HttpStatus: http.StatusOK,
	}, nil
}

func GetFlightDetails(ctx context.Context, flightId string, api api.API) (flight model.Flight, err error) {

	flight, err = mongo.GetFlightDetails(ctx, flightId)
	if err != nil {
		log.Error("GetFlightDetails: Error in getting Flight details from mongo - ", err)
		return
	}

	return
}

func GetAllFlights(ctx context.Context, req *http.Request, api api.API) (_ interface{}, err error) {

	var flights []model.Flight

	fromCity := req.Header.Get("from_city")
	toCity := req.Header.Get("to_city")

	allFlights, err := mongo.GetAllFlights(fromCity, toCity)

	currentTime := time.Now()

	date := req.Header.Get("search_date")
	for _, flight := range allFlights {
		if currentTime.After(flight.Departure) {
			go func() {
				flightDeletionErr := mongo.DeleteFlight(flight)
				if err != nil {
					log.Error("GetAllFlights: Error in getting Flight details from redis - ", flightDeletionErr, " flightId - ", flight.Id)
				}
			}()
		} else {
			flightDate := flight.Departure.Format("2006-01-02")

			if len(date) > 0 && flightDate != date {
				continue
			}

			flights = append(flights, flight)
		}
	}

	return flights, nil
}

func SendUpdates(flightUpdateReq model.UpdateFlightRequest, flight model.Flight, updateType string) (err error) {

	var updateMessage string
	switch updateType {
	case constants.Cancellation:
		updateMessage = fmt.Sprintf("Your Indigo flight %s from %s to %s has been cancelled due to unforseen circumstances.", flight.Id, flight.From, flight.To)
	case constants.Delay:
		updateMessage = fmt.Sprintf("Your Indigo flight %s from %s to %s has been delayed by %v minutes. The expected departure time is %v UTC", flight.Id, flight.From, flight.To, flight.Departure)
	case constants.GateChange:
		updateMessage = fmt.Sprintf("Your Indigo flight %s from %s to %s will be departing from Gate %v.", flight.Id, flight.From, flight.To, flight.Gate)
	}
	updateMessage += "We apologize for the inconvenience"

	err = kafkaProducer.SendUpdateToUser(context.Background(), updateMessage)
	if err != nil {
		log.Error("Error in SendUpdates while sending message to kafka: ", err, " flightId : ", flight.Id)

	}

	return
}
