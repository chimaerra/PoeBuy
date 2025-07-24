package bot

import (
	"poebuy/config"
	"poebuy/modules/watchers"
	"poebuy/utils"
)

// App is the main application struct
type Bot struct {
	config              *config.Config
	Watchers            map[string]*watchers.ItemWatcher
	ErrChan             chan error
	logger              *utils.Logger
	UpdateCheckmarkFunc func(int)
}

// Init initializes the application
func NewBot(cfg *config.Config, logger *utils.Logger) (*Bot, error) {

	bot := &Bot{
		ErrChan:  make(chan error),
		config:   cfg,
		logger:   logger,
		Watchers: make(map[string]*watchers.ItemWatcher),
	}

	go bot.errorWriter()

	return bot, nil
}

func (bot *Bot) WatchItem(code string, delay int64) error {

	var index int

	for i, _ := range bot.config.Trade.Links {
		if bot.config.Trade.Links[i].Code == code {
			index = i
			break
		}
	}

	soundFile := bot.config.General.SoundFile
	watcher, err := watchers.NewItemWatcher(
        bot.config.General.Poesessid,
        bot.config.Trade.League,
        code,
        bot.ErrChan,
        delay,
        index,
        bot.UpdateCheckmarkFunc,
        soundFile,
    )
	if err != nil {
		return err
	}

	bot.Watchers[code] = watcher

	go watcher.Watch()

	return nil
}

func (bot *Bot) StopWatcher(code string) {
	bot.Watchers[code].Stop()
	delete(bot.Watchers, code)
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
