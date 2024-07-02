package app

import (
	"fmt"
	"poebot/config"
	"poebot/modules/connections"
	"poebot/modules/handlers"
)

// App is the main application struct
type Bot struct {
	Config      *config.Config
	ItemHandler []*handlers.ItemHandler
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

	errChan := make(chan error)
	contChan := make(chan int)

	for _, link := range cfg.Trade.Links {

		conn, err := connections.NewWSConnection(bot.Config.User.Poesessid, link)
		if err != nil {
			return nil, fmt.Errorf("create ws listener failed: %v", err)
		}

		itemHandler := handlers.NewItemHandler(bot.Config.User.Poesessid, contChan, errChan, conn)
		bot.ItemHandler = append(bot.ItemHandler, itemHandler)
	}

	return bot, nil
}

// Run starts the application
func (bot *Bot) Run() error {

	for _, handler := range bot.ItemHandler {

		go handler.Serve()

	}

	return nil
}

// Stop closes the application and cleans up
func (bot *Bot) Stop() {

	for _, handler := range bot.ItemHandler {

		handler.Close()

	}

}
