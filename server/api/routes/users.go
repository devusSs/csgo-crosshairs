package routes

import (
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

	var adminTokenProvided bool
	if registerUser.AdminToken != "" {
		if registerUser.AdminToken != CFG.AdminKey {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "Invalid admin token provided."
			resp.SendErrorResponse(c)
			return
		}
		adminTokenProvided = true
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
		RegisterIP:       c.RemoteIP(),
	}

	if adminTokenProvided {
		newUser.Role = "admin"
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

	if err := utils.SendEmail(&newUser, emailData, CFG); err != nil {
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

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusOK
	resp.Data = "Please confirm your e-mail address."
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
	resp.Data = "Successfully verified your e-mail address. You may now log in!"
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
	user.LoginIP = c.RemoteIP()

	_, err = Svc.UpdateUserLogin(user)
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
	resp.Data = "Successfully logged in."
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
