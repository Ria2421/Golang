package middleware

import (
	"login-go/auth"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(jwter auth.IJwtParser) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := jwter.SetAuthToContext(c); err != nil {
				return err
			}

			return next(c)
		}
	}
}
