package mdw

import (
	"strings"
	myJwt "ytb/src/controller/jwt"

	"github.com/labstack/echo/v4"
)

func IsValidToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" || token == "null" || token == "undefined" {
			return echo.ErrUnauthorized
		} else {
			token = strings.Split(token, " ")[1]
			if _, err := myJwt.ParseATToken(token, c); err != nil {
				// log.Println("IsValidToken", err)
				return echo.ErrUnauthorized
			} else {
				return next(c)
			}
		}
	}
}
