package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/devusSs/crosshairs/api"
	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/api/routes"
	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

const (
	twitchUsersEndpoint = "https://api.twitch.tv/helix/users"
)

var (
	clientID     = ""
	clientSecret = ""
	redirectURL  = ""
	oauthConfig  *oauth2.Config
	state        string
)

type twitchUsersData struct {
	Data []struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Type            string    `json:"type"`
		BroadcasterType string    `json:"broadcaster_type"`
		Description     string    `json:"description"`
		ProfileImageURL string    `json:"profile_image_url"`
		OfflineImageURL string    `json:"offline_image_url"`
		ViewCount       int       `json:"view_count"`
		Email           string    `json:"email"`
		CreatedAt       time.Time `json:"created_at"`
	} `json:"data"`
}

func InitTwitchAuth(cfg *config.Config, api *api.API, hostURL string) error {
	if !strings.Contains(cfg.TwitchRedirectURL, hostURL) {
		hostURL = strings.Replace(hostURL, "127.0.0.1", "localhost", 1)

		if !strings.Contains(cfg.TwitchRedirectURL, hostURL) {
			return fmt.Errorf("twitch redirect url and host url mismatch: %s <-> %s", cfg.TwitchRedirectURL, hostURL)
		}
	}

	if cfg.TwitchClientID == "" || cfg.TwitchClientSecret == "" || cfg.TwitchRedirectURL == "" {
		return errors.New("missing Twitch variables in config")
	}

	clientID = cfg.TwitchClientID
	clientSecret = cfg.TwitchClientSecret
	redirectURL = cfg.TwitchRedirectURL

	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     endpoints.Twitch,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user:read:email"},
	}

	hostURL = strings.Replace(hostURL, "127.0.0.1", "localhost", 1)
	redirectURLHost := strings.Split(cfg.TwitchRedirectURL, hostURL)[1]

	api.Engine.GET("/api/integration/twitch/login", handleLogin)
	api.Engine.GET(redirectURLHost, handleCallback)

	return nil
}

func handleLogin(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("user") == nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "You are currently not logged in."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	_, err := uuid.Parse(fmt.Sprintf("%s", session.Get("user")))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse uuid."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	state = fmt.Sprintf("%d", time.Now().UnixNano())
	url := oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleCallback(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("user") == nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "You are currently not logged in."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	uuidUser, err := uuid.Parse(fmt.Sprintf("%s", session.Get("user")))
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_request"
		resp.Error.ErrorMessage = "Could not parse uuid."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	queryState := c.Query("state")
	if queryState != state {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_state"
		resp.Error.ErrorMessage = "Returned state did not match provided state."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	code := c.Query("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusBadRequest
		resp.Error.ErrorCode = "invalid_token_exchange"
		resp.Error.ErrorMessage = err.Error()
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	client := oauthConfig.Client(context.Background(), token)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, twitchUsersEndpoint, nil)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_sorry"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Client-Id", clientID)

	resp, err := client.Do(req)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_sorry"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respErr := responses.ErrorResponse{}
		respErr.Code = resp.StatusCode
		respErr.Error.ErrorCode = "unwanted_status"
		respErr.Error.ErrorMessage = resp.Status
		respErr.SendErrorResponse(c)
		c.Abort()
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_sorry"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	var data twitchUsersData

	if err := json.Unmarshal(body, &data); err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_sorry"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	userID := data.Data[0].ID
	userLogin := data.Data[0].Login
	userCreatedAt := data.Data[0].CreatedAt

	_, err = routes.Svc.AddUserTwitchDetails(&database.UserAccount{ID: uuidUser, TwitchID: userID, TwitchLogin: userLogin, TwitchCreatedAt: userCreatedAt})
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_sorry"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		resp.SendErrorResponse(c)
		c.Abort()
		return
	}

	respSucc := responses.SuccessResponse{}
	respSucc.Code = http.StatusOK
	respSucc.Data = gin.H{
		"message": "Successfully connected your Twitch account.",
	}
	respSucc.SendSuccessReponse(c)
}
