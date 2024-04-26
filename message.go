package glu

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func NewMessage(role Role, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}

func SystemMessage(content string) Message {
	return NewMessage(RoleSystem, content)
}

func UserMessage(content string) Message {
	return NewMessage(RoleUser, content)
}

func AssistantMessage(content string) Message {
	return NewMessage(RoleAssistant, content)
}
