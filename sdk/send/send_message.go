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
	if msg.Template == "" {
		msg.Template = DefaultTpl
	}
	return fmt.Sprintf(msg.Template, msg.GetTitle(level), msg.GetTime(), msg.GetContent(), msg.GetMsg(), msg.GetUser())
}

func (msg *Message) GetTitle(level string) string {
	return fmt.Sprintf("[%s]%s: %s", msg.SrvName, level, msg.Title)
}

func (msg *Message) GetUser() string {
	return fmt.Sprintf("# Operator: %s", msg.User)
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
