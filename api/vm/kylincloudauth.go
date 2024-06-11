package vm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func PostAuthRequestToken() (string, error) {

	var authRequest AuthRequest
	authRequest.Username = "system"
	authRequest.Password = "sys@cloud_"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 创建HTTP客户端
	client := &http.Client{
		Timeout:   time.Second * 10, // 设置超时时间
		Transport: tr,
	}

	// 定义请求的URL
	url := "https://10.44.61.74/api/auth"

	// 将AuthRequest结构体转换为JSON格式的请求体
	jsonBody, err := json.Marshal(authRequest)
	if err != nil {
		return "", err
	}

	// 创建HTTP POST请求
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	// 设置请求头部，指定发送的数据是JSON格式
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to authenticate: %v", resp.StatusCode)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应体，提取token
	var authResponse AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return "", err
	}

	// 返回获取到的token
	return authResponse.Data.Token, nil
}

func GetKylinCloudToken(c *gin.Context) {
	// 调用函数从指定的HTTPS URL获取token
	token, err := PostAuthRequestToken()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token"})
		return
	}

	// 将获取到的token返回给客户端
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": token}})
}
