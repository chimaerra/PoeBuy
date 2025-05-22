package main

import (
	"poebuy/config"
	"poebuy/modules/bot"
	"poebuy/modules/ui"
	"poebuy/utils"
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

	// Create a new bot instance
	botInstance, err := bot.NewBot(cfg, logger)
	if err != nil {
		logger.Errorf("error app initialisation: %v", err)
		return
	}

	ui.ShowUI(cfg, logger, botInstance)

}
