package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey = "asupersecretekeythatshouldntbehardcodedhere"
)

type ITokenService interface {
	CreateToken(user User) (string, error)
	ValidateAccessToken(token string) (UserSafe, error)
	GenerateTokenPairForUser(user User) (TokenPair, error)
	ValidateRefreshToken(refreshToken string) (string, error)
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenService struct{}

var TS ITokenService

func (t *TokenService) GenerateTokenPairForUser(user User) (TokenPair, error) {
	tokenPair := TokenPair{}
	accessToken, err := t.CreateToken(user)
	if err != nil {
		log.Println("Error creating access token: ", err)
		return tokenPair, err
	}
	refreshToken, err := t.CreateRefreshTokenForUser(user)
	if err != nil {
		log.Println("Error creating refresh token: ", err)
		return tokenPair, err
	}
	tokenPair.AccessToken = accessToken
	tokenPair.RefreshToken = refreshToken

	return tokenPair, nil
}

// Creates a new Token for a User
func (t *TokenService) CreateToken(user User) (string, error) {
	if (User{}) == user {
		return "", fmt.Errorf("cannot create token for empty user")
	}
	safeUser := UserSafe{}
	safeUser.ID = user.ID
	safeUser.Email = user.Email
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": safeUser,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// TODO: Remove hardcoded signed string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Error signing token: ", err)
		return "", err
	}
	return tokenString, nil
}

func (t *TokenService) CreateRefreshTokenForUser(user User) (string, error) {
	if (User{}) == user {
		return "", fmt.Errorf("cannot create token for empty user")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("error signing token", err)
		return "", err
	}
	return tokenString, nil
}

func (t *TokenService) ValidateRefreshToken(refreshToken string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if token.Valid {
		if reflect.TypeOf(claims["sub"]) != reflect.TypeOf("") {
			return "", fmt.Errorf("malformed sub string")
		}
		userId := claims["sub"].(string)
		return userId, nil
	}
	return "", nil
}

func (t *TokenService) ValidateAccessToken(tokenString string) (UserSafe, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return UserSafe{}, err
	}
	if token.Valid {
		//TODO: get rid of the sub - parse what we expect to see
		var userSafe UserSafe
		if reflect.TypeOf(claims["sub"]) != reflect.TypeOf(map[string]interface{}{}) {
			return UserSafe{}, fmt.Errorf("malformed sub string")
		}
		data := claims["sub"].(map[string]interface{})
		userSafe.ID = data["id"].(string)
		userSafe.Email = data["email"].(string)
		if err != nil {
			return UserSafe{}, fmt.Errorf("malformed sub string")
		}
		return userSafe, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return UserSafe{}, fmt.Errorf("malformed Token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return UserSafe{}, fmt.Errorf("token is Expired")
		} else {
			return UserSafe{}, fmt.Errorf("could not handle the token")
		}
	} else {
		return UserSafe{}, fmt.Errorf("error handling token: %v", err)
	}
}
