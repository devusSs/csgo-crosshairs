package routes

import (
	"fmt"
	"net/http"

	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/database"
	"github.com/devusSs/crosshairs/stats"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetTotalStatsRoute(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("user") == nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "You are currently not logged in."
		resp.SendErrorResponse(c)
		return
	}

	uuidUser, err := uuid.Parse(fmt.Sprintf("%s", session.Get("user")))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse uuid."
		resp.SendErrorResponse(c)
		return
	}

	user, err := Svc.GetUserByUID(&database.UserAccount{ID: uuidUser})
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	if user.Role != "admin" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "You are not an admin."
		resp.SendErrorResponse(c)
		return
	}

	total, err := stats.GetStatsAllTime(Svc)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusOK
	resp.Data = total
	resp.SendSuccessReponse(c)
}

func Get24HourStatsRoute(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("user") == nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "You are currently not logged in."
		resp.SendErrorResponse(c)
		return
	}

	uuidUser, err := uuid.Parse(fmt.Sprintf("%s", session.Get("user")))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse uuid."
		resp.SendErrorResponse(c)
		return
	}

	user, err := Svc.GetUserByUID(&database.UserAccount{ID: uuidUser})
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	if user.Role != "admin" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "You are not an admin."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusOK
	resp.Data = stats.GetStats24Hours()
	resp.SendSuccessReponse(c)
}

func GetSystemStatsRoute(c *gin.Context) {
	// No need to check the session here since the user is a vaidated engineer.
	data, err := stats.CollectSystemStats(Svc, StorageSvc)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusOK
	resp.Data = data
	resp.SendSuccessReponse(c)
}
