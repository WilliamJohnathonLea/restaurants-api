package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/WilliamJohnathonLea/restaurants-api/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authorizationHeaderKey = "Authorization"

	bearerScheme = "bearer"
)

type AuthData struct {
	Scheme string
	Token  string
}

func Authenticated(tokenKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHdr := ctx.GetHeader(authorizationHeaderKey)
		if authHdr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		auth, err := getSchemeAndToken(authHdr)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if auth.Scheme != bearerScheme {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parsedTkn, err := jwt.ParseWithClaims(auth.Token, &types.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(tokenKey), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := parsedTkn.Claims.(*types.UserClaims); ok {
			ctx.Set("user_claims", *claims)
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func getSchemeAndToken(authHeader string) (AuthData, error) {
	data := AuthData{}
	strSlice := strings.Split(authHeader, " ")

	if len(strSlice) != 2 {
		return data, errors.New("authorization header is incorrectly formatted")
	}

	data.Scheme = strings.ToLower(strSlice[0])
	data.Token = strSlice[1]

	return data, nil
}
