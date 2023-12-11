package telegram

import (
	"context"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/plus-core/core/v2/send"
)

const (
	Name = "Telegram"
)

type Telegram struct {
	Cfg    *SendConf    `json:"cfg"`
	Client *telebot.Bot `json:"client"`
}

func NewTelegram(cfg *SendConf, client *telebot.Bot) *Telegram {
	return &Telegram{
		Cfg:    cfg,
		Client: client,
	}
}

func (t *Telegram) String() string {
	return Name
}

func (t *Telegram) send(chatId int64, text string) (err error) {
	_, err = t.Client.Send(&telebot.User{ID: chatId}, text, &telebot.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             telebot.ModeHTML,
	})
	return err
}

func (t *Telegram) Info(ctx context.Context, msg send.ISendMsg) error {
	return t.send(t.Cfg.Info.ChatId, msg.Format("Info"))
}

func (t *Telegram) Warn(ctx context.Context, msg send.ISendMsg) error {
	return t.send(t.Cfg.Warn.ChatId, msg.Format("Warn"))
}

func (t *Telegram) Error(ctx context.Context, msg send.ISendMsg) error {
	return t.send(t.Cfg.Error.ChatId, msg.Format("Error"))
}
