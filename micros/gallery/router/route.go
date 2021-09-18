// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package router

import (
	"github.com/GMcD/ts-serverless/micros/gallery/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/middleware/authcookie"
	"github.com/red-gold/telar-core/middleware/authhmac"
	"github.com/red-gold/telar-core/types"

	"github.com/GMcD/cognito-jwt/verify"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {

	// Middleware
	authHMACMiddleware := func(hmacWithCookie bool) func(*fiber.Ctx) error {
		var Next func(c *fiber.Ctx) bool
		if hmacWithCookie {
			Next = func(c *fiber.Ctx) bool {
				if c.Get(types.HeaderHMACAuthenticate) != "" {
					return false
				}
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
					return true
				}
				return false
			}
		}
		return authcookie.New(authcookie.Config{
			Next:         Next,
			JWTSecretKey: []byte(*config.AppConfig.PublicKey),
			Authorizer:   verify.VerifyJWT,
		})
	}

	hmacCookieHandlers := []func(*fiber.Ctx) error{authHMACMiddleware(true), authCookieMiddleware(true)}

	// Routers
	app.Post("/", append(hmacCookieHandlers, handlers.CreateMediaHandle)...)
	app.Post("/list", append(hmacCookieHandlers, handlers.CreateMediaListHandle)...)
	app.Put("/", append(hmacCookieHandlers, handlers.UpdateMediaHandle)...)
	app.Delete("/id/:mediaId", append(hmacCookieHandlers, handlers.DeleteMediaHandle)...)
	app.Delete("/dir/:dir", append(hmacCookieHandlers, handlers.DeleteDirectoryHandle)...)
	app.Get("/", append(hmacCookieHandlers, handlers.QueryAlbumHandle)...)
	app.Get("/id/:mediaId", append(hmacCookieHandlers, handlers.GetMediaHandle)...)
	app.Get("/dir/:dir", append(hmacCookieHandlers, handlers.GetMediaByDirectoryHandle)...)
}
