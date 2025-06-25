package handlers

import (
	"telegram-bot/auth"
	"telegram-bot/audio"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, chatID int64) {
	info, ok := auth.GetAuth(chatID)
	if !ok {
		StartLogin(bot, chatID)
		return
	}
	audio.SendNextAudio(bot, chatID, info)
}

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"sync"
// 	"telegram-bot/config"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// // --- Structlar ---
// type AudioSegment struct {
// 	AudioID   int    `json:"audio_id"`
// 	AudioName string `json:"audio_name"`
// 	CreatedAt string `json:"created_at"`
// 	FilePath  string `json:"file_path"`
// 	ID        int    `json:"id"`
// 	Status    string `json:"status"`
// }

// type AudioResponse struct {
// 	AudioSegments []AudioSegment `json:"audio_segments"`
// 	Count         int            `json:"count"`
// }

// type AuthRequest struct {
// 	Login    string `json:"login"`
// 	Password string `json:"password"`
// }

// type AuthResponse struct {
// 	AccessToken string `json:"access_token"`
// }

// type UserAuthInfo struct {
// 	Token  string `json:"token"`
// 	UserID string `json:"user_id"`
// }

// var chunkId int
// // var userStates = make(map[int64]string)
// var userTokens = make(map[int64]UserAuthInfo)
// var mu sync.Mutex

// const tokenFile = "user_tokens.json"

// func saveTokens() {
// 	data, _ := json.MarshalIndent(userTokens, "", "  ")
// 	os.WriteFile(tokenFile, data, 0644)
// }

// func loadTokens() {
// 	data, err := os.ReadFile(tokenFile)
// 	if err == nil {
// 		json.Unmarshal(data, &userTokens)
// 	}
// }

// func ensureToken(chatID int64, bot *tgbotapi.BotAPI) (string, error) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	// if token, ok := userTokens[chatID]; ok {
// 		// return token, nil
// 	// }

// 	bot.Send(tgbotapi.NewMessage(chatID, "🔐 Login:"))
// 	userStates[chatID] = "awaiting_login"
// 	return "", fmt.Errorf("token not found")
// }

