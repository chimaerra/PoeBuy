package main

import (
	"log"
	"os"
	"os/signal"
	"poebot/app"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {

	time.Sleep(time.Second)

	robotgo.MoveSmooth(100, 100)
	return
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
