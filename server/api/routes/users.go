package routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devusSs/crosshairs/api/models"
	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/database"
	"github.com/devusSs/crosshairs/stats"
	"github.com/devusSs/crosshairs/storage"
	"github.com/devusSs/crosshairs/updater"
	"github.com/devusSs/crosshairs/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	lenPasswordNeeded = 8
	tmpDir            = "./tmp"
)

var (
	SRVAddr           string
	UsingReverseProxy bool = false
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

	resendMail := c.Query("action")

	if resendMail == "resend" {
		user, err := Svc.GetUserByEmail(&database.UserAccount{EMail: registerUser.EMail})
		if err != nil {
			errString := database.CheckDatabaseError(err)
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = errString
			resp.SendErrorResponse(c)
			return
		}

		if user.VerifiedMail {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = "E-Mail address has already been verified."
			resp.SendErrorResponse(c)
			return
		}

		if time.Since(user.CreatedAt) < 5*time.Minute {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = fmt.Sprintf("Please wait %.2f second(s) before retrying.", time.Until(user.CreatedAt).Seconds())
			resp.SendErrorResponse(c)
			return
		}

		if !user.RequestNewVerifyMailTime.IsZero() && time.Since(user.RequestNewVerifyMailTime) < 10*time.Minute {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusBadRequest
			resp.Error.ErrorCode = "invalid_request"
			resp.Error.ErrorMessage = fmt.Sprintf("Please wait %.2f second(s) before retrying.", time.Until(user.RequestNewVerifyMailTime).Seconds())
			resp.SendErrorResponse(c)
			return
		}

		var emailData *utils.EmailData

		if updater.BuildMode == "dev" {
			emailData = &utils.EmailData{
				// URL = backend.
				URL:     fmt.Sprintf("http://%s/api/users/verifyMail?code=%s", SRVAddr, utils.Encode(user.VerificationCode)),
				Subject: "dropawp.com - E-Mail verification",
			}
		} else {
			emailData = &utils.EmailData{
				// URL = frontend.
				URL:     fmt.Sprintf("%s/users/register?code=%s", CFG.Domain, utils.Encode(user.VerificationCode)),
				Subject: "dropawp.com - E-Mail verification",
			}
		}

		_, err = Svc.UpdateVerifyMailResendTime(&database.UserAccount{EMail: user.EMail, RequestNewVerifyMailTime: time.Now()})
		if err != nil {
			errString := database.CheckDatabaseError(err)
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusInternalServerError
			resp.Error.ErrorCode = "internal_error"
			resp.Error.ErrorMessage = errString
			resp.SendErrorResponse(c)
			return
		}

		if err := utils.SendVerificationMail(user, emailData); err != nil {
			resp := responses.ErrorResponse{}
			resp.Code = http.StatusInternalServerError
			resp.Error.ErrorCode = "internal_error"
			resp.Error.ErrorMessage = "Could not send confirmation email."
			resp.SendErrorResponse(c)
			return
		}

		resp := responses.SuccessResponse{}
		resp.Code = http.StatusOK
		resp.Data = responses.GeneralUserResponse{
			Message: "Please check your e-mails.",
		}
		resp.SendSuccessReponse(c)
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
	}

	if UsingReverseProxy {
		newUser.RegisterIP = c.Request.Header.Get("X-Forwarded-For")
	} else {
		newUser.RegisterIP = c.RemoteIP()
	}

	_, err = Svc.AddUser(&newUser)
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	var emailData *utils.EmailData

	if updater.BuildMode == "dev" {
		emailData = &utils.EmailData{
			// URL = backend.
			URL:     fmt.Sprintf("http://%s/api/users/verifyMail?code=%s", SRVAddr, utils.Encode(verificationCode)),
			Subject: "dropawp.com - E-Mail verification",
		}
	} else {
		emailData = &utils.EmailData{
			// URL = frontend.
			URL:     fmt.Sprintf("%s/users/register?code=%s", CFG.Domain, utils.Encode(verificationCode)),
			Subject: "dropawp.com - E-Mail verification",
		}
	}

	if err := utils.SendVerificationMail(&newUser, emailData); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not send confirmation email."
		resp.SendErrorResponse(c)
		return
	}

	var event database.Event

	event.Type = database.UserRegistered
	event.Data.URL = c.Request.RequestURI
	event.Data.Method = c.Request.Method
	event.Data.IssuerIP = c.Request.Header.Get("X-Forwarded-For")
	event.Timestamp = time.Now()

	_, err = Svc.AddEvent(&event)
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
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
	codeRaw := c.Query("code")

	if codeRaw == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Missing code."
		resp.SendErrorResponse(c)
		return
	}

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
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	if user.VerifiedMail {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "E-Mail has already been verified."
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
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	stats.UsersRegisteredLast24Hours++

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
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = errString
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

	if UsingReverseProxy {
		user.LoginIP = c.Request.Header.Get("X-Forwarded-For")
	} else {
		user.LoginIP = c.RemoteIP()
	}

	userDB, err := Svc.UpdateUserLogin(user)
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
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

	stats.UsersLoggedInLast24Hours++

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

	profilePictureLink := user.AvatarURL

	if profilePictureLink == "" {
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

		profilePictureLink = defaultAvatar
	}

	var userReturn models.ReturnUser
	userReturn.ProfilePictureLink = profilePictureLink

	// Relevant for Docker only.
	userReturn.ProfilePictureLink = strings.Replace(userReturn.ProfilePictureLink, "http://minio:", fmt.Sprintf("http://%s:", "localhost"), 1)

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

	if session.Get("user") == nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "You are currently not logged in."
		resp.SendErrorResponse(c)
		return
	}

	session.Set("user", "")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
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
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusNotFound
		resp.Error.ErrorCode = "not_found"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	if !user.PasswordResetCodeTime.IsZero() && time.Since(user.PasswordResetCodeTime) < 10*time.Minute {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = fmt.Sprintf("Already requested a password reset. Please wait %.2f second(s).", time.Until(user.PasswordResetCodeTime).Seconds())
		resp.SendErrorResponse(c)
		return
	}

	verificationCode := utils.RandomString(25)
	resetCodeTime := time.Now()

	_, err = Svc.AddResetPasswordCodeAndTime(&database.UserAccount{EMail: resetPass.EMail, PasswordResetCode: verificationCode, PasswordResetCodeTime: resetCodeTime})
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	var emailData *utils.EmailData

	if updater.BuildMode == "dev" {
		emailData = &utils.EmailData{
			URL:     fmt.Sprintf("http://%s/api/users/resetPass?email=%s&code=%s", SRVAddr, resetPass.EMail, utils.Encode(verificationCode)),
			Subject: "dropawp.com - Reset your password",
		}
	} else {
		emailData = &utils.EmailData{
			URL:     fmt.Sprintf("http://%s/users/reset-password?email=%s&code=%s", CFG.Domain, resetPass.EMail, utils.Encode(verificationCode)),
			Subject: "dropawp.com - Reset your password",
		}
	}

	if err := utils.SendVerificationMail(user, emailData); err != nil {
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
	email := c.Query("email")
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
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = errString
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
	email := c.Query("email")
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
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = errString
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
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	var event database.Event

	event.Type = database.UserChangedPassword
	event.Data.URL = c.Request.RequestURI
	event.Data.Method = c.Request.Method
	event.Data.IssuerIP = c.Request.Header.Get("X-Forwarded-For")
	event.Timestamp = time.Now()

	_, err = Svc.AddEvent(&event)
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
		Data: responses.GeneralUserResponse{
			Message: "Successfully reset your password.",
		},
	}
	resp.SendSuccessReponse(c)
}

