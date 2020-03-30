package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	GitHubOauthTokenUrl    = "https://github.com/login/oauth/access_token"
	GitHubOauthUserInfoUrl = "https://api.github.com/user"
	ClientId               = ""
	ClientSecret           = ""
)

func getToken(code string) (string, error) {
	client := &http.Client{}
	uri := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s", GitHubOauthTokenUrl, ClientId, ClientSecret, code)
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := map[string]string{}
	json.Unmarshal(body, &data)
	return data["access_token"], nil
}

func getUserInfo(token string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", GitHubOauthUserInfoUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Add("accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := map[string]interface{}{}
	json.Unmarshal(body, &data)
	return data, err
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/oauth/redirect", func(c *gin.Context) {
		code := c.Query("code")
		fmt.Printf("oauth2 code: %s\n", code)
		token, err := getToken(code)
		if err != nil {
			c.JSON(200, gin.H{"message": err.Error()})
			return
		}
		fmt.Printf("oauth2 token: %s\n", token)
		userInfo, err := getUserInfo(token)
		if err != nil {
			c.JSON(200, gin.H{"message": err.Error()})
			return
		}
		fmt.Printf("oauth2 user info: %v\n", userInfo)
		c.JSON(200, gin.H{"message": "ok", "token": token, "user_info": userInfo})
		return
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
