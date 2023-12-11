package telegram

type GroupConf struct {
	Group  string `yaml:"group" json:"group"`
	ChatId int64  `yaml:"chatId" json:"chatId"`
}

type InfoConf GroupConf

type WarnConf GroupConf

type ErrorConf GroupConf

type SendConf struct {
	Info  InfoConf  `yaml:"info" json:"info"`
	Warn  WarnConf  `yaml:"warn" json:"warn"`
	Error ErrorConf `yaml:"error" json:"error"`
}
