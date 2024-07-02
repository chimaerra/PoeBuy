package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"poebot/app"
	"syscall"
)

//go:embed config.yaml.example
var template string

func main() {

	if _, err := os.Stat("./config.yaml"); errors.Is(err, os.ErrNotExist) {
		os.WriteFile("config.yaml", []byte(template), 0644)
		fmt.Println("You should fill configuration file (config.yaml)\nPress CTRL-C to exit...")
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-stop
		return
	}

	// Create a new bot instance
	bot, err := app.Init()
	if err != nil {
		log.Fatalf("error app initialisation: %v", err)
	}

	// Start the bot
	err = bot.Run()
	if err != nil {
		log.Fatalf("error app run: %v", err)
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
