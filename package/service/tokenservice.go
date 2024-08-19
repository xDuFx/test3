package service

import (
	"encoding/base64"
	"test3/package/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	jwt.MapClaims
}

func NewAccessToken(ip, loginUser string) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS512)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["sub"] = loginUser
	claims["ip"] = ip

	return accessToken.SignedString([]byte(("TOP")))
}

func NewRefreshToken(loginUser, ip string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS512)
	claims := refreshToken.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	claims["sub"] = loginUser
	claims["ip"] = ip
	return refreshToken.SignedString([]byte(("SECRET")))
}

func ParseAccessToken(accessTokenHash string) (jwt.MapClaims, error) {
	accessToken, err := base64.StdEncoding.DecodeString(accessTokenHash)
	if err != nil {
		return jwt.MapClaims{}, err
	}
	parsedAccessToken, err := jwt.Parse(string(accessToken), func(token *jwt.Token) (interface{}, error) {
		return []byte("TOP"), nil
	})
	if err != nil || !parsedAccessToken.Valid {
		return jwt.MapClaims{}, err
	}

	return parsedAccessToken.Claims.(jwt.MapClaims), nil

}

func ParseRefreshToken(refreshTokenHash string) (jwt.MapClaims, error) {
	refreshToken, err := base64.StdEncoding.DecodeString(refreshTokenHash)
	if err != nil {
		return jwt.MapClaims{}, err
	}
	parsedRefreshToken, err := jwt.Parse(string(refreshToken), func(token *jwt.Token) (interface{}, error) {
		return []byte("SECRET"), nil
	})
	if err != nil || !parsedRefreshToken.Valid {
		return jwt.MapClaims{}, err
	}

	return parsedRefreshToken.Claims.(jwt.MapClaims), err
}
func CreateToken(loginUser, ip string) (*models.Token, error) {
	t, err := NewAccessToken(ip, loginUser)
	if err != nil {
		return &models.Token{}, err
	}
	ref, err := NewRefreshToken(loginUser, ip)
	if err != nil {
		return &models.Token{}, err
	}
	return &models.Token{
		AccessToken:  t,
		RefreshToken: ref,
	}, nil
}

func CompareToken(accessTokenHash, refreshTokenHash string) (bool, error) {
	claimsAccess, err := ParseAccessToken(accessTokenHash)
	if err != nil {
		return false, err
	}
	claimsRefresh, err := ParseRefreshToken(refreshTokenHash)
	if err != nil {
		return false, err
	}
	if claimsAccess["ip"] == claimsRefresh["ip"] && claimsAccess["sub"] == claimsRefresh["sub"] {
		return true, nil
	}
	return false, nil
}
