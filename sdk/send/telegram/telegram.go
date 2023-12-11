package telegram

import (
	"context"
	"fmt"
	telebot "github.com/jxo-me/gfbot"
)

const (
	Name = "Telegram"
)

type formatFunc func(level string, msg Message) string

type Telegram struct {
	Format formatFunc   `json:"format"`
	Cfg    *SendConf    `json:"cfg"`
	Client *telebot.Bot `json:"client"`
}

type Message struct {
	SrvName string `json:"srv_name"`
	UserId  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Msg     string `json:"msg"`
}

func NewTelegram(cfg *SendConf, client *telebot.Bot, f formatFunc) *Telegram {
	if f == nil {
		f = format
	}
	return &Telegram{
		Format: f,
		Cfg:    cfg,
		Client: client,
	}
}

func (t *Telegram) String() string {
	return Name
}

func format(level string, msg Message) string {
	var operator = ""
	if msg.UserId != 0 {
		operator = fmt.Sprintf("# Operator: %d", msg.UserId)
	}
	msg.Title = fmt.Sprintf("[%s]%s: %s", msg.SrvName, level, msg.Title)
	return fmt.Sprintf("<b>%s</b> \n\n# Time: %s\n# Content: %s\n# Message: %s \n%s",
		msg.Title, msg.Time, msg.Content, msg.Msg, operator)
}

func (t *Telegram) send(chatId int64, text string) (err error) {
	_, err = t.Client.Send(&telebot.User{ID: chatId}, text, &telebot.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             telebot.ModeHTML,
	})
	return err
}

func (t *Telegram) Info(ctx context.Context, msg Message) error {
	return t.send(t.Cfg.Info.ChatId, t.Format("Info", msg))
}

func (t *Telegram) Warn(ctx context.Context, msg Message) error {
	return t.send(t.Cfg.Warn.ChatId, t.Format("Warn", msg))
}

func (t *Telegram) Error(ctx context.Context, msg Message) error {
	return t.send(t.Cfg.Error.ChatId, t.Format("Error", msg))
}
