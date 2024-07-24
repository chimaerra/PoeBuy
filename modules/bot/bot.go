package bot

import (
	"poebuy/config"
	"poebuy/modules/watchers"
	"poebuy/utils"
)

// App is the main application struct
type Bot struct {
	config   *config.Config
	Watchers []*watchers.ItemWatcher
	stopChan chan string
	ErrChan  chan error
	logger   *utils.Logger
}

// Init initializes the application
func NewBot(cfg *config.Config, logger *utils.Logger) (*Bot, error) {

	bot := &Bot{
		stopChan: make(chan string),
		ErrChan:  make(chan error),
		config:   cfg,
		logger:   logger,
	}

	go bot.errorWriter()

	return bot, nil
}

func (bot *Bot) WatchItem(code string) error {
	watcher, err := watchers.NewItemWatcher(bot.config.General.Poesessid, bot.config.Trade.League, code, bot.stopChan, bot.ErrChan)
	if err != nil {
		return err
	}

	bot.Watchers = append(bot.Watchers, watcher)

	go watcher.Watch()

	return nil
}

func (bot *Bot) StopWatcher(code string) {
	bot.stopChan <- code
}

// Stop closes the application and cleans up
func (bot *Bot) StopAllWatchers() {

	for _, watcher := range bot.Watchers {
		watcher.Stop()
	}

}

func (bot *Bot) errorWriter() {
	for {
		err := <-bot.ErrChan
		bot.logger.Error(err.Error())
	}
}
