package routes

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/devusSs/crosshairs/api/models"
	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	lenNoteMin       = 3
	shareCodePattern = `^CSGO-[A-Za-z0-9]{5}-[A-Za-z0-9]{5}-[A-Za-z0-9]{5}-[A-Za-z0-9]{5}-[A-Za-z0-9]{5}$`
	crosshairsMax    = 20
)

func AddCrosshairRoute(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")

	if fmt.Sprintf("%s", userID) == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Missing user id."
		resp.SendErrorResponse(c)
		return
	}

	userUID, err := uuid.Parse(fmt.Sprintf("%s", userID))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse user id."
		resp.SendErrorResponse(c)
		return
	}

	var addCrosshair models.AddCrosshair

	if err := c.BindJSON(&addCrosshair); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid JSON body provided."
		resp.SendErrorResponse(c)
		return
	}

	user, err := Svc.GetUserByUID(&database.UserAccount{ID: userUID})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not find user."
		resp.SendErrorResponse(c)
		return
	}

	if user.CrosshairsRegistered > crosshairsMax {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Already added maximum number of crosshairs."
		resp.SendErrorResponse(c)
		return
	}

	matched, err := regexp.Match(shareCodePattern, []byte(addCrosshair.Code))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse crosshair code."
		resp.SendErrorResponse(c)
		return
	}

	if !matched {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid crosshair code provided."
		resp.SendErrorResponse(c)
		return
	}

	if len(addCrosshair.Note) < lenNoteMin {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = fmt.Sprintf("Crosshair note needs to be at least %d characters long.", lenNoteMin)
		resp.SendErrorResponse(c)
		return
	}

	crosshair := &database.Crosshair{
		RegistrantID: userUID,
		Code:         addCrosshair.Code,
		Note:         addCrosshair.Note,
		RegisterIP:   c.Request.Header.Get("X-Forwarded-For"),
	}

	_, err = Svc.AddCrosshair(crosshair)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	user, err = Svc.UpdateUserCrosshairCount(&database.UserAccount{ID: userUID})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{
		Code: http.StatusCreated,
		Data: responses.CrosshairResponse{
			Status:      "Successfully added crosshair",
			CHsOnRecord: user.CrosshairsRegistered + 1,
		},
	}
	resp.SendSuccessReponse(c)
}

func GetAllCrosshairsFromUserRoute(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")

	if fmt.Sprintf("%s", userID) == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Missing user id."
		resp.SendErrorResponse(c)
		return
	}

	userUID, err := uuid.Parse(fmt.Sprintf("%s", userID))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse user id."
		resp.SendErrorResponse(c)
		return
	}

	crosshairs, err := Svc.GetAllCrosshairsFromUser(userUID)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	crosshairCode := c.Query("code")
	if crosshairCode != "" {
		matched, err := regexp.Match(shareCodePattern, []byte(crosshairCode))
		if err != nil {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Could not parse crosshair code."
			resp.SendErrorResponse(c)
			return
		}

		if !matched {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Invalid crosshair code provided."
			resp.SendErrorResponse(c)
			return
		}

		for _, ch := range crosshairs {
			var crosshair models.Crosshair
			if ch.Code == crosshairCode {
				crosshair.ID = ch.ID
				crosshair.Added = ch.CreatedAt
				crosshair.Code = ch.Code
				crosshair.Note = ch.Note

				resp := responses.SuccessResponse{}
				resp.Code = http.StatusOK
				resp.Data = crosshair
				resp.SendSuccessReponse(c)
				return
			}
		}

		resp := responses.ErrorResponse{}
		resp.Code = http.StatusNotFound
		resp.Error.ErrorCode = "not_found"
		resp.Error.ErrorMessage = "No matching crosshair found."
		resp.SendErrorResponse(c)
		return
	}

	startDate := c.Query("start")
	endDate := c.Query("end")

	if startDate != "" && endDate != "" {
		startTime, err := time.Parse(time.DateOnly, startDate)
		if err != nil {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Invalid start date specified."
			resp.SendErrorResponse(c)
			return
		}

		endTime, err := time.Parse(time.DateOnly, endDate)
		if err != nil {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Invalid start date specified."
			resp.SendErrorResponse(c)
			return
		}

		var returnCrosshairs []models.Crosshair

		for _, ch := range crosshairs {
			var crosshair models.Crosshair
			if ch.CreatedAt.After(startTime) && ch.CreatedAt.Before(endTime) {
				crosshair.ID = ch.ID
				crosshair.Added = ch.CreatedAt
				crosshair.Code = ch.Code
				crosshair.Note = ch.Note
				returnCrosshairs = append(returnCrosshairs, crosshair)
			}
		}

		if len(returnCrosshairs) > 0 {
			resp := responses.SuccessResponse{
				Code: http.StatusOK,
				Data: models.GetMultipleCrosshairs{Crosshairs: returnCrosshairs},
			}
			resp.SendSuccessReponse(c)
			return
		}

		resp := responses.ErrorResponse{}
		resp.Code = http.StatusNotFound
		resp.Error.ErrorCode = "not_found"
		resp.Error.ErrorMessage = "No matching crosshairs found."
		resp.SendErrorResponse(c)
		return
	}

	var returnCrosshairs []models.Crosshair

	for _, ch := range crosshairs {
		var crosshair models.Crosshair
		crosshair.ID = ch.ID
		crosshair.Added = ch.CreatedAt
		crosshair.Code = ch.Code
		crosshair.Note = ch.Note
		returnCrosshairs = append(returnCrosshairs, crosshair)
	}

	resp := responses.SuccessResponse{
		Code: http.StatusOK,
		Data: models.GetMultipleCrosshairs{Crosshairs: returnCrosshairs},
	}
	resp.SendSuccessReponse(c)
}

func DeleteAllCrosshairsFromUserRoute(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")

	if fmt.Sprintf("%s", userID) == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Missing user id."
		resp.SendErrorResponse(c)
		return
	}

	userUID, err := uuid.Parse(fmt.Sprintf("%s", userID))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse user id."
		resp.SendErrorResponse(c)
		return
	}

	if err := Svc.DeleteAllCrosshairsFromUser(userUID); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something, went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{
		Code: http.StatusNoContent,
	}
	resp.SendSuccessReponse(c)
}

func DeleteOneCrosshairFromUserRoute(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")

	if fmt.Sprintf("%s", userID) == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Missing user id."
		resp.SendErrorResponse(c)
		return
	}

	userUID, err := uuid.Parse(fmt.Sprintf("%s", userID))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse user id."
		resp.SendErrorResponse(c)
		return
	}

	code := c.Param("code")

	if code == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Did not specify crosshair code."
		resp.SendErrorResponse(c)
		return
	}

	if err := Svc.DeleteCrosshairFromUserByCode(userUID, code); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusNotFound
		resp.Error.ErrorCode = "not_found"
		resp.Error.ErrorMessage = "Could not find crosshair."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{
		Code: http.StatusNoContent,
	}
	resp.SendSuccessReponse(c)
}
