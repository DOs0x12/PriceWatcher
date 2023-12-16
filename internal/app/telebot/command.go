package telebot

type command struct {
	name        string
	description string
	action      func() string
}

var commands = []command{
	{
		name:        "/start",
		description: "Starts the session",
		action: func() string {
			return "The session is started!"
		},
	}, {
		name:        "/hello",
		description: "Say hello to the bot",
		action: func() string {
			return "Hello there!"
		},
	},
}
