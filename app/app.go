package app

import (
	"fmt"
	"poebot/config"
	"poebot/connections"
	"poebot/handlers"
)

// App is the main application struct
type Bot struct {
	Config      *config.Config
	ItemHandler *handlers.ItemHandler
}

// Init initializes the application
func Init() (*Bot, error) {

	bot := &Bot{}

	// Load config from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config load failed: %v", err)
	}
	bot.Config = cfg

	conn, err := connections.NewWSConnection(bot.Config.User.Poesessid, "LQEV88GUn")
	if err != nil {
		return nil, fmt.Errorf("create ws listener failed: %v", err)
	}
	errChan := make(chan error)
	contChan := make(chan int)

	itemHandler := handlers.NewItemHandler(bot.Config.User.Poesessid, contChan, errChan, conn)
	bot.ItemHandler = itemHandler

	return bot, nil
}

// Run starts the application
func (bot *Bot) Run() error {

	go bot.ItemHandler.Serve()

	return nil
}

// Stop closes the application and cleans up
func (bot *Bot) Stop() {

	bot.ItemHandler.Close()

}
