package ui

import (
	"fmt"
	"poebuy/config"
	"poebuy/modules/bot"
	"poebuy/modules/connections"
	"poebuy/modules/connections/models"
	"poebuy/resources"
	"poebuy/utils"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

type UI struct {
	app             fyne.App
	mainWindow      *MainWindow
	poesessidwindow *PoessidWindow
	delayWindow     *DelayWindow
	cfg             *config.Config
	info            *models.TradeInfo
	bot             *bot.Bot
}

const _appId = "com.kelaron.poebuy"

func ShowUI(cfg *config.Config, logger *utils.Logger, bot *bot.Bot) {

	ui := &UI{
		cfg: cfg,
		bot: bot,
	}

	app := app.NewWithID(_appId)
	app.Settings().SetTheme(theme.DarkTheme())
	app.SetIcon(resources.ResourceDivineIco)
	ui.app = app

	if cfg.General.Poesessid == "" {
		ui.ShowPoessidWindow()
	} else {
		info, err := connections.GetTradeInfo(ui.cfg.General.Poesessid)
		if err != nil && err != connections.ErrorBadPoessid {
			logger.Error(err.Error())
			return
		}
		if err == connections.ErrorBadPoessid {
			ui.ShowPoessidWindow()
		} else {
			ui.info = info
			ui.ShowMainWindow()
		}
	}
	bot.UpdateCheckmarkFunc = ui.updateCheckmark
	ui.app.Run()
}

func (ui *UI) updateCheckmark(row int) {
	ui.cfg.Trade.Links[row].IsActiv = false
	fyne.Do(func() { ui.mainWindow.tradeTable.RefreshItem(widget.TableCellID{Col: 2, Row: row}) })
}

func (ui *UI) ShowPoessidWindow() {
	ui.poesessidwindow = NewPoessidWindow(ui.app)
	ui.poesessidwindow.OnConfirmPoessid(ui.savePoessid)
	ui.poesessidwindow.OnClose(ui.Close)
	ui.poesessidwindow.Show()
}

func (ui *UI) ShowMainWindow() {
	ui.mainWindow = NewMainWindow(ui.app, ui.info, ui.cfg)
	ui.mainWindow.SetOnClosed(ui.closeApp)
	ui.mainWindow.OnAddTrade(ui.addTrade)
	ui.mainWindow.OnTableCellClick(ui.tableCellClick)
	ui.mainWindow.Show()
}

func (ui *UI) ShowDelayWindow(delay int64, linkId int) {
	ui.delayWindow = NewDelayWindow(ui.app, delay, linkId)
	ui.delayWindow.OnConfirmDelay(ui.saveDelay)
	ui.delayWindow.Show()
}

func (ui *UI) savePoessid() {
	info, err := connections.GetTradeInfo(ui.poesessidwindow.poesessidEntry.Text)
	if err != nil {
		dialog.Message("Wrong POESSID: %v", err).Title("Ooops!").Error()
		ui.poesessidwindow.poesessidEntry.SetText("")
		return
	}
	ui.info = info
	ui.cfg.General.Poesessid = ui.poesessidwindow.poesessidEntry.Text
	ui.ShowMainWindow()
	ui.poesessidwindow.Close()
}

func (ui *UI) Close() {
	ui.app.Quit()
}

func (ui *UI) addTrade() {

	var inputLink string

	if strings.Contains(ui.mainWindow.linkEntry.Text, "pathofexile.com") {
		inputLink = regexp.MustCompile("[A-Za-z0-9-_]+$").FindString(ui.mainWindow.linkEntry.Text)
	} else {
		inputLink = ui.mainWindow.linkEntry.Text
	}

	for _, link := range ui.cfg.Trade.Links {
		if link.Code == inputLink {
			dialog.Message("This link has already been added").Title("PoeBuy").Info()
			return
		}
	}
	ui.cfg.Trade.Links = append(ui.cfg.Trade.Links, config.Link{Name: ui.mainWindow.nameEntry.Text, Code: inputLink})
	err := ui.cfg.Save()
		if err != nil {
		dialog.Message(fmt.Sprintf("Error saving config: %v", err)).Title("Error").Error()
	}
	ui.mainWindow.nameEntry.SetText("")
	ui.mainWindow.linkEntry.SetText("")
}

func (ui *UI) tableCellClick(id widget.TableCellID) {

	ui.mainWindow.tradeTable.Unselect(id)

	if id.Row < 0 {
		return
	}

	switch id.Col {
	case 2:
		if ui.cfg.Trade.Links[id.Row].IsActiv {
			ui.bot.StopWatcher(ui.cfg.Trade.Links[id.Row].Code)
			ui.cfg.Trade.Links[id.Row].IsActiv = false
		} else {
			err := ui.bot.WatchItem(ui.cfg.Trade.Links[id.Row].Code, ui.cfg.Trade.Links[id.Row].Delay)
			if err != nil {
				dialog.Message("Link connection error:\n%v", err).Title("Live search error").Error()
				return
			}
			ui.cfg.Trade.Links[id.Row].IsActiv = true
		}
	case 3:
		if ui.cfg.Trade.Links[id.Row].IsActiv {
			ui.bot.StopWatcher(ui.cfg.Trade.Links[id.Row].Code)
			ui.cfg.Trade.Links[id.Row].IsActiv = false
		}
		ui.cfg.Trade.Links = append(ui.cfg.Trade.Links[:id.Row], ui.cfg.Trade.Links[id.Row+1:]...)
		err := ui.cfg.Save()
			if err != nil {
				dialog.Message(fmt.Sprintf("Error saving config: %v", err)).Title("Error").Error()
		}
	case 4:
		ui.ShowDelayWindow(ui.cfg.Trade.Links[id.Row].Delay, id.Row)
	default:
		return
	}
	ui.mainWindow.tradeTable.Refresh()
}

func (ui *UI) closeApp() {
    if err := ui.cfg.Save(); err != nil {
        // Log error - adjust as needed
        println("Error saving config on app close:", err.Error())
    }
	ui.bot.StopAllWatchers()
}

func (ui *UI) saveDelay() {
	if ui.delayWindow.delayEntry.Validate() != nil {
		dialog.Message("Enter valid delay value").Title("Error").Error()
		return
	}
	delay, _ := strconv.Atoi(ui.delayWindow.delayEntry.Text)
	ui.cfg.Trade.Links[ui.delayWindow.linkID].Delay = int64(delay)
	err := ui.cfg.Save()
		if err != nil {
			dialog.Message(fmt.Sprintf("Error saving config: %v", err)).Title("Error").Error()
		}
	ui.delayWindow.Close()
}
