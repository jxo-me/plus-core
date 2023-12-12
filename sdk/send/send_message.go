package send

import "fmt"

type Message struct {
	SrvName string `json:"srv_name"`
	User    string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Msg     string `json:"msg"`
}

func (msg *Message) Format(level string) string {
	var operator = ""
	if msg.User != "" {
		operator = fmt.Sprintf("# Operator: %s", msg.User)
	}
	msg.Title = fmt.Sprintf("[%s]%s: %s", msg.SrvName, level, msg.Title)
	return fmt.Sprintf("<b>%s</b> \n\n# Time: %s\n# Content: %s\n# Message: %s \n%s",
		msg.Title, msg.Time, msg.Content, msg.Msg, operator)
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
