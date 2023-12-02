package publicapi

import (
	"labs/htmx-blog/internal/configs"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func ServiceRegistration() func(app *fiber.App) {
	return func(app *fiber.App) {
		app.Get("/", UiDashboard)
		app.Get("/partials/room_table", UiPartialTable)
		apiGroup := app.Group("/api")
		apiGroup.Get("/rooms", ApiRoomsHandler)
		apiGroup.Post("/rooms", ApiAddRoomHandler)

	}
}

func Middlewares(configs *configs.HttpConfigs) []interface{} {
	return []interface{}{
		limiter.New(limiter.Config{
			Max:        10,
			Expiration: 1 * time.Second,
			LimitReached: func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusTooManyRequests)
			},
		}),
		cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				switch origin {
				case "http://*":
				case "https://*":
				default:
					return false
				}
				return true
			},
			AllowMethods: "GET,HEAD",
			AllowHeaders: "Origin, Content-Type, Accept-Encoding, Host",
		}),
		etag.New(),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		recover.New(recover.ConfigDefault),
		logger.New(logger.Config{
			DisableColors: true,
			Format:        "PUBLIC ${time} [${ip}:${port}] ${latency} ${method} ${status} ${path}\n",
			TimeFormat:    time.RFC3339,
		}),
		favicon.New(favicon.Config{
			URL:  "/favicon.ico",
			File: "./public/favicon/favicon.ico",
		}),
		cache.New(cache.Config{
			Expiration:   time.Minute * 1,
			CacheControl: false,
			CacheHeader:  "X-Cache",
			Methods: []string{
				fiber.MethodGet,
				fiber.MethodHead,
			},
		}),
		helmet.New(helmet.ConfigDefault),
	}
}
