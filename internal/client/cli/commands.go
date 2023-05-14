package cli

import "github.com/chzyer/readline"

var authCompleters = readline.NewPrefixCompleter(
	readline.PcItem(registerCommand),
	readline.PcItem(loginCommand),

	readline.PcItem(exitCommand),
	readline.PcItem(helpCommand),
)

var secretsCompleters = readline.NewPrefixCompleter(
	readline.PcItem(secretsCommand),

	readline.PcItem(createCommand),
	readline.PcItem(getCommand+" <number>"),
	// readline.PcItem(updateCommand+" <number>"),
	readline.PcItem(deleteCommand+" <number>"),

	readline.PcItem(exitCommand),
	readline.PcItem(helpCommand),
)

const (
	exitCommand = "exit"
	helpCommand = "help"

	registerCommand = "register"
	loginCommand    = "login"

	secretsCommand = "secrets"

	createCommand = "create"
	getCommand    = "get"
	// updateCommand = "update"
	deleteCommand = "delete"
)
