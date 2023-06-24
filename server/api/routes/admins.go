package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/devusSs/crosshairs/api/models"
	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/database"
	"github.com/devusSs/crosshairs/logging"
	"github.com/devusSs/crosshairs/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type zapLogFormat struct {
	Level     string  `json:"level"`
	Timestamp float64 `json:"ts"`
	Message   string  `json:"msg"`
	Status    int     `json:"status"`
	Method    string  `json:"method"`
	Path      string  `json:"path"`
	Query     string  `json:"query"`
	IP        string  `json:"ip"`
	UserAgent string  `json:"user-agent"`
	Latency   float64 `json:"latency"`
	Time      string  `json:"time"`
}

func GetAllUsersRoute(c *gin.Context) {
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
			errString := database.CheckDatabaseError(err)
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = errString
			resp.SendErrorResponse(c)
			return
		}

		var returnUser models.ReturnUserAdmin

		if user.AvatarURL == "" {
			defaultAvatar, err := StorageSvc.GetUserProfilePictureLink("sample")
			if err != nil {
				resp := responses.ErrorResponse{}
				resp.Code = http.StatusInternalServerError
				resp.Error.ErrorCode = "internal_error"
				resp.Error.ErrorMessage = "Something went wrong, sorry."
				resp.SendErrorResponse(c)
				return
			}

			// Relevant for Docker only.
			defaultAvatar = strings.Replace(defaultAvatar, "http://minio:", fmt.Sprintf("http://%s:", "localhost"), 1)

			returnUser.AvatarURL = defaultAvatar
		}

		// Relevant for Docker only.
		returnUser.AvatarURL = strings.Replace(returnUser.AvatarURL, "http://minio:", fmt.Sprintf("http://%s:", "localhost"), 1)

		returnUser.ID = user.ID
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
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	var usersReturn models.MultipleUsersAdmin

	for _, u := range users {
		var user models.ReturnUserAdmin
		user.ID = u.ID
		user.CreatedAt = u.CreatedAt
		user.UpdatedAt = u.UpdatedAt
		user.EMail = u.EMail
		user.Role = u.Role
		user.VerifiedMail = u.VerifiedMail
		user.RegisterIP = u.RegisterIP
		user.LoginIP = u.LoginIP
		user.LastLogin = u.LastLogin
		user.CrosshairsRegistered = u.CrosshairsRegistered
		user.AvatarURL = u.AvatarURL

		if user.AvatarURL == "" {
			defaultAvatar, err := StorageSvc.GetUserProfilePictureLink("sample")
			if err != nil {
				resp := responses.ErrorResponse{}
				resp.Code = http.StatusInternalServerError
				resp.Error.ErrorCode = "internal_error"
				resp.Error.ErrorMessage = "Something went wrong, sorry."
				resp.SendErrorResponse(c)
				return
			}

			// Relevant for Docker only.
			defaultAvatar = strings.Replace(defaultAvatar, "http://minio:", fmt.Sprintf("http://%s:", "localhost"), 1)

			user.AvatarURL = defaultAvatar
		}

		// Relevant for Docker only.
		user.AvatarURL = strings.Replace(user.AvatarURL, "http://minio:", fmt.Sprintf("http://%s:", "localhost"), 1)

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

	crosshairsDB, err := Svc.GetAllCrosshairs()
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
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
			errString := database.CheckDatabaseError(err)
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = errString
			resp.SendErrorResponse(c)
			return
		}

		for _, ch := range crosshairsDB {
			var crosshair models.Crosshair
			if ch.RegistrantID == user.ID {
				crosshair.ID = ch.ID
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
			ID:    ch.ID,
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

func GetAllEventsOrByTypeRoute(c *gin.Context) {
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

	limit := c.Query("limit")
	eventType := c.Query("type")

	if eventType != "" {
		if eventType != "user_registered" && eventType != "user_password_change" && eventType != "user_uploaded_avatar" {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Specified invalid event types."
			resp.SendErrorResponse(c)
			return
		}

		if limit != "" {
			limitInt, err := strconv.Atoi(limit)
			if err != nil {
				resp := responses.ErrorResponse{}
				resp.Code = http.StatusBadRequest
				resp.Error.ErrorCode = "invalid_request"
				resp.Error.ErrorMessage = "Could not parse limit."
				resp.SendErrorResponse(c)
				return
			}

			events, err := Svc.GetEventsByTypeWithLimit(eventType, limitInt)
			if err != nil {
				errString := database.CheckDatabaseError(err)
				resp := responses.ErrorResponse{}
				resp.Code = http.StatusInternalServerError
				resp.Error.ErrorCode = "internal_error"
				resp.Error.ErrorMessage = errString
				resp.SendErrorResponse(c)
				return
			}

			resp := responses.SuccessResponse{
				Code: http.StatusOK,
				Data: events,
			}
			resp.SendSuccessReponse(c)
			return
		}

		events, err := Svc.GetEventsByType(eventType)
		if err != nil {
			errString := database.CheckDatabaseError(err)
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusInternalServerError
			resp.Error.ErrorCode = "internal_error"
			resp.Error.ErrorMessage = errString
			resp.SendErrorResponse(c)
			return
		}

		resp := responses.SuccessResponse{
			Code: http.StatusOK,
			Data: events,
		}
		resp.SendSuccessReponse(c)
		return
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Could not parse limit."
			resp.SendErrorResponse(c)
			return
		}

		events, err := Svc.GetEventsWithLimit(limitInt)
		if err != nil {
			errString := database.CheckDatabaseError(err)
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusInternalServerError
			resp.Error.ErrorCode = "internal_error"
			resp.Error.ErrorMessage = errString
			resp.SendErrorResponse(c)
			return
		}

		resp := responses.SuccessResponse{
			Code: http.StatusOK,
			Data: events,
		}
		resp.SendSuccessReponse(c)
		return
	}

	events, err := Svc.GetEvents()
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{
		Code: http.StatusOK,
		Data: events,
	}
	resp.SendSuccessReponse(c)
}

func GetAPILogsRoute(c *gin.Context) {
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

	lines, err := logging.ReadZapLog()
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	if len(lines) == 0 {
		resp := responses.SuccessResponse{}
		resp.Code = http.StatusOK
		resp.Data = gin.H{"logs": "No logs to read so far."}
		resp.SendSuccessReponse(c)
		return
	}

	var noJSONLines []zapLogFormat

	for _, line := range lines {
		var noJSONLine zapLogFormat

		if err := json.Unmarshal([]byte(line), &noJSONLine); err != nil {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusInternalServerError
			resp.Error.ErrorCode = "internal_error"
			resp.Error.ErrorMessage = "Something went wrong, sorry."
			resp.SendErrorResponse(c)
			return
		}

		noJSONLines = append(noJSONLines, noJSONLine)
	}

	sort.Slice(noJSONLines, func(i, j int) bool {
		t1, err := time.Parse("2006-01-02T15:04:05Z", noJSONLines[i].Time)
		if err != nil {
			log.Printf("%s Error parsing date: %s\n", logging.ErrSign, err.Error())
		}

		t2, err := time.Parse("2006-01-02T15:04:05Z", noJSONLines[j].Time)
		if err != nil {
			log.Printf("%s Error parsing date: %s\n", logging.ErrSign, err.Error())
		}

		return t1.After(t2)
	})

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusOK
	resp.Data = gin.H{"logs": noJSONLines}
	resp.SendSuccessReponse(c)
}
