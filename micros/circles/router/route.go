// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package router

import (
	"log"

	"github.com/GMcD/ts-serverless/micros/circles/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/middleware/authcookie"
	"github.com/red-gold/telar-core/middleware/authhmac"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {

	// Middleware
	authHMACMiddleware := func(hmacWithCookie bool) func(*fiber.Ctx) error {
		var Next func(c *fiber.Ctx) bool
		if hmacWithCookie {
			Next = func(c *fiber.Ctx) bool {
				if c.Get(types.HeaderHMACAuthenticate) != "" {
					log.Println("Have HMAC, returning FALSE...")
					return false
				}
				log.Println("No HMAC, returning TRUE...")
				return true
			}
		}
		return authhmac.New(authhmac.Config{
			Next:          Next,
			PayloadSecret: *config.AppConfig.PayloadSecret,
		})
	}

	authCookieMiddleware := func(hmacWithCookie bool) func(*fiber.Ctx) error {
		var Next func(c *fiber.Ctx) bool
		if hmacWithCookie {
			Next = func(c *fiber.Ctx) bool {
				if c.Get(types.HeaderHMACAuthenticate) != "" {
					log.Println("HMAC present and accounted for...")
					return true
				}
				log.Println("HMAC absent...")
				return false
			}
		}
		return authcookie.New(authcookie.Config{
			Next:         Next,
			JWTSecretKey: []byte(*config.AppConfig.PublicKey),
			Authorizer:   utils.VerifyJWT,
		})
	}

	hmacCookieHandlers := []func(*fiber.Ctx) error{authHMACMiddleware(true), authCookieMiddleware(true)}

	// Routers
	app.Post("/following/:userId", authHMACMiddleware(false), handlers.CreateFollowingHandle)
	app.Post("/", append(hmacCookieHandlers, handlers.CreateCircleHandle)...)
	app.Put("/", append(hmacCookieHandlers, handlers.UpdateCircleHandle)...)
	app.Delete("/:circleId", append(hmacCookieHandlers, handlers.DeleteCircleHandle)...)
	app.Get("/my", append(hmacCookieHandlers, handlers.GetMyCircleHandle)...)
	app.Get("/id/:circleId", append(hmacCookieHandlers, handlers.GetCircleHandle)...)
}