func ResetPasswordWhenLoggedInRoute(c *gin.Context) {
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

	var requestBody models.RequestPWResetLoggedIn

	if err := c.BindJSON(&requestBody); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Invalid request body."
		resp.SendErrorResponse(c)
		return
	}

	if len(requestBody.CurrentPassword) < lenPasswordNeeded || len(requestBody.NewPassword) < lenPasswordNeeded {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = fmt.Sprintf("A password needs to be at least %d characters long.", lenPasswordNeeded)
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

	if err := utils.VerifyPassword(user.Password, requestBody.CurrentPassword); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Old passwords do not match."
		resp.SendErrorResponse(c)
		return
	}

	hashedNewPassword, err := utils.HashPassword(requestBody.NewPassword)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Could not hash password."
		resp.SendErrorResponse(c)
		return
	}

	_, err = Svc.UpdateUserPasswordRaw(&database.UserAccount{EMail: user.EMail, Password: hashedNewPassword})
	if err != nil {
		errString := database.CheckDatabaseError(err)
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = errString
		resp.SendErrorResponse(c)
		return
	}

	var event database.Event

	event.Type = database.UserChangedPassword
	event.Data.URL = c.Request.RequestURI
	event.Data.Method = c.Request.Method
	event.Data.IssuerIP = c.Request.Header.Get("X-Forwarded-For")
	event.Timestamp = time.Now()

	_, err = Svc.AddEvent(&event)
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
		Data: responses.GeneralUserResponse{
			Message: "Successfully updated your password.",
		},
	}
	resp.SendSuccessReponse(c)
}

func UploadUserAvatarRoute(c *gin.Context) {
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

	file, err := c.FormFile("avatar")
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Missing form file in request."
		resp.SendErrorResponse(c)
		return
	}

	if err := storage.CheckFileValid(file); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = err.Error()
		resp.SendErrorResponse(c)
		return
	}

	fileName := fmt.Sprintf("%s.png", uuidUser.String())
	filePath := filepath.Join(tmpDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	if err := StorageSvc.UpdateUserProfilePicture(fileName, filePath); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	link, err := StorageSvc.GetUserProfilePictureLink(uuidUser.String())
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	_, err = Svc.UpdateUserAvatarURL(&database.UserAccount{ID: uuidUser, AvatarURL: link})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	if err := os.Remove(filePath); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	var event database.Event

	event.Type = database.UserUploadedAvatar
	event.Data.URL = c.Request.RequestURI
	event.Data.Method = c.Request.Method
	event.Data.IssuerIP = c.Request.Header.Get("X-Forwarded-For")
	event.Timestamp = time.Now()

	_, err = Svc.AddEvent(&event)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	var userReturn models.ReturnUserAvatar
	userReturn.ID = uuidUser
	userReturn.AvatarURL = link

	// Relevant for Docker only.
	userReturn.AvatarURL = strings.Replace(userReturn.AvatarURL, "http://minio:", fmt.Sprintf("http://%s:", "localhost"), 1)

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusOK
	resp.Data = userReturn
	resp.SendSuccessReponse(c)
}

func DeleteUserAvatarRoute(c *gin.Context) {
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

	if err := StorageSvc.DeleteUserProfilePicture(uuidUser.String()); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "User does not own an avatar"
		resp.SendErrorResponse(c)
		return
	}

	_, err = Svc.UpdateUserAvatarURL(&database.UserAccount{ID: uuidUser, AvatarURL: ""})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		return
	}

	resp := responses.SuccessResponse{}
	resp.Code = http.StatusNoContent
	resp.SendSuccessReponse(c)
}
