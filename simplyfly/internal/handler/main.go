package handler

import (
	"simplifly/internal/api"
)

type Handler struct {
	API    api.API
	IsTest bool
}
