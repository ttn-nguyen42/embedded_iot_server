package privateapi

import (
	"labs/htmx-blog/internal/configs"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func ServiceRegistration() func(app *fiber.App) {
	return func(app *fiber.App) {
		aboutGroup := app.Group("/manage")
		aboutGroup.Get("/home", ManageHomeHandler)
	}
}

func Middlewares(configs *configs.HttpConfigs) []interface{} {
	username := configs.Auth.Username
	token := configs.Auth.Token

	return []interface{}{
		limiter.New(limiter.Config{
			Max:        5,
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
			AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowHeaders: "Origin, Content-Type, Accept-Encoding, Host, Authorization",
		}),
		etag.New(),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		basicauth.New(basicauth.Config{
			Users: map[string]string{
				username: token,
			},
			ContextUsername: "username",
			ContextPassword: "token",
			Unauthorized: func(c *fiber.Ctx) error {
				return c.SendStatus(http.StatusUnauthorized)
			},
		}),
		recover.New(recover.ConfigDefault),
		logger.New(logger.Config{
			DisableColors: true,
			Format:        "ADMIN ${time} [${ip}:${port}] (${latency}) [${locals:username}:${locals:token}] ${method} ${status}  ${path}\n",
			TimeFormat:    time.RFC3339,
		}),
		favicon.New(favicon.Config{
			URL:  "/favicon.ico",
			File: "./public/favicon/favicon.ico",
		}),
	}
}
