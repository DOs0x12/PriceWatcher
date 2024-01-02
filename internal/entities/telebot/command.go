package telebot

type Command struct {
	Name        string
	Description string
	Action      func() string
}

type CommandWithInput struct {
	Name        string
	Description string
	Action      func(input string) string
}
