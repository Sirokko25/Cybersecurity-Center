package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const tokenExp = time.Hour * 3
const secretKey = "supersecretkey"

func BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Неопознанный метод шифрования: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		fmt.Println("Токен не валиден.")
		return ""
	}

	return claims.UserID
}

func WithAuth(c *gin.Context) {
	jwtAuth, err := c.Cookie("jwt_auth")
	if err != nil {
		if err == http.ErrNoCookie {
			token, err := BuildJWTString(uuid.NewString())
			if err != nil {
				c.JSON(http.StatusInternalServerError, "Ошибка при создании JWT-токена")
				return
			}
			cookie := &http.Cookie{
				Path:   "/",
				Name:   "jwt_auth",
				Value:  token,
				MaxAge: 300,
			}
			c.Writer.Header().Set("user-id-auth", GetUserID(token))
			c.Writer.Header().Set("is-new-user", "true")
			c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, "localhost", false, true)
			log.Info().Str("user_id", GetUserID(token)).Msg("")
			c.Next()
			return
		}
		c.JSON(http.StatusInternalServerError, "Ошибка при создании куки")
		return
	}

	c.Writer.Header().Set("user-id-auth", GetUserID(jwtAuth))
	log.Info().Str("user_id", GetUserID(jwtAuth)).Msg("Аунтификация прошла успешно.")
	c.Next()
}
