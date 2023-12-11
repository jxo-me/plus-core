package send

import "fmt"

type Message struct {
	SrvName string `json:"srv_name"`
	UserId  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Msg     string `json:"msg"`
}

func (msg *Message) Format(level string) string {
	var operator = ""
	if msg.UserId != 0 {
		operator = fmt.Sprintf("# Operator: %d", msg.UserId)
	}
	msg.Title = fmt.Sprintf("[%s]%s: %s", msg.SrvName, level, msg.Title)
	return fmt.Sprintf("<b>%s</b> \n\n# Time: %s\n# Content: %s\n# Message: %s \n%s",
		msg.Title, msg.Time, msg.Content, msg.Msg, operator)
}
