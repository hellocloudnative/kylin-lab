package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"time"
)

// CustomClaims 是一个自定义结构，用于存储JWT的声明
type CustomClaims struct {
	// 可以添加任何你需要的声明
	jwt.StandardClaims
}

// GetJWTSecret 从环境变量中获取JWT密钥
func GetJWTSecret() []byte {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return []byte(jwtSecret)
}

// GenerateToken 生成一个带有时效的token
func GenerateToken() (string, error) {
	// 设置token的过期时间
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// JWTMiddleware 是一个Gin中间件，用于验证token
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "Authorization header is missing"})
			c.Abort()
			return
		}
		claims := &CustomClaims{}
		// 解析token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return GetJWTSecret(), nil
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Malformed token"})
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "Token expired"})
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "Token not active yet"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Could not handle this token"})
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": err.Error()})
			}
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "Token is invalid"})
			c.Abort()
			return
		}
		// 如果token有效，继续处理请求
		c.Next()
	}
}
