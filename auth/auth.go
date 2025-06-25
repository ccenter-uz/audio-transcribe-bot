package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
)

const tokenFile = "user_tokens.json"

type UserAuthInfo struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

var (
	userTokens = make(map[int64]UserAuthInfo)
	mu         sync.Mutex
)

func LoadTokens() {
	data, err := os.ReadFile(tokenFile)
	if err == nil {
		_ = json.Unmarshal(data, &userTokens)
	}
}

func SaveTokens() {
	data, _ := json.MarshalIndent(userTokens, "", "  ")
	_ = os.WriteFile(tokenFile, data, 0644)
}

func GetAuth(chatID int64) (UserAuthInfo, bool) {
	mu.Lock()
	defer mu.Unlock()
	val, ok := userTokens[chatID]
	return val, ok
}

func GetToken(chatID int64) string {
	mu.Lock()
	defer mu.Unlock()
	if val, ok := userTokens[chatID]; ok {
		return val.Token
	}
	return ""
}

func GetUserID(chatID int64) string {
	mu.Lock()
	defer mu.Unlock()
	if val, ok := userTokens[chatID]; ok {
		return val.UserID
	}
	return ""
}

func SetAuth(chatID int64, info UserAuthInfo) {
	mu.Lock()
	defer mu.Unlock()
	userTokens[chatID] = info
	SaveTokens()
}

func RemoveAuth(chatID int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(userTokens, chatID)
	SaveTokens()
}

func LoginUser(login, password string) (UserAuthInfo, error) {
	type authReq struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	type authResp struct {
		AccessToken string `json:"access_token"`
		User        struct {
			ID string `json:"agent_id"`
		} `json:"user"`
	}

	reqBody, _ := json.Marshal(authReq{Login: login, Password: password})
	req, _ := http.NewRequest("POST", "https://transcriber-bk.ccenter.uz/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return UserAuthInfo{}, errors.New("login failed")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var parsed authResp
	_ = json.Unmarshal(body, &parsed)

	return UserAuthInfo{
		Token:  parsed.AccessToken,
		UserID: parsed.User.ID,
	}, nil
}
