package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	jwt2 "kylin-lab/jwt"
	"kylin-lab/models"
	"net/http"
	"strconv"
	"time"
)

type UserInfoResponse struct {
	Data struct {
		Token    string   `json:"token"`
		UserInfo UserInfo `json:"user_info"`
	} `json:"data"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func LabToken(c *gin.Context) {
	var creds jwt2.UserCredentials
	creds.Username = "system"
	creds.Password = "sys@cloud_"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Second * 10, // 设置超时时间
		Transport: tr,
	}

	url := "https://10.44.61.74/api/auth"

	jsonBody, err := json.Marshal(creds)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	var userInfoResponse UserInfoResponse
	if err := json.Unmarshal(body, &userInfoResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// 创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt2.JWTClaims{
		UserID:   strconv.Itoa(userInfoResponse.Data.UserInfo.ID),
		Username: userInfoResponse.Data.UserInfo.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "YourIssuer",
		},
	})

	// 签名令牌
	tokenString, err := token.SignedString(jwt2.SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	if userInfoResponse.Data.Token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token not found in response"})
		return
	}

	var userprofile models.LabUser
	userprofile.UserId = userInfoResponse.Data.UserInfo.ID
	userprofile.Username = userInfoResponse.Data.UserInfo.Username

	_, err = userprofile.Insert()
	if err != nil {
		log.Warning(err)
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
