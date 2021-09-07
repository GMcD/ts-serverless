// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/middleware/authcookie"
	"github.com/red-gold/telar-core/middleware/authhmac"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/ts-serverless/micros/posts/handlers"

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
	app.Post("/", append(hmacCookieHandlers, handlers.CreatePostHandle)...)
	app.Post("/index", authHMACMiddleware(false), handlers.InitPostIndexHandle)
	app.Put("/", append(hmacCookieHandlers, handlers.UpdatePostHandle)...)
	app.Put("/profile", append(hmacCookieHandlers, handlers.UpdatePostProfileHandle)...)
	app.Put("/score", authHMACMiddleware(false), handlers.IncrementScoreHandle)
	app.Put("/comment/count", authHMACMiddleware(false), handlers.IncrementCommentHandle)
	app.Put("/comment/disable", append(hmacCookieHandlers, handlers.DisableCommentHandle)...)
	app.Put("/share/disable", append(hmacCookieHandlers, handlers.DisableSharingHandle)...)
	app.Put("/urlkey/:postId", append(hmacCookieHandlers, handlers.GeneratePostURLKeyHandle)...)
	app.Delete("/:postId", append(hmacCookieHandlers, handlers.DeletePostHandle)...)
	app.Get("/", append(hmacCookieHandlers, handlers.QueryPostHandle)...)
	app.Get("/:postId", append(hmacCookieHandlers, handlers.GetPostHandle)...)
	app.Get("/urlkey/:urlkey", append(hmacCookieHandlers, handlers.GetPostByURLKeyHandle)...)
}
