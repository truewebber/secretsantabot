package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mgutz/logxi/v1"

	"lib/config"
	"lib/game_factory"
	"lib/model"
	"lib/storage"
)

var (
	conf = config.Get()
)

func main() {
	bot, err := tgbotapi.NewBotAPI(conf.GetString("token"))
	if err != nil {
		log.Error("Error create new bot", "error", err.Error())

		return
	}

	log.Debug("Authorized", "_", bot.Self.UserName)

	u := tgbotapi.NewUpdate(20)
	u.Timeout = 5

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error("Error get updates", "error", err.Error())

		return
	}

	strg, err := storage.NewRedisStorage()
	if err != nil {
		log.Error("Error storage", "error", err.Error())

		return
	}

	gameFactory := game_factory.New(strg)

	for u := range updates {
		if u.Message == nil {
			log.Debug("NOT MSG")

			continue
		}

		msg := u.Message

		//lock to chat
		if lockedWithChat(msg.Chat.ID) {
			log.Error("NOT TRUST CHAT", "chat-id", msg.Chat.ID, "chat", msg.Chat.Title)

			continue
		}

		log.Debug("MSG", "chat-id", msg.Chat.ID, "chat", msg.Chat.Title,
			"user-id", msg.From.ID, "user", msg.From.UserName, "_", msg.Text)

		cmdText := strings.TrimSuffix(msg.Text, fmt.Sprintf("@%s", model.TGBotName))

		switch cmdText {
		case "/start":
			{
				log.Debug("Enroll", "_", fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName))

				err := gameFactory.Enroll(&model.HellMan{
					TelegramId: msg.From.ID,
					Username:   msg.From.UserName,
					FirstName:  msg.From.FirstName,
					LastName:   msg.From.LastName,
					EnrollAt:   time.Now(),
				})
				if err != nil && err != game_factory.ErrorAlreadyEnroll && err != game_factory.ErrorMagicWasAlreadyDone {
					log.Error("Error enroll user", "error", err.Error())

					continue
				} else if err == game_factory.ErrorMagicWasAlreadyDone {
					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "You can't enroll after magic was done.")
					bot.Send(replyMsg)

					continue
				} else if err == game_factory.ErrorAlreadyEnroll {
					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Already enrolled.")
					bot.Send(replyMsg)

					continue
				}

				text := fmt.Sprintf("Congratulations!\n%s %s is having part in Secret Santa",
					msg.From.FirstName, msg.From.LastName)
				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				bot.Send(replyMsg)
			}
		case "/end":
			{
				log.Debug("End Enroll", "_", fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName))

				err := gameFactory.DropEnroll(&model.HellMan{
					TelegramId: msg.From.ID,
					Username:   msg.From.UserName,
					FirstName:  msg.From.FirstName,
					LastName:   msg.From.LastName,
					EnrollAt:   time.Now(),
				})
				if err != nil && err != game_factory.ErrorAlreadyEnroll && err != game_factory.ErrorMagicWasAlreadyDone {
					log.Error("Error enroll user", "error", err.Error())

					continue
				} else if err == game_factory.ErrorMagicWasAlreadyDone {
					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "You can't leave game after magic was done.")
					bot.Send(replyMsg)

					continue
				} else if err == game_factory.ErrorAlreadyEnroll {
					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "You are not in game.")
					bot.Send(replyMsg)

					continue
				}

				text := fmt.Sprintf("Sad to see you leaving =(\n%s %s is not in game from now.",
					msg.From.FirstName, msg.From.LastName)
				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				bot.Send(replyMsg)
			}
		case "/list":
			{
				list, err := gameFactory.ListEnrolled()
				if err != nil {
					log.Error("Error list enrolled", "error", err.Error())

					continue
				}

				text := ""
				for i, u := range list {
					text += fmt.Sprintf("%s %s", u.FirstName, u.LastName)
					if u.Username != "" {
						text += fmt.Sprintf(" (@%s)", u.Username)
					}
					if i != len(list)-1 {
						text += "\n"
					}
				}

				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				bot.Send(replyMsg)
			}
		case "/magic":
			{
				if lockedWithUser(msg.From.ID) {
					log.Error("NOT ADMIN REQUEST", "user-id", msg.From.ID, "user", msg.From.UserName)

					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "I feel your fear! Don't disturb me any more!")
					bot.Send(replyMsg)

					continue
				}

				result, err := gameFactory.Magic()
				if err != nil && err != game_factory.ErrorMagicWasAlreadyDone {
					log.Error("Error MAGIC", "error", err.Error())

					continue
				} else if err == game_factory.ErrorMagicWasAlreadyDone {
					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "One magic was already done.")
					bot.Send(replyMsg)

					continue
				}

				text := "TEST RESULTS:\n\n"
				for santa, man := range result {
					text += fmt.Sprintf("%s to %s\n", santa.FirstName, man.FirstName)
				}

				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				bot.Send(replyMsg)
			}
		case "/help":
			{
				text := "/start - enroll the game\n" +
					"/end - stop your enroll (only before magic starts)\n" +
					"/magic - start the game (only admin)\n" +
					"/list - list all enrolling people\n" +
					"/help - show this message\n"

				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				bot.Send(replyMsg)
			}
		}
	}
}

func lockedWithChat(chatId int64) bool {
	lockChatId := conf.GetInt64("lock-on-chat-id")
	if lockChatId == 0 {
		return false
	}

	if chatId == lockChatId {
		return false
	}

	return true
}

func lockedWithUser(userId int) bool {
	lockUserId := conf.GetInt("admin-user-id")
	if lockUserId == 0 {
		return false
	}

	if userId == lockUserId {
		return false
	}

	return true
}
