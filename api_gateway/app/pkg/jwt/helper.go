package jwt

import (
	"context"
	"encoding/json"
	"github.com/WTC-SYSTEM/apperror"
	"github.com/WTC-SYSTEM/logging"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/client/user_service"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/config"
	"github.com/cristalhq/jwt/v3"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"time"
)

var _ Helper = &helper{}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

type RT struct {
	RefreshToken string `json:"refresh_token"`
}

type helper struct {
	Logger logging.Logger
	DB     *redis.Client
}

func NewHelper(DB *redis.Client, logger logging.Logger) Helper {
	return &helper{DB: DB, Logger: logger}
}

type Helper interface {
	GenerateAccessToken(u user_service.User) ([]byte, error)
	UpdateRefreshToken(rt RT) ([]byte, error)
}

func (h *helper) UpdateRefreshToken(rt RT) ([]byte, error) {
	ctx := context.TODO()
	defer h.DB.Del(ctx, rt.RefreshToken)

	userBytes, err := h.DB.Get(ctx, rt.RefreshToken).Result()
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var u user_service.User
	err = json.Unmarshal([]byte(userBytes), &u)
	if err != nil {
		return nil, err
	}
	return h.GenerateAccessToken(u)
}

func (h *helper) GenerateAccessToken(u user_service.User) ([]byte, error) {
	key := []byte(config.GetConfig().JWT.Secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return nil, err
	}
	builder := jwt.NewBuilder(signer)

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        u.ID,
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		Email: u.Email,
	}
	token, err := builder.Build(claims)
	if err != nil {
		return nil, err
	}

	h.Logger.Info("create refresh token")
	refreshTokenUuid := uuid.New()
	userBytes, _ := json.Marshal(u)
	err = h.DB.Set(context.Background(), refreshTokenUuid.String(), userBytes, 0).Err()
	if err != nil {
		h.Logger.Error(err)
		return nil, err
	}

	jsonBytes, err := json.Marshal(map[string]string{
		"token":         token.String(),
		"refresh_token": refreshTokenUuid.String(),
	})
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
