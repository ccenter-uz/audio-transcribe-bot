package handlers

// import (
// 	"bufio"
// 	"log"
// 	"log/slog"
// 	"os"
// 	"regexp"
// 	"strings"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// var (
// 	youtubeRegex  = regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/[^\s]+$`)
// 	transcriptDir = "transcripts"
// 	linksFile     = transcriptDir + "/links.txt"
// )

// func YTLinks(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
// 	os.MkdirAll(transcriptDir, 0755)
// 	if _, err := os.Stat(linksFile); os.IsNotExist(err) {
// 		os.WriteFile(linksFile, []byte(""), 0644)
// 	}

// 	for update := range updates {
// 		if update.Message == nil {
// 			continue
// 		}

// 		msg := update.Message
// 		text := strings.TrimSpace(msg.Text)

// 		if strings.HasPrefix(text, "#report") && msg.ReplyToMessage != nil {
// 			replyText := strings.TrimSpace(msg.ReplyToMessage.Text)
// 			if youtubeRegex.MatchString(replyText) {
// 				err := removeLinkFromFile(replyText)
// 				if err != nil {
// 					slog.Error("Error while delete link")
// 				} else {
// 					delMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.ReplyToMessage.MessageID)
// 					_, delErr := bot.Request(delMsg)
// 					if delErr != nil {
// 						log.Printf("❌ Xabarni o‘chirishda xatolik: %v", delErr)
// 					}
// 					slog.Info("Link o‘chirildi")
// 				}
// 			}
// 			continue
// 		}

// 		if youtubeRegex.MatchString(text) {
// 			if isLinkAlreadySaved(text) {
// 				slog.Info("Link allaqachon saqlangan", "link", text)
// 			} else {
// 				appendLinkToFile(text)
// 				slog.Info("Link saqlandi", "link", text)
// 			}
// 		} else {
// 			del := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
// 			bot.Request(del)
// 		}
// 	}
// }

// func isLinkAlreadySaved(link string) bool {
// 	file, err := os.Open(linksFile)
// 	if err != nil {
// 		log.Printf("❌ Fayl o‘qishda xatolik: %v", err)
// 		return false
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(link) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func appendLinkToFile(link string) {
// 	file, err := os.OpenFile(linksFile, os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Printf("❌ Yozishda xatolik: %v", err)
// 		return
// 	}
// 	defer file.Close()

// 	_, err = file.WriteString(link + "\n")
// 	if err != nil {
// 		log.Printf("❌ Link yozishda xatolik: %v", err)
// 	}
// }

// func removeLinkFromFile(link string) error {
// 	file, err := os.ReadFile(linksFile)
// 	if err != nil {
// 		return err
// 	}

// 	lines := strings.Split(string(file), "\n")
// 	var newLines []string
// 	for _, line := range lines {
// 		if strings.TrimSpace(line) != strings.TrimSpace(link) && line != "" {
// 			newLines = append(newLines, line)
// 		}
// 	}

// 	return os.WriteFile(linksFile, []byte(strings.Join(newLines, "\n")+"\n"), 0644)
// }
