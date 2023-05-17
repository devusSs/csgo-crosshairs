package routes

import (
	"fmt"
	"net/http"

	"github.com/devusSs/crosshairs/api/models"
	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/database"
	"github.com/devusSs/crosshairs/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllUsersRoute(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")

	if fmt.Sprintf("%s", userID) == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Invalid session id."
		resp.SendErrorResponse(c)
		return
	}

	uuidUser, err := uuid.Parse(fmt.Sprintf("%s", userID))
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
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "User could not be found."
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

	email := c.Query("email")

	if email != "" {
		if !utils.IsEmailValid(email) {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Invalid e-mail address provided."
			resp.SendErrorResponse(c)
			return
		}

		user, err := Svc.GetUserByEmail(&database.UserAccount{EMail: email})
		if err != nil {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "User does not exist."
			resp.SendErrorResponse(c)
			return
		}

		var returnUser models.ReturnUserAdmin

		returnUser.CreatedAt = user.CreatedAt
		returnUser.UpdatedAt = user.UpdatedAt
		returnUser.EMail = user.EMail
		returnUser.Role = user.Role
		returnUser.VerifiedMail = user.VerifiedMail
		returnUser.RegisterIP = user.RegisterIP
		returnUser.LoginIP = user.LoginIP
		returnUser.LastLogin = user.LastLogin
		returnUser.CrosshairsRegistered = user.CrosshairsRegistered

		resp := responses.SuccessResponse{
			Code: http.StatusOK,
			Data: returnUser,
		}
		resp.SendSuccessReponse(c)
		return
	}

	users, err := Svc.GetAllUsers()
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	var usersReturn models.MultipleUsersAdmin

	for _, u := range users {
		var user models.ReturnUserAdmin
		user.CreatedAt = u.CreatedAt
		user.UpdatedAt = u.UpdatedAt
		user.EMail = u.EMail
		user.Role = u.Role
		user.VerifiedMail = u.VerifiedMail
		user.RegisterIP = u.RegisterIP
		user.LoginIP = u.LoginIP
		user.LastLogin = u.LastLogin
		user.CrosshairsRegistered = u.CrosshairsRegistered
		usersReturn.Users = append(usersReturn.Users, user)
	}

	resp := responses.SuccessResponse{
		Code: http.StatusOK,
		Data: usersReturn,
	}
	resp.SendSuccessReponse(c)
}

func GetAllCrosshairsRoute(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")

	if fmt.Sprintf("%s", userID) == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Invalid session id."
		resp.SendErrorResponse(c)
		return
	}

	uuidUser, err := uuid.Parse(fmt.Sprintf("%s", userID))
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
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "User could not be found."
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

	crosshairsDB, err := Svc.GetAllCrosshairs()
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	var crosshairs models.GetMultipleCrosshairs

	email := c.Query("email")

	if email != "" {
		if !utils.IsEmailValid(email) {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Invalid e-mail address provided."
			resp.SendErrorResponse(c)
			return
		}

		user, err := Svc.GetUserByEmail(&database.UserAccount{EMail: email})
		if err != nil {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "User does not exist."
			resp.SendErrorResponse(c)
			return
		}

		for _, ch := range crosshairsDB {
			var crosshair models.Crosshair
			if ch.RegistrantID == user.ID {
				crosshair.Added = ch.CreatedAt
				crosshair.Code = ch.Code
				crosshair.Note = ch.Note
				crosshairs.Crosshairs = append(crosshairs.Crosshairs, crosshair)
			}
		}

		resp := responses.SuccessResponse{
			Code: http.StatusOK,
			Data: crosshairs,
		}
		resp.SendSuccessReponse(c)
		return
	}

	for _, ch := range crosshairsDB {
		crosshair := models.Crosshair{
			Added: ch.CreatedAt,
			Code:  ch.Code,
			Note:  ch.Note,
		}
		crosshairs.Crosshairs = append(crosshairs.Crosshairs, crosshair)
	}

	resp := responses.SuccessResponse{
		Code: http.StatusOK,
		Data: crosshairs,
	}
	resp.SendSuccessReponse(c)
}
