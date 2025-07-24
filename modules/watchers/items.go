package watchers

import (
	"net/http"
	"poebuy/modules/connections"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type ItemWatcher struct {
	WSConnection        *websocket.Conn
	Fetcher             *connections.Fetcher
	Whisper             *connections.Whisper
	Code                string
	ErrChan             chan error
	Working             bool
	Delay               time.Duration
	readReady           bool
	index               int
	UpdateCheckmarkFunc func(int)
}

func NewItemWatcher(poesseid string, league string, code string, errChan chan error, delay int64, index int, updateCheckmarkFunc func(int), soundFile string) (*ItemWatcher, error) {

	client := &http.Client{}

	wsConn, err := connections.NewWSConnection(poesseid, league, code)
	if err != nil {
		return nil, err
	}

	watcher := &ItemWatcher{
		WSConnection:        wsConn,
		Fetcher:             connections.NewFetcher(client, headers.GetFetchitemHeaders(poesseid)),
		Whisper:             connections.NewWhisper(client, headers.GetWhisperHeaders(poesseid), soundFile),
		Code:                code,
		ErrChan:             errChan,
		Working:             false,
		Delay:               time.Millisecond * time.Duration(delay),
		readReady:           true,
		index:               index,
		UpdateCheckmarkFunc: updateCheckmarkFunc,
	}

	return watcher, nil

}

func (w *ItemWatcher) Watch() {

	w.Working = true

	if w.Delay > 0 {
		go w.delayer()
	}

	for {
		if !w.Working {
			return
		}
		var ls models.LivesearchNewItem
		err := w.WSConnection.ReadJSON(&ls)
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				w.ErrChan <- err
			}
			break
		}
		if w.Delay > 0 && !w.readReady {
			continue
		}

		length := len(ls.New)
		if w.Delay != 0 && length > 3 {
			length = 3
		}

		itemsInfo, err := w.Fetcher.FetchItems(ls.New[:length], w.Code)
		if err != nil {
			w.ErrChan <- err
			continue
		}

		for _, itemInfo := range itemsInfo {
			err := w.Whisper.Whisper(itemInfo.Result[0].Listing.WhisperToken)
			if err != nil {
				w.ErrChan <- err
				continue
			}
		}

		w.readReady = false
	}

	w.Stop()
}

func (w *ItemWatcher) Stop() {

	w.WSConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	w.WSConnection.Close()
	w.Working = false
	w.UpdateCheckmarkFunc(w.index)
}

func (w *ItemWatcher) delayer() {
	for {
		if !w.Working {
			return
		}
		w.readReady = true
		time.Sleep(w.Delay)
	}
}
