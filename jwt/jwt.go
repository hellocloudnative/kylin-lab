package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserProfile 用于存储第三方平台的用户信息
type UserProfile struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	// 其他需要的字段 ...
}

// JWTClaims 用于定义JWT的载荷
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var SecretKey = []byte("lab")

// AuthMiddleware 是用于保护路由的中间件
func AuthMiddleware(c *gin.Context) {
	// 从Header中获取JWT令牌
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		c.Abort()
		return
	}

	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	// 将用户信息设置到Gin的上下文中
	claims := token.Claims.(*JWTClaims)
	c.Set("userProfile", UserProfile{UserID: claims.UserID, Username: claims.Username})
	c.Next()
}
