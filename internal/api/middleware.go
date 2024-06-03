package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func errJson(ctx echo.Context, code int, message string) error {
	return ctx.JSON(code, echo.Map{
		"error": echo.Map{
			"code":    code,
			"message": message,
		},
	})
}

func (serv *Server) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authorizationHeader := ctx.Request().Header.Get(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			return errJson(ctx, http.StatusUnauthorized, "authorization header is not provided")
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return errJson(ctx, http.StatusUnauthorized, "invalid authorization header format")
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return errJson(ctx, http.StatusUnauthorized, fmt.Sprintf("unsupported authorization type %s", authorizationType))
		}

		accessToken := fields[1]
		payload, err := serv.token.Verify(accessToken)
		if err != nil {
			return errJson(ctx, http.StatusUnauthorized, "cannot verify access token")
		}

		ctx.Set(authorizationPayloadKey, payload)
		return next(ctx)
	}
}
