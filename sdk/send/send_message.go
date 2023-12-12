package send

import "fmt"

const DefaultTpl = "<b>%s</b> \n\n# Time: %s\n# Content: %s\n# Message: %s \n%s"

type Message struct {
	Template string `json:"template"`
	SrvName  string `json:"srv_name"`
	User     string `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Time     string `json:"time"`
	Msg      string `json:"msg"`
}

func (msg *Message) Format(level string) string {
	if msg.User != "" {
		msg.User = fmt.Sprintf("# Operator: %s", msg.User)
	}
	if msg.Title != "" && level != "" {
		msg.Title = fmt.Sprintf("[%s]%s: %s", msg.SrvName, level, msg.Title)
	}
	if msg.Template == "" {
		msg.Template = DefaultTpl
	}

	return fmt.Sprintf(msg.Template, msg.GetTitle(), msg.GetTime(), msg.GetContent(), msg.GetMsg(), msg.GetUser())
}

func (msg *Message) SetTemplate(tpl string) {
	msg.Template = tpl
}

func (msg *Message) GetTitle() string {
	return msg.Title
}

func (msg *Message) GetUser() string {
	return msg.User
}

func (msg *Message) GetTime() string {
	return msg.Time
}

func (msg *Message) GetContent() string {
	return msg.Content
}

func (msg *Message) GetMsg() string {
	return msg.Msg
}
