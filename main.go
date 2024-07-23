package main

import (
	"log"
	"os"
	"os/signal"
	"poebuy/config"
	"poebuy/modules/bot"
	"poebuy/modules/ui"
	"poebuy/utils"
	"syscall"
)

func main() {

	logger := utils.NewLogger()
	defer logger.Close()

	// Load config from environment variables
	cfg, err := config.LoadConfig()
	if err != nil && err != config.ErrorNoConfigFile {
		logger.Errorf("config load failed: %v", err)
		return
	}

	ui.ShowUI(cfg, logger)
	return

	// Create a new bot instance
	bot, err := bot.NewBot(cfg)
	if err != nil {
		logger.Errorf("error app initialisation: %v", err)
		return
	}

	// Start the bot
	err = bot.Run()
	if err != nil {
		logger.Errorf("error app run: %v", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	log.Println("Shutting down...")

	// Cleanly close down
	bot.Stop()

}
