package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"simplifly/partner"
)

func (h *Handler) RegisterFlightHandler(c echo.Context) (err error) {

	req := c.Request()
	ctx := req.Context()

	res, err := partner.RegisterFlight(ctx, req, h.API)
	if err != nil {
		return
	}
	// send response payload
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateFlightHandler(c echo.Context) (err error) {

	req := c.Request()
	ctx := req.Context()

	res, err := partner.UpdateFlight(ctx, req, h.API)
	if err != nil {
		return
	}
	// send response payload
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetFlightDetailsHandler(c echo.Context) (err error) {

	req := c.Request()
	ctx := req.Context()

	flightId := c.Param("id")

	res, err := partner.GetFlightDetails(ctx, flightId, h.API)
	if err != nil {
		return
	}
	// send response payload
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllFlightsHandler(c echo.Context) (err error) {

	req := c.Request()
	ctx := req.Context()

	res, err := partner.GetAllFlights(ctx, req, h.API)
	if err != nil {
		return
	}
	// send response payload
	return c.JSON(http.StatusOK, res)
}
