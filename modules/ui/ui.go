package ui

import (
	"poebuy/config"
	"poebuy/modules/connections"
	"poebuy/modules/connections/models"
	"poebuy/resources"
	"poebuy/utils"
	"regexp"
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
	cfg             *config.Config
	info            *models.TradeInfo
}

const _appId = "com.kelaron.poebuy"

func ShowUI(cfg *config.Config, logger *utils.Logger) {

	ui := &UI{}

	ui.cfg = cfg

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
	ui.app.Run()
}

func (ui *UI) ShowPoessidWindow() {
	ui.poesessidwindow = NewPoessidWindow(ui.app)
	ui.poesessidwindow.OnConfirmPoessid(ui.savePoessid)
	ui.poesessidwindow.OnClose(ui.Close)
	ui.poesessidwindow.Show()
}

func (ui *UI) ShowMainWindow() {
	ui.mainWindow = NewMainWindow(ui.app, ui.info, ui.cfg)
	ui.mainWindow.SetOnClosed(ui.cfg.Save)
	ui.mainWindow.OnAddTrade(ui.addTrade)
	ui.mainWindow.OnTableCellClick(ui.tableCellClick)
	ui.mainWindow.Show()
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
	ui.mainWindow.nameEntry.SetText("")
	ui.mainWindow.linkEntry.SetText("")
}

func (ui *UI) tableCellClick(id widget.TableCellID) {
	ui.mainWindow.tradeTable.Unselect(id)

	switch id.Col {
	case 2:
		if ui.cfg.Trade.Links[id.Row].IsActiv {
			ui.cfg.Trade.Links[id.Row].IsActiv = false
		} else {
			ui.cfg.Trade.Links[id.Row].IsActiv = true
		}
	case 3:
		ui.cfg.Trade.Links = append(ui.cfg.Trade.Links[:id.Row], ui.cfg.Trade.Links[id.Row+1:]...)
	default:
		return
	}
}
