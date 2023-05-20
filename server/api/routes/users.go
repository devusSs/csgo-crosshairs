package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devusSs/crosshairs/api/models"
	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/database"
	"github.com/devusSs/crosshairs/updater"
	"github.com/devusSs/crosshairs/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	lenPasswordNeeded = 8
)

var (
	SRVAddr string
)

func RegisterUserRoute(c *gin.Context) {
	var registerUser models.RegisterUser

	if err := c.BindJSON(&registerUser); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid JSON body provided."
		resp.SendErrorResponse(c)
		return
	}

	if !utils.IsEmailValid(registerUser.EMail) {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid e-mail address provided."
		resp.SendErrorResponse(c)
		return
	}

	if len(registerUser.Password) < lenPasswordNeeded {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = fmt.Sprintf("Password needs to be at least %d characters long.", lenPasswordNeeded)
		resp.SendErrorResponse(c)
		return
	}

	hashedPassword, err := utils.HashPassword(registerUser.Password)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not hash password."
		resp.SendErrorResponse(c)
		return
	}

	verificationCode := utils.RandomString(25)

	newUser := database.UserAccount{
		EMail:            registerUser.EMail,
		Password:         hashedPassword,
		Role:             "user",
		VerificationCode: verificationCode,
		VerifiedMail:     false,
		RegisterIP:       c.Request.Header.Get("X-Forwarded-For"),
	}

	var emailData *utils.EmailData

	if updater.BuildMode == "dev" {
		emailData = &utils.EmailData{
			URL:     fmt.Sprintf("http://%s/api/users/verifyMail/%s", SRVAddr, utils.Encode(verificationCode)),
			Subject: "E-Mail verification",
		}
	} else {
		emailData = &utils.EmailData{
			URL:     fmt.Sprintf("http://%s/api/users/verifyMail/%s", CFG.BackendDomain, utils.Encode(verificationCode)),
			Subject: "E-Mail verification",
		}
	}

	if err := utils.SendEmail(&newUser, emailData, CFG, utils.MailVerfication); err != nil {
		log.Println(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not send confirmation email."
		resp.SendErrorResponse(c)
		return
	}

	_, err = Svc.AddUser(&newUser)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = err.Error()
		resp.SendErrorResponse(c)
		return
	}

	var event database.CreateEvent

	event.Type = database.UserRegistered
	event.Data.URL = c.Request.RequestURI
	event.Data.Method = c.Request.Method
	event.Data.IssuerIP = c.Request.Header.Get("X-Forwarded-For")
	eventData, err := json.Marshal(registerUser)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}
	event.Data.Data = eventData
	event.Timestamp = time.Now()

	var eventDB database.Event
	encodedEventType, err := json.Marshal(event.Type)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}
	encodedEventData, err := json.Marshal(event.Data)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	eventDB.Type = string(encodedEventType)
	eventDB.Data = string(encodedEventData)
	eventDB.Timestamp = event.Timestamp

	_, err = Svc.AddEvent(&eventDB)
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
	resp.Data = responses.GeneralUserResponse{
		Message: "Please check your e-mails.",
	}
	resp.SendSuccessReponse(c)
}

func VerifyUserEMailRoute(c *gin.Context) {
	codeRaw := c.Params.ByName("code")
	code, err := utils.Decode(codeRaw)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Error decoding verification token."
		resp.SendErrorResponse(c)
		return
	}

	user, err := Svc.GetUserByVerificationCode(&database.UserAccount{VerificationCode: code})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "User does not exist."
		resp.SendErrorResponse(c)
		return
	}

	if code != user.VerificationCode {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid verification code."
		resp.SendErrorResponse(c)
		return
	}

	user.VerifiedMail = true

	if _, err := Svc.UpdateUserVerification(user); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Error updating user verification."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusOK
	resp.Data = responses.GeneralUserResponse{
		Message: "Successfully verified your e-mail address. Please login now.",
	}
	resp.SendSuccessReponse(c)
}

func LoginUserRoute(c *gin.Context) {
	var loginUser models.LoginUser

	if err := c.BindJSON(&loginUser); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid JSON body provided."
		resp.SendErrorResponse(c)
		return
	}

	if !utils.IsEmailValid(loginUser.EMail) {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid e-mail address provided."
		resp.SendErrorResponse(c)
		return
	}

	if len(loginUser.Password) < lenPasswordNeeded {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Password provided is too short."
		resp.SendErrorResponse(c)
		return
	}

	user, err := Svc.GetUserByEmail(&database.UserAccount{EMail: loginUser.EMail})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "E-Mail address does not exist."
		resp.SendErrorResponse(c)
		return
	}

	if !user.VerifiedMail {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Please confirm your e-mail address first."
		resp.SendErrorResponse(c)
		return
	}

	if err := utils.VerifyPassword(user.Password, loginUser.Password); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Passwords do not match."
		resp.SendErrorResponse(c)
		return
	}

	user.LastLogin = time.Now()
	user.LoginIP = c.Request.Header.Get("X-Forwarded-For")

	userDB, err := Svc.UpdateUserLogin(user)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	session := sessions.Default(c)
	session.Set("user", user.ID.String())
	if err := session.Save(); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not set session cookie."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusOK
	resp.Data = responses.LoginUserResponse{
		Message: "Successfully logged in.",
		Role:    userDB.Role,
	}
	resp.SendSuccessReponse(c)
}

