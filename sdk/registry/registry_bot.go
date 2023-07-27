package registry

import telebot "github.com/jxo-me/gfbot"

type BotRegistry struct {
	registry[*telebot.Bot]
}
