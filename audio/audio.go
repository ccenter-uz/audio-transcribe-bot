package audio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"telegram-bot/auth"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AudioSegment struct {
	AudioID   int    `json:"audio_id"`
	AudioName string `json:"audio_name"`
	CreatedAt string `json:"created_at"`
	FilePath  string `json:"file_path"`
	ID        int    `json:"id"`
	Status    string `json:"status"`
}

type AudioResponse struct {
	AudioSegments []AudioSegment `json:"audio_segments"`
	Count         int            `json:"count"`
}

var userChunkIDs = make(map[int64]int)
var userStates = make(map[int64]string)

func HandleCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	data := cq.Data

	bot.Request(tgbotapi.NewCallback(cq.ID, "")) // Loadingni yopish

	switch data {
	case "start_transcribe", "next_audio":
		if info, ok := auth.GetAuth(chatID); ok {
			go SendNextAudio(bot, chatID, info)
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "❌ Token topilmadi. /start buyrug‘i orqali qayta login qiling."))
		}
	case "valid", "invalid":
		userStates[chatID] = data
		bot.Send(tgbotapi.NewMessage(chatID, "✏️ Matnni kiriting:"))
	}
}

func HandleText(bot *tgbotapi.BotAPI, chatID int64, msgText string) {
	state, ok := userStates[chatID]
	if !ok {
		return
	}
	delete(userStates, chatID)

	info, ok := auth.GetAuth(chatID)
	if !ok {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Token topilmadi. Iltimos, qayta login qiling."))
		return
	}

	payload := map[string]interface{}{
		"entire_audio_invalid": state == "invalid",
	}
	if state == "valid" {
		payload["transcribe_text"] = msgText
	} else {
		payload["report_text"] = msgText
	}

	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://transcriber-bk.ccenter.uz/api/v1/transcript/update?id=%d", userChunkIDs[chatID])

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+info.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		return
	}
	defer resp.Body.Close()

	nextBtn := tgbotapi.NewMessage(chatID, "✅ Saqlandi. Davom ettirasizmi?")
	nextBtn.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➡️ Next Audio", "next_audio"),
		),
	)
	bot.Send(nextBtn)
}

func SendNextAudio(bot *tgbotapi.BotAPI, chatID int64, info auth.UserAuthInfo) {
	url := fmt.Sprintf("https://transcriber-bk.ccenter.uz/api/v1/audio_segment?user_id=%s", info.UserID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+info.Token)
	req.Header.Set("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ HTTP error: "+err.Error()))
		return
	}
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		auth.RemoveAuth(chatID)
		bot.Send(tgbotapi.NewMessage(chatID, "🔁 Token topilmadi. /start buyrug‘i orqali qayta login qiling."))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result AudioResponse
	json.Unmarshal(body, &result)

	for _, seg := range result.AudioSegments {
		if seg.Status == "ready" {
			err := downloadAndConvertAndSend(bot, chatID, seg)
			if err != nil {
				slog.Error("Error downloading or converting audio", "error", err)
			}
			userChunkIDs[chatID] = seg.ID
			return
		}
	}
	bot.Send(tgbotapi.NewMessage(chatID, "🔍 Not found ready audio"))
}

func downloadAndConvertAndSend(bot *tgbotapi.BotAPI, chatID int64, seg AudioSegment) error {
	resp, err := http.Get(seg.FilePath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	wavPath := "temp.wav"
	mp3Path := "temp.mp3"

	wavFile, _ := os.Create(wavPath)
	io.Copy(wavFile, resp.Body)
	wavFile.Close()

	exec.Command("ffmpeg", "-y", "-i", wavPath, mp3Path).Run()
	mp3Data, err := os.ReadFile(mp3Path)
	if err != nil {
		return err
	}

	file := tgbotapi.FileBytes{
		Name:  fmt.Sprintf("%d.mp3", seg.ID),
		Bytes: mp3Data,
	}
	audio := tgbotapi.NewAudio(chatID, file)
	audio.Caption = fmt.Sprintf("🎧 %s", seg.AudioName)
	bot.Send(audio)

	os.Remove(wavPath)
	os.Remove(mp3Path)

	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Transcript yozish", "valid"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Xabar berish", "invalid"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "✅ Audio holati:")
	msg.ReplyMarkup = buttons
	bot.Send(msg)

	return nil
}