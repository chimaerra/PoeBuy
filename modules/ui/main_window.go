package ui

import (
	"fmt"
	"image/color"
	"poebuy/config"
	"poebuy/modules/connections/models"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TappableLabel is a label that can be tapped and shows pointer cursor
type TappableLabel struct {
	widget.Label
	OnTapped func()
}

func NewTappableLabel(text string, tapped func()) *TappableLabel {
	t := &TappableLabel{
		Label: widget.Label{
			Text: text,
		},
		OnTapped: tapped,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *TappableLabel) Tapped(_ *fyne.PointEvent) {
	if t.OnTapped != nil {
		t.OnTapped()
	}
}

func (t *TappableLabel) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

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

	tradeTable := widget.NewTable(
		func() (int, int) {
			return len(cfg.Trade.Links), 5
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

			// Clear any existing tappable labels beyond the first two objects
			if len(container.Objects) > 2 {
				container.Objects = container.Objects[:2]
			}

			switch i.Col {
			case 0:
				label.Show()
				icon.Hide()
				label.SetText(cfg.Trade.Links[i.Row].Name)
			case 1: // Link column
				label.Hide()
				icon.Hide()
				
				// Create tappable label for copying
				tappable := NewTappableLabel(cfg.Trade.Links[i.Row].Code, func() {
					mw.Clipboard().SetContent(cfg.Trade.Links[i.Row].Code)
				})
				tappable.Wrapping = fyne.TextTruncate
				container.Objects = append(container.Objects, tappable)
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
			case 4:
				label.Show()
				icon.Hide()
				label.SetText(millisecondsToHumanReadable(cfg.Trade.Links[i.Row].Delay))
			}
		})
	mw.tradeTable = tradeTable
	tradeTable.Move(fyne.NewPos(15, 300))
	tradeTable.Resize(fyne.NewSize(765, 280))
	tradeTable.SetColumnWidth(0, 430)
	tradeTable.SetColumnWidth(1, 100)
	tradeTable.SetColumnWidth(2, 55)
	tradeTable.SetColumnWidth(3, 60)
	tradeTable.SetColumnWidth(4, 100)
	tradeTable.CreateHeader = func() fyne.CanvasObject {
		label := widget.NewLabel("")
		return label
	}
	tradeTable.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		label := o.(*widget.Label)
		switch id.Col {
		case 0:
			label.SetText("Name")
		case 1:
			label.SetText("Link")
		case 2:
			label.SetText("Active")
		case 3:
			label.SetText("Delete")
		case 4:
			label.SetText("Delay")
		}
	}
	tradeTable.ShowHeaderRow = true
	tradeTable.Refresh()

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

func millisecondsToHumanReadable(ms int64) string {

	if ms == 0 {
		return "no delay"
	}

	t := time.Duration(ms) * time.Millisecond

	switch {
	case t < time.Second:
		return fmt.Sprintf("%v ms", float32(ms))
	case t < time.Minute:
		return fmt.Sprintf("%v s", float32(ms)/1000)
	case t < time.Hour:
		return fmt.Sprintf("%v min", float32(ms)/1000/60)
	default:
		return fmt.Sprintf("%v h", float32(ms)/1000/60/60)
	}
}
