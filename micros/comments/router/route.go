// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package router

import (
	"github.com/GMcD/ts-serverless/micros/comments/handlers"
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
			Authorizer:   utils.VerifyJWT,
		})
	}

	hmacCookieHandlers := []func(*fiber.Ctx) error{authHMACMiddleware(true), authCookieMiddleware(true)}

	// Routers
	app.Post("/", append(hmacCookieHandlers, handlers.CreateCommentHandle)...)
	app.Put("/", append(hmacCookieHandlers, handlers.UpdateCommentHandle)...)
	app.Put("/profile", append(hmacCookieHandlers, handlers.UpdateCommentProfileHandle)...)
	app.Delete("/id/:commentId/post/:postId", append(hmacCookieHandlers, handlers.DeleteCommentHandle)...)
	app.Delete("/post/:postId", append(hmacCookieHandlers, handlers.DeleteCommentByPostIdHandle)...)
	app.Get("/", append(hmacCookieHandlers, handlers.GetCommentsByPostIdHandle)...)
	app.Get("/:commentId", append(hmacCookieHandlers, handlers.GetCommentHandle)...)
}
