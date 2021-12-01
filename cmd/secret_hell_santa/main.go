package main

//import (
//	"fmt"
//	"github.com/truewebber/secretsantabot/internal/log"
//	"strings"
//	"time"
//
//	"github.com/go-telegram-bot-api/telegram-bot-api"
//	"github.com/mgutz/logxi/v1"
//
//	"lib/config"
//	"lib/game_factory"
//	"lib/model"
//	"lib/storage"
//)
//
//var (
//	conf     = config.Get()
//	selfName string
//)
//
//const (
//	EnrollCommand = "enroll"
//	EndCommand    = "end"
//	ListCommand   = "list"
//	MagicCommand  = "magic"
//	MyCommand     = "my"
//	HelpCommand   = "help"
//)
//
//func main() {
//	if conf.GetInt64("lock-on-chat-id") == 0 {
//		log.Error("lock-on-chat-id is required param")
//
//		return
//	}
//
//	bot, err := tgbotapi.NewBotAPI(conf.GetString("token"))
//	if err != nil {
//		log.Error("Error create new bot", "error", err.Error())
//
//		return
//	}
//
//	log.Debug("Authorized", "_", bot.Self.UserName)
//	selfName = bot.Self.UserName
//
//	u := tgbotapi.NewUpdate(20)
//	u.Timeout = 5
//
//	updates, err := bot.GetUpdatesChan(u)
//	if err != nil {
//		log.Error("Error get updates", "error", err.Error())
//
//		return
//	}
//
//	strg, err := storage.NewRedisStorage()
//	if err != nil {
//		log.Error("Error storage", "error", err.Error())
//
//		return
//	}
//
//	gameFactory := game_factory.New(strg)
//
//	for u := range updates {
//		if u.Message == nil {
//			log.Debug("NOT MSG")
//
//			continue
//		}
//
//		msg := u.Message
//
//		log.Debug("MSG", "chat-id", msg.Chat.ID, "chat", msg.Chat.Title,
//			"user-id", msg.From.ID, "user", msg.From.UserName, "_", msg.Text)
//
//		command := getCommand(msg.Text)
//
//		//lock to chat
//		if lockedWithChat(msg.Chat.ID) && command != MyCommand {
//			log.Error("NOT TRUST CHAT", "chat-id", msg.Chat.ID, "chat", msg.Chat.Title)
//
//			continue
//		}
//
//		switch command {
//		case EnrollCommand:
//			{
//				log.Debug("Enroll", "_", fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName))
//
//				err := gameFactory.Enroll(&model.HellMan{
//					TelegramId: msg.From.ID,
//					Username:   msg.From.UserName,
//					FirstName:  msg.From.FirstName,
//					LastName:   msg.From.LastName,
//					EnrollAt:   time.Now(),
//				})
//				if err != nil && err != game_factory.ErrorAlreadyEnroll && err != game_factory.ErrorMagicWasAlreadyDone {
//					log.Error("Error enroll user", "error", err.Error())
//
//					continue
//				} else if err == game_factory.ErrorMagicWasAlreadyDone {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "You can't enroll after magic was done.")
//					bot.Send(replyMsg)
//
//					continue
//				} else if err == game_factory.ErrorAlreadyEnroll {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Already enrolled.")
//					bot.Send(replyMsg)
//
//					continue
//				}
//
//				text := fmt.Sprintf("Congratulations!\n%s %s is having part in Secret Santa",
//					msg.From.FirstName, msg.From.LastName)
//				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
//				bot.Send(replyMsg)
//			}
//		case EndCommand:
//			{
//				log.Debug("End Enroll", "_", fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName))
//
//				err := gameFactory.DropEnroll(&model.HellMan{
//					TelegramId: msg.From.ID,
//					Username:   msg.From.UserName,
//					FirstName:  msg.From.FirstName,
//					LastName:   msg.From.LastName,
//					EnrollAt:   time.Now(),
//				})
//				if err != nil && err != game_factory.ErrorAlreadyEnroll && err != game_factory.ErrorMagicWasAlreadyDone {
//					log.Error("Error enroll user", "error", err.Error())
//
//					continue
//				} else if err == game_factory.ErrorMagicWasAlreadyDone {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "You can't leave game after magic was done.")
//					bot.Send(replyMsg)
//
//					continue
//				} else if err == game_factory.ErrorAlreadyEnroll {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "You are not in game.")
//					bot.Send(replyMsg)
//
//					continue
//				}
//
//				text := fmt.Sprintf("Sad to see you leaving =(\n%s %s is not in game from now.",
//					msg.From.FirstName, msg.From.LastName)
//				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
//				bot.Send(replyMsg)
//			}
//		case ListCommand:
//			{
//				log.Debug("List", "_", fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName))
//
//				list, err := gameFactory.ListEnrolled()
//				if err != nil {
//					log.Error("Error list enrolled", "error", err.Error())
//
//					continue
//				}
//
//				text := ""
//				if len(list) != 0 {
//					for i, u := range list {
//						text += fmt.Sprintf("%s %s", u.FirstName, u.LastName)
//						if u.Username != "" {
//							text += fmt.Sprintf(" (@%s)", u.Username)
//						}
//						if i != len(list)-1 {
//							text += "\n"
//						}
//					}
//				} else {
//					text = "List is empty..."
//				}
//
//				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
//				bot.Send(replyMsg)
//			}
//		case MagicCommand:
//			{
//				log.Debug("Magic", "_", fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName))
//
//				if lockedWithUser(msg.From.ID) {
//					log.Error("NOT ADMIN REQUEST", "user-id", msg.From.ID, "user", msg.From.UserName)
//
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "I sense great fear in you! Don't disturb me any more!")
//					bot.Send(replyMsg)
//
//					continue
//				}
//
//				result, err := gameFactory.Magic()
//				if err != nil && err != game_factory.ErrorMagicWasAlreadyDone && err != game_factory.ErrorNotEnough {
//					log.Error("Error MAGIC", "error", err.Error())
//
//					continue
//				} else if err == game_factory.ErrorMagicWasAlreadyDone {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "One magic was already done.")
//					bot.Send(replyMsg)
//
//					continue
//				} else if err == game_factory.ErrorNotEnough {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Not enough users enrolled to start game.")
//					bot.Send(replyMsg)
//
//					continue
//				}
//
//				for santa, man := range result {
//					text := fmt.Sprintf("Hi! Your target is `%s %s`.", man.FirstName, man.LastName)
//
//					replyMsg := tgbotapi.NewMessage(int64(santa.TelegramId), text)
//					_, err = bot.Send(replyMsg)
//					if err != nil {
//						log.Error("Error send magic private message", "_", err.Error())
//					}
//				}
//
//				text := "Magic done.\n\n" +
//					"In case you didn't receive message from me, write strait to me.\n" +
//					"Press on my name and press start on the bottom of window, then send me `/my` command."
//
//				replyMsg := tgbotapi.NewMessage(conf.GetInt64("lock-on-chat-id"), text)
//				_, err = bot.Send(replyMsg)
//				if err != nil {
//					log.Error("Error send message", "_", err.Error())
//				}
//
//			}
//		case MyCommand:
//			{
//				log.Debug("My", "_", fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName))
//				if msg.Chat.ID != int64(msg.From.ID) {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Please text me in private chat.")
//					replyMsg.ReplyToMessageID = msg.MessageID
//					bot.Send(replyMsg)
//
//					continue
//				}
//
//				man, err := gameFactory.GetMyMagic(&model.HellMan{
//					TelegramId: msg.From.ID,
//					Username:   msg.From.UserName,
//					FirstName:  msg.From.FirstName,
//					LastName:   msg.From.LastName,
//					EnrollAt:   time.Now(),
//				})
//				if err != nil && err != game_factory.ErrorMagicIsNotProceed && err != game_factory.ErrorYouAreNotInThisGame {
//					log.Error("Error get santa's man", "error", err.Error())
//
//					continue
//				} else if err == game_factory.ErrorMagicIsNotProceed {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Wait until magic will be done.")
//					bot.Send(replyMsg)
//
//					continue
//				} else if err == game_factory.ErrorYouAreNotInThisGame {
//					replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Sorry, you are not in this game.")
//					bot.Send(replyMsg)
//
//					continue
//				}
//
//				text := fmt.Sprintf("Hey! Your target is `%s %s`", man.FirstName, man.LastName)
//				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
//				bot.Send(replyMsg)
//			}
//		case HelpCommand:
//			{
//				text := "/enroll - enroll the game\n" +
//					"/end - stop your enroll (only before magic starts)\n" +
//					"/list - list all enrolling people\n" +
//					"/magic - start the game (only admin)\n" +
//					"/my - SecretHelSanta will resend magic info for you (only in private chat wi me)\n" +
//					"/help - show this message\n"
//
//				replyMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
//				bot.Send(replyMsg)
//			}
//		}
//	}
//}
//
//// ---------------------------------------------------------------------------------------------------------------------
//
//func getCommand(text string) string {
//	text = strings.TrimSuffix(text, fmt.Sprintf("@%s", selfName))
//	text = strings.Trim(text, "/")
//
//	return text
//}
//
//// ---------------------------------------------------------------------------------------------------------------------
//
//func lockedWithChat(chatId int64) bool {
//	lockChatId := conf.GetInt64("lock-on-chat-id")
//	if lockChatId == 0 {
//		return false
//	}
//
//	if chatId == lockChatId {
//		return false
//	}
//
//	return true
//}
//
//func lockedWithUser(userId int) bool {
//	lockUserId := conf.GetInt("admin-user-id")
//	if lockUserId == 0 {
//		return false
//	}
//
//	if userId == lockUserId {
//		return false
//	}
//
//	return true
//}
