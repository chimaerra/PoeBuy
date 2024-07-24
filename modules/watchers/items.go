package watchers

import (
	"net/http"
	"poebuy/modules/connections"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"
	"strings"

	"github.com/gorilla/websocket"
)

type ItemWatcher struct {
	WSConnection *websocket.Conn
	Fetcher      *connections.Fetcher
	Whisper      *connections.Whisper
	Code         string
	StopChan     chan string
	ErrChan      chan error
	Working      bool
}

func NewItemWatcher(poesseid string, league string, code string, stopChan chan string, errChan chan error) (*ItemWatcher, error) {

	client := &http.Client{}

	wsConn, err := connections.NewWSConnection(poesseid, league, code)
	if err != nil {
		return nil, err
	}

	watcher := &ItemWatcher{
		WSConnection: wsConn,
		Fetcher:      connections.NewFetcher(client, headers.GetFetchitemHeaders(poesseid)),
		Whisper:      connections.NewWhisper(client, headers.GetWhisperHeaders(poesseid)),
		Code:         code,
		StopChan:     stopChan,
		ErrChan:      errChan,
		Working:      false,
	}

	go watcher.Stopper()

	return watcher, nil

}

func (w *ItemWatcher) Watch() {

	w.Working = true

	for {
		if !w.Working {
			return
		}
		var ls models.LivesearchNewItem
		err := w.WSConnection.ReadJSON(&ls)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			w.ErrChan <- err
			continue
		}

		itemsInfo, err := w.Fetcher.FetchItems(ls.New, w.Code)
		if err != nil {
			w.ErrChan <- err
			continue
		}
		for _, itemInfo := range *itemsInfo {
			err := w.Whisper.Whisper(itemInfo.Result[0].Listing.WhisperToken)
			if err != nil {
				w.ErrChan <- err
				continue
			}
		}
	}
}

func (w *ItemWatcher) Stop() {

	w.WSConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	w.WSConnection.Close()
	w.Working = false
}

func (w *ItemWatcher) Stopper() {

	for {
		forClose := <-w.StopChan
		if forClose == w.Code {
			w.WSConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			w.WSConnection.Close()
			w.Working = false
		}
	}

}
