package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const contextHospitalKey = "hospital"

func JWTAuth(secret string) gin.HandlerFunc {
	key := []byte(secret)
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if !strings.HasPrefix(authz, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenInvalidClaims
			}
			return key, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		hospitalAny, ok := claims["hospital"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "hospital not found in token"})
			return
		}
		hospital, ok := hospitalAny.(string)
		if !ok || strings.TrimSpace(hospital) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "hospital not found in token"})
			return
		}

		c.Set(contextHospitalKey, hospital)
		c.Next()
	}
}

func HospitalFromContext(c *gin.Context) string {
	v, ok := c.Get(contextHospitalKey)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
