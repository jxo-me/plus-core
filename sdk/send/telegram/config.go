package telegram

type GroupConf struct {
	ChatId int64 `yaml:"chatId" json:"chatId"`
}

type InfoConf GroupConf

type WarnConf GroupConf

type ErrorConf GroupConf

type SendConf struct {
	Group string    `yaml:"group" json:"group"`
	Info  InfoConf  `yaml:"info" json:"info"`
	Warn  WarnConf  `yaml:"warn" json:"warn"`
	Error ErrorConf `yaml:"error" json:"error"`
}