// func sendNextAudio(bot *tgbotapi.BotAPI, chatID int64, token string) {
// 	url := "https://transcriber-bk.ccenter.uz/api/v1/audio_segment?user_id=185a1daa-06cf-4edb-89f5-c850ab8187cb"
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	req.Header.Set("accept", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		bot.Send(tgbotapi.NewMessage(chatID, "❌ HTTP error: "+err.Error()))
// 		return
// 	}
// 	if resp.StatusCode == 403 || resp.StatusCode == 401 {
// 		mu.Lock()
// 		delete(userTokens, chatID)
// 		saveTokens()
// 		mu.Unlock()
// 		bot.Send(tgbotapi.NewMessage(chatID, "🔁 Token expired. Login again"))
// 		ensureToken(chatID, bot)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	var result AudioResponse
// 	json.Unmarshal(body, &result)

// 	for _, seg := range result.AudioSegments {
// 		if seg.Status == "ready" {
// 			resp, err := http.Get(seg.FilePath)
// 			if err != nil {
// 				bot.Send(tgbotapi.NewMessage(chatID, "❌ Error while downloading audio: "+err.Error()))
// 				return
// 			}
// 			defer resp.Body.Close()

// 			wavPath := "temp.wav"
// 			mp3Path := "temp.mp3"

// 			wavFile, _ := os.Create(wavPath)
// 			io.Copy(wavFile, resp.Body)
// 			wavFile.Close()

// 			exec.Command("ffmpeg", "-y", "-i", wavPath, mp3Path).Run()
// 			mp3Data, _ := os.ReadFile(mp3Path)

// 			file := tgbotapi.FileBytes{
// 				Name:  fmt.Sprintf("%d.mp3", seg.ID),
// 				Bytes: mp3Data,
// 			}
// 			audio := tgbotapi.NewAudio(chatID, file)
// 			audio.Caption = fmt.Sprintf("🎧 %s", seg.AudioName)
// 			bot.Send(audio)

// 			chunkId = seg.ID
// 			os.Remove(wavPath)
// 			os.Remove(mp3Path)

// 			buttons := tgbotapi.NewInlineKeyboardMarkup(
// 				tgbotapi.NewInlineKeyboardRow(
// 					tgbotapi.NewInlineKeyboardButtonData("✅ Valid", "valid"),
// 					tgbotapi.NewInlineKeyboardButtonData("❌ Invalid", "invalid"),
// 				),
// 			)
// 			msg := tgbotapi.NewMessage(chatID, "✅ Choose audio status:")
// 			msg.ReplyMarkup = buttons
// 			bot.Send(msg)
// 			return
// 		}
// 	}
// 	bot.Send(tgbotapi.NewMessage(chatID, "🔍 Not found ready audio."))
// }

// func HandleUpdate(cfg *config.Config, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
// 	if update.CallbackQuery != nil {
// 		handleCallback(bot, update.CallbackQuery)
// 		return
// 	}

// 	if update.Message == nil {
// 		return
// 	}

// 	chatID := update.Message.Chat.ID
// 	msgText := update.Message.Text

// 	if state, ok := userStates[chatID]; ok {
// 		switch state {
// 		case "awaiting_login":
// 			userStates[chatID] = "awaiting_password"
// 			userStates[chatID+1000000000] = msgText
// 			bot.Send(tgbotapi.NewMessage(chatID, "🔑 Password:"))
// 			return
// 		case "awaiting_password":
// 			login := userStates[chatID+1000000000]
// 			password := msgText

// 			authPayload := AuthRequest{Login: login, Password: password}
// 			data, _ := json.Marshal(authPayload)
// 			req, _ := http.NewRequest("POST", "https://transcriber-bk.ccenter.uz/api/v1/auth/login", bytes.NewBuffer(data))
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("accept", "application/json")

// 			resp, err := http.DefaultClient.Do(req)
// 			if err != nil || resp.StatusCode != 200 {
// 				bot.Send(tgbotapi.NewMessage(chatID, "❌ Error during authorization"))
// 				delete(userStates, chatID)
// 				return
// 			}
// 			defer resp.Body.Close()
// 			body, _ := io.ReadAll(resp.Body)
// 			var authResp AuthResponse
// 			json.Unmarshal(body, &authResp)

// 			mu.Lock()
// 			// userTokens[chatID] = authResp.AccessToken
// 			mu.Unlock()
// 			saveTokens()

// 			bot.Send(tgbotapi.NewMessage(chatID, "✅ Successfully logged in!"))
// 			delete(userStates, chatID)
// 			go sendNextAudio(bot, chatID, authResp.AccessToken)
// 			return
// 		}
// 	}

// 	if msgText == "/start" {
// 		loadTokens()
// 		token, err := ensureToken(chatID, bot)
// 		if err != nil {
// 			return
// 		}

// 		url := "https://transcriber-bk.ccenter.uz/api/v1/audio_segment?user_id=185a1daa-06cf-4edb-89f5-c850ab8187cb"
// 		req, _ := http.NewRequest("GET", url, nil)
// 		req.Header.Set("Authorization", "Bearer "+token)
// 		req.Header.Set("accept", "application/json")

// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			bot.Send(tgbotapi.NewMessage(chatID, "❌ HTTP error: "+err.Error()))
// 			return
// 		}
// 		if resp.StatusCode == 403 || resp.StatusCode == 401 {
// 			delete(userTokens, chatID)
// 			saveTokens()
// 			bot.Send(tgbotapi.NewMessage(chatID, "🔁 Token expired. Login again"))
// 			ensureToken(chatID, bot)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		body, _ := io.ReadAll(resp.Body)
// 		var result AudioResponse
// 		json.Unmarshal(body, &result)

// 		for _, seg := range result.AudioSegments {
// 			if seg.Status == "ready" {
// 				resp, err := http.Get(seg.FilePath)
// 				if err != nil {
// 					bot.Send(tgbotapi.NewMessage(chatID, "❌ Error while downloading audio: "+err.Error()))
// 					return
// 				}
// 				defer resp.Body.Close()

// 				wavPath := "temp.wav"
// 				mp3Path := "temp.mp3"

// 				wavFile, _ := os.Create(wavPath)
// 				io.Copy(wavFile, resp.Body)
// 				wavFile.Close()

// 				exec.Command("ffmpeg", "-y", "-i", wavPath, mp3Path).Run()
// 				mp3Data, _ := os.ReadFile(mp3Path)

// 				file := tgbotapi.FileBytes{
// 					Name:  fmt.Sprintf("%d.mp3", seg.ID),
// 					Bytes: mp3Data,
// 				}
// 				audio := tgbotapi.NewAudio(chatID, file)
// 				audio.Caption = fmt.Sprintf("🎧 %s", seg.AudioName)
// 				bot.Send(audio)

// 				chunkId = seg.ID
// 				os.Remove(wavPath)
// 				os.Remove(mp3Path)

// 				buttons := tgbotapi.NewInlineKeyboardMarkup(
// 					tgbotapi.NewInlineKeyboardRow(
// 						tgbotapi.NewInlineKeyboardButtonData("✅ Valid", "valid"),
// 						tgbotapi.NewInlineKeyboardButtonData("❌ Invalid", "invalid"),
// 					),
// 				)
// 				msg := tgbotapi.NewMessage(chatID, "✅ Choose audio status:")
// 				msg.ReplyMarkup = buttons
// 				bot.Send(msg)
// 				return
// 			}
// 		}
// 		bot.Send(tgbotapi.NewMessage(chatID, "🔍 Not found ready audio."))
// 		return
// 	}

// 	if state, ok := userStates[chatID]; ok {
// 		delete(userStates, chatID)
// 		// token := userTokens[chatID]

// 		payload := map[string]interface{}{
// 			"entire_audio_invalid": state == "invalid",
// 		}
// 		if state == "valid" {
// 			payload["transcribe_text"] = msgText
// 		} else {
// 			payload["report_text"] = msgText
// 		}

// 		data, _ := json.Marshal(payload)
// 		url := fmt.Sprintf("https://transcriber-bk.ccenter.uz/api/v1/transcript/update?id=%d", chunkId)

// 		req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
// 		// req.Header.Set("Authorization", "Bearer "+token)
// 		req.Header.Set("Content-Type", "application/json")
// 		req.Header.Set("accept", "application/json")

// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			bot.Send(tgbotapi.NewMessage(chatID, "❌ Error request: "+err.Error()))
// 			return
// 		}
// 		defer resp.Body.Close()

// 		bot.Send(tgbotapi.NewMessage(chatID, "✅ Saved."))
// 		// go sendNextAudio(bot, chatID, token)
// 		return
// 	}
// }

// func handleCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
// 	chatID := cq.Message.Chat.ID
// 	data := cq.Data

// 	bot.Request(tgbotapi.NewCallback(cq.ID, ""))

// 	if data == "valid" || data == "invalid" {
// 		userStates[chatID] = data
// 		bot.Send(tgbotapi.NewMessage(chatID, "✏️ Write a transcript:"))
// 	}
// }
