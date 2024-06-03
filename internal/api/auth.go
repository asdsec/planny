package api

import (
	"com.github/asdsec/planny/internal/model"
	"com.github/asdsec/planny/internal/security"
	db "com.github/asdsec/planny/internal/store"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (serv *Server) register(ctx echo.Context) error {
	var req registerRequest
	if err := ctx.Bind(&req); err != nil {
		return serv.err(ctx, http.StatusBadRequest, "invalid request")
	}

	password, err := security.HashPassword(req.Password)
	if err != nil {
		return serv.err(ctx, http.StatusInternalServerError, "failed to hash password")
	}

	arg := db.CreateUserArg{
		Username:       req.Username,
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		HashedPassword: password,
	}
	user, err := serv.store.CreateUser(ctx.Request().Context(), arg)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return serv.err(ctx, http.StatusBadRequest, "username or email already exists")
		}
		return serv.err(ctx, http.StatusInternalServerError, "failed to create user")
	}

	return ctx.JSON(http.StatusOK, newRegisterResponse(&user))
}

type (
	registerRequest struct {
		Username  string `json:"username" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		FirstName string `json:"first_name" validate:"required,min=1,max=32"`
		LastName  string `json:"last_name" validate:"required,min=1,max=32"`
		Password  string `json:"password" validate:"required,min=6,max=32"`
	}

	registerResponse struct {
		Username    string    `json:"username"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedAt   time.Time `json:"created_at"`
	}
)

func newRegisterResponse(user *model.User) *registerResponse {
	return &registerResponse{
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.FirstName + " " + user.LastName,
		UpdatedAt:   user.UpdatedAt,
		CreatedAt:   user.CreatedAt,
	}
}

func (serv *Server) login(ctx echo.Context) error {
	var req loginRequest
	if err := ctx.Bind(&req); err != nil {
		return serv.err(ctx, http.StatusBadRequest, "invalid request")
	}
	if req.Username == "" && req.Email == "" {
		return serv.err(ctx, http.StatusBadRequest, "username or email is required")
	}

	var user model.User
	var err error
	if req.Username != "" {
		user, err = serv.store.GetUserByUsername(ctx.Request().Context(), req.Username)
		if err != nil {
			if err.Error() == db.ErrRecordNotFound {
				return serv.err(ctx, http.StatusNotFound, "user not found")
			}
			return serv.err(ctx, http.StatusInternalServerError, "failed to get user by username")
		}
	} else {
		user, err = serv.store.GetUserByEmail(ctx.Request().Context(), req.Email)
		if err != nil {
			if err.Error() == db.ErrRecordNotFound {
				return serv.err(ctx, http.StatusNotFound, "user not found")
			}
			return serv.err(ctx, http.StatusInternalServerError, "failed to get user by email")
		}
	}

	err = security.CheckPassword(req.Password, user.Password)
	if err != nil {
		return serv.err(ctx, http.StatusUnauthorized, "incorrect password")
	}

	accessToken, accessPayload, err := serv.token.Generate(
		user.ID,
		user.Username,
		serv.conf.AccessTokenDuration,
	)
	if err != nil {
		return serv.err(ctx, http.StatusInternalServerError, "failed to generate access token")
	}
	refreshToken, refreshPayload, err := serv.token.Generate(
		user.ID,
		user.Username,
		serv.conf.RefreshTokenDuration,
	)
	if err != nil {
		return serv.err(ctx, http.StatusInternalServerError, "failed to generate refresh token")
	}

	arg := db.CreateSessionArg{
		ID:           refreshPayload.ID,
		UserID:       refreshPayload.UserID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request().UserAgent(),
		ClientIp:     ctx.RealIP(),
		ExpiresAt:    refreshPayload.ExpiresAt,
	}
	session, err := serv.store.CreateSession(ctx.Request().Context(), arg)
	if err != nil {
		return serv.err(ctx, http.StatusInternalServerError, "failed to create session")
	}

	credentials := loginCredentials{
		SessionID:             session.ID.String(),
		AccessToken:           accessToken,
		ExpiresAt:             accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
	}
	return ctx.JSON(http.StatusOK, newLoginResponse(&user, &credentials))
}

type (
	loginRequest struct {
		Username string `json:"username" validate:"omitempty"`
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"required"`
	}

	loginCredentials struct {
		SessionID             string    `json:"session_id"`
		AccessToken           string    `json:"access_token"`
		ExpiresAt             time.Time `json:"expires_at"`
		RefreshToken          string    `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	}

	loginUser struct {
		ID          uint      `json:"id"`
		Username    string    `json:"username"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedAt   time.Time `json:"created_at"`
	}

	loginResponse struct {
		Credentials loginCredentials `json:"credentials"`
		User        loginUser        `json:"user"`
	}
)

func newLoginResponse(user *model.User, credentials *loginCredentials) *loginResponse {
	return &loginResponse{
		Credentials: *credentials,
		User: loginUser{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.FirstName + " " + user.LastName,
			UpdatedAt:   user.UpdatedAt,
			CreatedAt:   user.CreatedAt,
		},
	}
}

func (serv *Server) renewAccess(ctx echo.Context) error {
	var req renewAccessRequest
	if err := ctx.Bind(&req); err != nil {
		return serv.err(ctx, http.StatusBadRequest, "invalid request")
	}

	refreshPayload, err := serv.token.Verify(req.RefreshToken)
	if err != nil {
		return serv.err(ctx, http.StatusUnauthorized, "invalid refresh token")
	}

	session, err := serv.store.GetSessionByRefreshToken(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return serv.err(ctx, http.StatusNotFound, "session not found")
		}
		return serv.err(ctx, http.StatusInternalServerError, "failed to get session by refresh token")
	}
	if session.Username != refreshPayload.Username {
		return serv.err(ctx, http.StatusUnauthorized, "incorrect session user")
	}
	if session.RefreshToken != req.RefreshToken {
		return serv.err(ctx, http.StatusUnauthorized, "incorrect refresh token")
	}
	if time.Now().After(session.ExpiresAt) {
		return serv.err(ctx, http.StatusUnauthorized, "session expired")
	}

	accessToken, accessPayload, err := serv.token.Generate(
		refreshPayload.UserID,
		refreshPayload.Username,
		serv.conf.AccessTokenDuration,
	)
	if err != nil {
		return serv.err(ctx, http.StatusInternalServerError, "failed to generate access token")
	}

	rsp := renewAccessResponse{
		SessionID:             session.ID.String(),
		AccessToken:           accessToken,
		ExpiresAt:             accessPayload.ExpiresAt,
		RefreshToken:          req.RefreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
	}
	return ctx.JSON(http.StatusOK, rsp)
}

type (
	renewAccessRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	renewAccessResponse loginCredentials
)
