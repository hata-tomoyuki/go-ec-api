package auth

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var ErrInvalidClaims = errors.New("invalid token claims")

// UserID は JWT の "sub" クレームからユーザーIDを取得する。
// 全ハンドラーで共通して使える。
func UserID(r *http.Request) (int64, error) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, ErrInvalidClaims
	}
	return strconv.ParseInt(sub, 10, 64)
}

// LogoutClaims は Logout ハンドラーで必要なクレームをまとめて取得する。
type LogoutClaims struct {
	JTI            string
	ExpiredAt      time.Time
	RefreshTokenID int64
}

// GetLogoutClaims は JWT から Logout に必要な全クレームを取得する。
func GetLogoutClaims(r *http.Request) (LogoutClaims, error) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || token == nil {
		return LogoutClaims{}, ErrInvalidClaims
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return LogoutClaims{}, ErrInvalidClaims
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return LogoutClaims{}, ErrInvalidClaims
	}

	var refreshTokenID int64
	if rtidStr, ok := claims["rtid"].(string); ok {
		rtid, err := strconv.ParseInt(rtidStr, 10, 64)
		if err != nil {
			return LogoutClaims{}, ErrInvalidClaims
		}
		refreshTokenID = rtid
	}

	return LogoutClaims{
		JTI:            jti,
		ExpiredAt:      time.Unix(int64(exp), 0),
		RefreshTokenID: refreshTokenID,
	}, nil
}
