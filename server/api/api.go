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
	"strings"
	"syscall"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/devusSs/crosshairs/api/middleware"
	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/api/routes"
	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database"
	"github.com/devusSs/crosshairs/logging"
	"github.com/devusSs/crosshairs/stats"
	"github.com/devusSs/crosshairs/storage"
	"github.com/devusSs/crosshairs/updater"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	cors "github.com/rs/cors/wrapper/gin"
)

var (
	UsingReverseProxy bool = false
)

type API struct {
	Host   string
	Port   int
	Engine *gin.Engine
}

func NewAPIInstance(cfg *config.Config) (*API, error) {
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

	engine.MaxMultipartMemory = 2 << 20 // 2 MiB maximum file size

	if err := engine.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, err
	}

	return &API{
		cfg.APIHost,
		cfg.APIPort,
		engine,
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
		})
	} else {
		store.Options(sessions.Options{
			Path:     "/",
			Domain:   strings.Replace(cfg.Domain, "https://", "", 1),
			HttpOnly: true,
			MaxAge:   30 * 24 * 60 * 60 * 1000, // 30 days until expiry
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	api.Engine.Use(sessions.Sessions("sessions", store))

	return nil
}

func (api *API) SetupRedisRateLimiting(cfg *config.Config) error {
	rClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
	})

	redisServerVersion, err := rClient.Do(context.Background(), "info", "server").Result()
	if err != nil {
		return err
	}

	redisServerVersionSplit := strings.Split(fmt.Sprintf("%v", redisServerVersion), "\n")

	redisServerVersionFinal := ""

	for _, line := range redisServerVersionSplit {
		if strings.Contains(line, "redis_version:") {
			redisServerVersionFinal = strings.TrimSpace(line)
			continue
		}

		if strings.Contains(line, "os:") {
			redisServerVersionFinal = redisServerVersionFinal + " " + strings.TrimSpace(line)
			continue
		}
	}

	redisServerVersionFinal = redisServerVersionFinal + " go-redis client version:" + redis.Version()

	stats.RedisVersion = redisServerVersionFinal

	if err := rClient.Close(); err != nil {
		return err
	}

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

	return nil
}

func (api *API) SetupCors(cfg *config.Config) {
	var c gin.HandlerFunc

	if updater.BuildMode == "dev" {
		c = cors.New(cors.Options{
			AllowedOrigins:      []string{"http://localhost:5173"}, // Used for vite projects.
			AllowedMethods:      []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
			AllowedHeaders:      []string{"Content-Type", "Content-Length"},
			AllowPrivateNetwork: true,
			AllowCredentials:    true,
			MaxAge:              0,
			Debug:               true,
		})
	} else {
		c = cors.New(cors.Options{
			AllowedOrigins:   []string{cfg.Domain},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
			AllowedHeaders:   []string{"Content-Type", "Content-Length"},
			AllowCredentials: true,
			MaxAge:           43200, // 12 hours caching for preflight requests
		})
	}

	api.Engine.Use(c)
}

func (api *API) SetupRoutes(db database.Service, strSvc *storage.Service, cfg *config.Config, logsDir string, debug bool) error {
	logger := logging.InitZapAPILogger(logsDir, debug)

	api.Engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	api.Engine.Use(ginzap.RecoveryWithZap(logger, true))

	routes.CFG = cfg
	routes.Svc = db
	routes.StorageSvc = strSvc

	if err := middleware.SetupPrivateIPBlock(); err != nil {
		return err
	}

	middleware.AllowedDomain = cfg.AllowedDomain

	api.Engine.NoRoute(routes.NotFoundRoute)
	api.Engine.NoMethod(routes.MethodNotAllowedRoute)

	base := api.Engine.Group("/api")
	{
		base.Use(middleware.CountRequestsMiddleware)

		base.GET("/", routes.HomeRoute)

		users := base.Group("/users")
		{
			users.POST("/register", routes.RegisterUserRoute)
			users.GET("/verifyMail", routes.VerifyUserEMailRoute)
			users.POST("/login", routes.LoginUserRoute)
			users.GET("/me", routes.GetUserRoute)
			users.GET("/logout", routes.LogoutUserRoute)
			users.POST("/resetPass", routes.ResetPasswordRoute)
			users.GET("/resetPass", routes.VerifyUserPasswordCodeRoute)
			users.PATCH("/resetPass", routes.ResetPasswordRouteFinal)
			users.PATCH("/newPass", routes.ResetPasswordWhenLoggedInRoute)

			users.POST("/avatar", routes.UploadUserAvatarRoute)
			users.DELETE("/avatar", routes.DeleteUserAvatarRoute)
		}

		crosshairs := base.Group("/crosshairs")
		{
			crosshairs.POST("/add", routes.AddCrosshairRoute)
			crosshairs.GET("", routes.GetAllCrosshairsFromUserRoute)
			crosshairs.DELETE("", routes.DeleteOneOrMultipleCrosshairs)
		}

		admins := base.Group("/admins")
		{
			admins.GET("/users", routes.GetAllUsersRoute)
			admins.GET("/crosshairs", routes.GetAllCrosshairsRoute)
			admins.GET("/logs", routes.GetAPILogsRoute)

			events := admins.Group("/events")
			{
				events.GET("", routes.GetAllEventsOrByTypeRoute)
			}

			stats := admins.Group("/stats")
			{
				stats.GET("/total", routes.GetTotalStatsRoute)
				stats.GET("/daily", routes.Get24HourStatsRoute)

				// Route is only accessable for engineers / users with ACTUAL database access.
				system := stats.Group("/system")
				{
					system.Use(middleware.CheckAllowedHostMiddleware)
					system.Use(middleware.VerifyEngineerMiddleware)
					system.GET("", routes.GetSystemStatsRoute)
				}
			}
		}
	}

	return nil
}

func (api *API) StartAPI() error {
	time.AfterFunc(stats.CalculateTimeUntilMidnight(), stats.Reset24Statistics)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", api.Host, api.Port),
		Handler:      api.Engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if srv.Addr == fmt.Sprintf(":%d", api.Port) {
		routes.SRVAddr = fmt.Sprintf("localhost:%d", api.Port)
	} else {
		routes.SRVAddr = srv.Addr
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[%s] Error starting API: %s\n", logging.ErrSign, err.Error())
		}
	}()

	var addr string

	addr = srv.Addr

	if srv.Addr == fmt.Sprintf(":%d", api.Port) {
		addr = fmt.Sprintf("localhost:%d", api.Port)
	}

	log.Printf("%s API started on 'http://%s'\n", logging.SucSign, addr)
	log.Printf("%s Press CTRL+C to exit any time\n", logging.InfSign)
	log.Printf("%s Please DO NOT force exit the app\n", logging.InfSign)

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

	return nil
}

func rateLimitGetIP(c *gin.Context) string {
	if UsingReverseProxy {
		return c.Request.Header.Get("X-Forwarded-For")
	}

	return c.RemoteIP()
}

func rateLimitError(c *gin.Context, info ratelimit.Info) {
	var resp responses.ErrorResponse
	resp.Code = http.StatusTooManyRequests
	resp.Error.ErrorCode = "flooding"
	resp.Error.ErrorMessage = fmt.Sprintf("Too many requests. Try again in %.2f second(s).", time.Until(info.ResetTime).Seconds())
	resp.SendErrorResponse(c)
}
