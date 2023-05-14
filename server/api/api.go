package api

import (
	"errors"
	"log"

	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/updater"
	"github.com/gin-gonic/gin"
)

type API struct {
	Host   string
	Port   int
	Engine *gin.Engine
	Logger *log.Logger
}

func NewAPIInstance(cfg *config.Config, logger *log.Logger) (*API, error) {
	switch updater.BuildMode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		return nil, errors.New("unknown build mode")
	}

	engine := gin.New()

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = false
	engine.HandleMethodNotAllowed = true
	engine.ForwardedByClientIP = true
	engine.UseRawPath = false
	engine.UnescapePathValues = true

	if err := engine.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, err
	}

	return &API{
		cfg.APIHost,
		cfg.APIPort,
		engine,
		logger,
	}, nil
}
