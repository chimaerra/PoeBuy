package ui

import (
	"image/color"
	"poebuy/config"
	"poebuy/modules/connections/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MainWindow struct {
	fyne.Window

	leagueSelect   *widget.Select
	nameEntry      *widget.Entry
	linkEntry      *widget.Entry
	addTradeButton *widget.Button
	tradeTable     *widget.Table
}

func NewMainWindow(app fyne.App, info *models.TradeInfo, cfg *config.Config) *MainWindow {

	mw := &MainWindow{}

	mw.Window = app.NewWindow("PoeBuy")
	mw.SetFixedSize(true)
	mw.Resize(fyne.NewSize(800, 600))

	leagueLabel := widget.NewLabel("League:")
	leagueLabel.Move(fyne.NewPos(15, 10))

	leagueSelect := widget.NewSelect(info.GetLeagues(), func(s string) { cfg.Trade.League = s })
	mw.leagueSelect = leagueSelect
	leagueSelect.Move(fyne.NewPos(15, 50))
	leagueSelect.Resize(fyne.NewSize(350, 35))
	leagueSelect.PlaceHolder = "Select legaue"
	leagueSelect.SetSelected(cfg.Trade.League)
	leagueSelect.Refresh()

	nicknameLabel := widget.NewLabel("Logged in as " + info.Nickname)
	nicknameLabel.Move(fyne.NewPos(630-float32(len(info.Nickname)*7), 10))
	nicknameLabel.TextStyle = fyne.TextStyle{Bold: true}

	addTradeLabel := widget.NewLabel("Add trade links:")
	addTradeLabel.Move(fyne.NewPos(15, 100))

	nameEntry := widget.NewEntry()
	mw.nameEntry = nameEntry
	nameEntry.Move(fyne.NewPos(15, 140))
	nameEntry.Resize(fyne.NewSize(765, 40))
	nameEntry.PlaceHolder = "Trade name"
	nameEntry.Refresh()

	linkEntry := widget.NewEntry()
	mw.linkEntry = linkEntry
	linkEntry.Move(fyne.NewPos(15, 190))
	linkEntry.Resize(fyne.NewSize(765, 40))
	linkEntry.PlaceHolder = "Trade link or code"
	linkEntry.Refresh()

	addTradeRectangle := canvas.NewRectangle(nil)
	addTradeRectangle.Move(fyne.NewPos(10, 135))
	addTradeRectangle.Resize(fyne.NewSize(775, 152))
	addTradeRectangle.StrokeWidth = 1
	addTradeRectangle.StrokeColor = color.RGBA{R: 60, G: 60, B: 60, A: 255}
	addTradeRectangle.CornerRadius = 5
	addTradeRectangle.Refresh()

	addTradeButton := widget.NewButton("Add", nil)
	mw.addTradeButton = addTradeButton
	addTradeButton.Move(fyne.NewPos(580, 240))
	addTradeButton.Resize(fyne.NewSize(200, 40))

	// tradeTable := widget.NewTableWithHeaders(nil, nil, nil)
	tradeTable := widget.NewTable(
		func() (int, int) {
			return len(cfg.Trade.Links), 4
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			icon := widget.NewIcon(nil)
			icon.Hide()
			return container.NewStack(label, icon)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {

			container := o.(*fyne.Container)
			label := container.Objects[0].(*widget.Label)
			icon := container.Objects[1].(*widget.Icon)

			switch i.Col {
			case 0:
				label.Show()
				icon.Hide()
				label.SetText(cfg.Trade.Links[i.Row].Name)
			case 1:
				label.Show()
				icon.Hide()
				label.SetText(cfg.Trade.Links[i.Row].Code)
			case 2:
				label.Hide()
				icon.Show()
				if cfg.Trade.Links[i.Row].IsActiv {
					icon.SetResource(theme.CheckButtonCheckedIcon())
				} else {
					icon.SetResource(theme.CheckButtonIcon())
				}
			case 3:
				label.Hide()
				icon.Show()
				icon.SetResource(theme.DeleteIcon())
			}
		})
	mw.tradeTable = tradeTable
	tradeTable.Move(fyne.NewPos(15, 300))
	tradeTable.Resize(fyne.NewSize(765, 280))
	tradeTable.SetColumnWidth(0, 500)
	tradeTable.SetColumnWidth(1, 100)
	tradeTable.SetColumnWidth(2, 30)
	tradeTable.SetColumnWidth(3, 30)

	mw.SetContent(container.NewWithoutLayout(
		leagueLabel,
		leagueSelect,
		nicknameLabel,
		addTradeLabel,
		nameEntry,
		linkEntry,
		addTradeRectangle,
		addTradeButton,
		tradeTable,
	))

	return mw
}

func (w *MainWindow) OnAddTrade(f func()) {
	w.addTradeButton.OnTapped = f
}

func (w *MainWindow) OnTableCellClick(f func(id widget.TableCellID)) {
	w.tradeTable.OnSelected = f
}
