package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/api/routes"
	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database"
	"github.com/devusSs/crosshairs/logging"
	"github.com/devusSs/crosshairs/updater"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	cors "github.com/rs/cors/wrapper/gin"
)

type API struct {
	Host            string
	Port            int
	Engine          *gin.Engine
	RequestsLogFile *os.File
}

func NewAPIInstance(cfg *config.Config, requestsLogFile *os.File) (*API, error) {
	switch updater.BuildMode {
	case "dev":
		gin.SetMode(gin.DebugMode)
		gin.DefaultWriter = os.Stdout
	case "release":
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = requestsLogFile
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
		requestsLogFile,
	}, nil
}

func (api *API) SetupSessions(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword,
		cfg.PostgresDB, cfg.PostgresPort)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	store, err := postgres.NewStore(db, []byte(cfg.SecretSessionsKey))
	if err != nil {
		return err
	}

	if updater.BuildMode == "dev" {
		store.Options(sessions.Options{
			Path:     "/",
			HttpOnly: true,
			MaxAge:   30 * 24 * 60 * 60 * 1000, // 30 days until expiry, does not really matter for dev
			Secure:   false,
		})
	} else {
		store.Options(sessions.Options{
			Path:     "/",
			Domain:   cfg.Domain,
			HttpOnly: true,
			MaxAge:   30 * 24 * 60 * 60 * 1000, // 30 days until expiry
			Secure:   true,
		})
	}

	api.Engine.Use(sessions.Sessions("sessions", store))

	return nil
}

func (api *API) SetupRedisRateLimiting(cfg *config.Config) {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword,
		}),
		Rate:  time.Second,
		Limit: 5,
	})

	rateMW := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: rateLimitError,
		KeyFunc:      rateLimitGetIP,
	})

	api.Engine.Use(rateMW)
}

func (api *API) SetupCors(cfg *config.Config) {
	var c gin.HandlerFunc

	if updater.BuildMode == "dev" {
		c = cors.New(cors.Options{
			AllowedOrigins:      []string{"http://localhost:5173"},
			AllowedMethods:      []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodHead},
			AllowedHeaders:      []string{"Content-Type", "Content-Length"},
			AllowPrivateNetwork: true,
			AllowCredentials:    true,
			MaxAge:              0,
			Debug:               true,
		})
	} else {
		c = cors.New(cors.Options{
			AllowedOrigins:   []string{cfg.Domain},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodHead},
			AllowedHeaders:   []string{"Content-Type", "Content-Length"},
			AllowCredentials: true,
			MaxAge:           43200, // 12 hours
		})
	}

	api.Engine.Use(c)
}

func (api *API) SetupRoutes(db database.Service, cfg *config.Config) {
	api.Engine.Use(gin.Recovery())
	api.Engine.Use(gin.Logger())

	routes.CFG = cfg
	routes.Svc = db

	api.Engine.NoRoute(routes.NotFoundRoute)
	api.Engine.NoMethod(routes.MethodNotAllowedRoute)

	base := api.Engine.Group("/api")
	{
		base.GET("/", routes.HomeRoute)

		users := base.Group("/users")
		{
			users.POST("/register", routes.RegisterUserRoute)
			users.GET("/verifyMail/:code", routes.VerifyUserEMailRoute)
			users.POST("/login", routes.LoginUserRoute)
			users.GET("/me", routes.GetUserRoute)
			users.GET("/logout", routes.LogoutUserRoute)
			users.POST("/resetPass", routes.ResetPasswordRoute)
			users.GET("/resetPass/:email", routes.VerifyUserPasswordCodeRoute)
			users.PATCH("/resetPass/:email", routes.ResetPasswordRouteFinal)
		}

		crosshairs := base.Group("/crosshairs")
		{
			crosshairs.POST("/add", routes.AddCrosshairRoute)
			crosshairs.GET("", routes.GetAllCrosshairsFromUserRoute)
			crosshairs.DELETE("", routes.DeleteAllCrosshairsFromUserRoute)
			crosshairs.DELETE("/single/:code", routes.DeleteOneCrosshairFromUserRoute)
		}

		admins := base.Group("/admins")
		{
			admins.GET("/users", routes.GetAllUsersRoute)
			admins.GET("/crosshairs", routes.GetAllCrosshairsRoute)

			events := admins.Group("/events")
			{
				events.GET("", routes.GetAllEventsRoute)
				events.GET("/:type", routes.GetEventsByTypeRoute)
			}
		}
	}
}

func (api *API) StartAPI() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", api.Host, api.Port),
		Handler: api.Engine,
	}

	routes.SRVAddr = srv.Addr

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[%s] Error starting API: %s\n", logging.ErrSign, err.Error())
		}
	}()

	log.Printf("[%s] API started on 'http://%s'\n", logging.InfSign, srv.Addr)
	log.Printf("[%s] Press CTRL+C to exit any time\n", logging.InfSign)
	log.Printf("[%s] Please DO NOT force exit the app\n", logging.InfSign)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	defer ctx.Done()

	fmt.Println("")

	return api.RequestsLogFile.Close()
}

func rateLimitGetIP(c *gin.Context) string {
	return c.ClientIP()
}

func rateLimitError(c *gin.Context, info ratelimit.Info) {
	var resp responses.ErrorResponse
	resp.Code = http.StatusTooManyRequests
	resp.Error.ErrorCode = "flooding"
	resp.Error.ErrorMessage = "Too many requests. Try again in " + time.Until(info.ResetTime).String()
	resp.SendErrorResponse(c)
}
