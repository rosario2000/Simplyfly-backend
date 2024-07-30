package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"simplifly/partner"
)

func (h *Handler) GetUserDetailsHandler(c echo.Context) (err error) {

	req := c.Request()
	ctx := req.Context()

	res, err := partner.GetUserDetails(ctx, req, h.API)
	if err != nil {
		return
	}
	// send response payload
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpsertUserHandler(c echo.Context) (err error) {
	req := c.Request()
	ctx := req.Context()

	res, err := partner.UpsertUser(ctx, req, h.API)
	if err != nil {
		return
	}
	// send response payload
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAllUsersHandler(c echo.Context) (err error) {
	req := c.Request()
	ctx := req.Context()

	res, err := partner.GetAllUsers(ctx, req, h.API)
	if err != nil {
		return
	}
	// send response payload
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) BookFlightByUserHandler(c echo.Context) (err error) {
	req := c.Request()
	ctx := req.Context()

	res, err := partner.BookFlight(ctx, req, h.API)
	if err != nil {
		return
	}
	// send response payload
	return c.JSON(http.StatusOK, res)
}