func GetUserRoute(c *gin.Context) {
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

	var userReturn models.ReturnUser
	userReturn.CreatedAt = user.CreatedAt
	userReturn.EMail = user.EMail
	userReturn.Role = user.Role

	resp := responses.SuccessResponse{
		Code: http.StatusOK,
		Data: userReturn,
	}
	resp.SendSuccessReponse(c)
}

func LogoutUserRoute(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("user", "")
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not remove session."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusNoContent
	resp.SendSuccessReponse(c)
}

func ResetPasswordRoute(c *gin.Context) {
	var resetPass models.ResetPassword

	if err := c.BindJSON(&resetPass); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid JSON body provided."
		resp.SendErrorResponse(c)
		return
	}

	if !utils.IsEmailValid(resetPass.EMail) {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid e-mail address provided."
		resp.SendErrorResponse(c)
		return
	}

	user, err := Svc.GetUserByEmail(&database.UserAccount{EMail: resetPass.EMail})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusNotFound
		resp.Error.ErrorCode = "not_found"
		resp.Error.ErrorMessage = "User could not be found."
		resp.SendErrorResponse(c)
		return
	}

	verificationCode := utils.RandomString(25)

	_, err = Svc.AddResetPasswordCode(&database.UserAccount{EMail: resetPass.EMail, PasswordResetCode: verificationCode})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	var emailData *utils.EmailData

	if updater.BuildMode == "dev" {
		emailData = &utils.EmailData{
			URL:     fmt.Sprintf("http://%s/api/users/resetPass/%s?code=%s", SRVAddr, resetPass.EMail, utils.Encode(verificationCode)),
			Subject: "Reset your password",
		}
	} else {
		emailData = &utils.EmailData{
			URL:     fmt.Sprintf("http://%s/api/users/resetPass/%s?code=%s", CFG.BackendDomain, resetPass.EMail, utils.Encode(verificationCode)),
			Subject: "Reset your password",
		}
	}

	if err := utils.SendEmail(user, emailData, CFG, utils.MailVerificationPassword); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not send confirmation email."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{
		Code: http.StatusOK,
		Data: responses.GeneralUserResponse{
			Message: "Please check your e-mails.",
		},
	}
	resp.SendSuccessReponse(c)
}

func VerifyUserPasswordCodeRoute(c *gin.Context) {
	email := c.Param("email")
	code := c.Query("code")

	if email == "" || code == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Missing e-mail address or code."
		resp.SendErrorResponse(c)
		return
	}

	code, err := utils.Decode(code)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	if !utils.IsEmailValid(email) {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid e-mail address provided."
		resp.SendErrorResponse(c)
		return
	}

	_, err = Svc.GetUserByResetpasswordCode(&database.UserAccount{EMail: email, PasswordResetCode: code})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "E-Mail or code mismatch."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{
		Code: http.StatusOK,
		Data: responses.GeneralUserResponse{
			Message: "Successfully confirmed your e-mail. You may now reset your password.",
		},
	}
	resp.SendSuccessReponse(c)
}

func ResetPasswordRouteFinal(c *gin.Context) {
	email := c.Param("email")
	code := c.Query("code")

	if email == "" || code == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Missing e-mail address or code."
		resp.SendErrorResponse(c)
		return
	}

	code, err := utils.Decode(code)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	if !utils.IsEmailValid(email) {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid e-mail address provided."
		resp.SendErrorResponse(c)
		return
	}

	var resetPasswordFinal models.ResetPasswordFinal

	if err := c.BindJSON(&resetPasswordFinal); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid JSON body provided."
		resp.SendErrorResponse(c)
		return
	}

	if len(resetPasswordFinal.Password) < lenPasswordNeeded {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = fmt.Sprintf("Password must be at least %d characters long.", lenPasswordNeeded)
		resp.SendErrorResponse(c)
		return
	}

	user, err := Svc.GetUserByResetpasswordCode(&database.UserAccount{EMail: email, PasswordResetCode: code})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "E-Mail or code mismatch."
		resp.SendErrorResponse(c)
		return
	}

	hashedPassword, err := utils.HashPassword(resetPasswordFinal.Password)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not hash password."
		resp.SendErrorResponse(c)
		return
	}

	user.EMail = email
	user.PasswordResetCode = code
	user.Password = hashedPassword

	_, err = Svc.UpdateUserPassword(user)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not update password."
		resp.SendErrorResponse(c)
		return
	}

	var event database.CreateEvent

	event.Type = database.UserChangedPassword
	event.Data.URL = c.Request.RequestURI
	event.Data.Method = c.Request.Method
	event.Data.IssuerIP = c.Request.Header.Get("X-Forwarded-For")
	eventData, err := json.Marshal(resetPasswordFinal)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}
	event.Data.Data = eventData
	event.Timestamp = time.Now()

	var eventDB database.Event
	encodedEventType, err := json.Marshal(event.Type)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}
	encodedEventData, err := json.Marshal(event.Data)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	eventDB.Type = string(encodedEventType)
	eventDB.Data = string(encodedEventData)
	eventDB.Timestamp = event.Timestamp

	_, err = Svc.AddEvent(&eventDB)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{
		Code: http.StatusOK,
		Data: responses.GeneralUserResponse{
			Message: "Successfully reset your password.",
		},
	}
	resp.SendSuccessReponse(c)
}
