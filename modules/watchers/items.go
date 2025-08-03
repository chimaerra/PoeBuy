package watchers

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"poebuy/modules/connections"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"

	"github.com/gorilla/websocket"
)

type ItemWatcher struct {
	WSConnection        *websocket.Conn
	Fetcher             *connections.Fetcher
	Whisper             *connections.Whisper
	Code                string
	ErrChan             chan error

	delay               time.Duration
	readReady           bool
	index               int
	UpdateCheckmarkFunc func(int)

	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.Mutex
	reconnectDelay time.Duration

	poesseid  string
	league    string
	soundFile string
}

func NewItemWatcher(poesseid string, league string, code string, errChan chan error, delay int64, index int, updateCheckmarkFunc func(int), soundFile string) (*ItemWatcher, error) {
	client := &http.Client{}
	ctx, cancel := context.WithCancel(context.Background())

	return &ItemWatcher{
		Fetcher:             connections.NewFetcher(client, headers.GetFetchitemHeaders(poesseid)),
		Whisper:             connections.NewWhisper(client, headers.GetWhisperHeaders(poesseid), soundFile),
		Code:                code,
		ErrChan:             errChan,
		delay:               time.Millisecond * time.Duration(delay),
		readReady:           true,
		index:               index,
		UpdateCheckmarkFunc: updateCheckmarkFunc,
		ctx:                 ctx,
		cancel:              cancel,
		reconnectDelay:      1 * time.Second,
		poesseid:            poesseid,
		league:              league,
		soundFile:           soundFile,
	}, nil
}

func (w *ItemWatcher) Watch() {
	defer func() {
		if r := recover(); r != nil {
			w.ErrChan <- connections.ErrWatcherPanicked
		}
	}()

	if w.delay > 0 {
		go w.delayer()
	}

	backoff := w.reconnectDelay

	for {
		select {
		case <-w.ctx.Done():
			w.closeConnection()
			w.UpdateCheckmarkFunc(w.index)
			return
		default:
			wsConn, err := connections.NewWSConnection(w.poesseid, w.league, w.Code)
			if err != nil {
				w.ErrChan <- err
				time.Sleep(backoff)
				backoff = increaseBackoff(backoff)
				continue
			}

			w.mu.Lock()
			w.WSConnection = wsConn
			w.mu.Unlock()
			backoff = w.reconnectDelay

			w.readReady = true

			for {
				select {
				case <-w.ctx.Done():
					w.closeConnection()
					w.UpdateCheckmarkFunc(w.index)
					return
				default:
					var ls models.LivesearchNewItem
					err := w.WSConnection.ReadJSON(&ls)
					if err != nil {
						if !strings.Contains(err.Error(), "use of closed network connection") {
							w.ErrChan <- err
						}
						w.closeConnection()
						break
					}

					if w.delay > 0 && !w.readReady {
						continue
					}

					length := len(ls.New)
					if w.delay != 0 && length > 3 {
						length = 3
					}

					itemsInfo, err := w.Fetcher.FetchItems(ls.New[:length], w.Code)
					if err != nil {
						w.ErrChan <- err
						continue
					}

					for _, itemInfo := range itemsInfo {
						if err := w.Whisper.Whisper(itemInfo.Result[0].Listing.WhisperToken); err != nil {
							w.ErrChan <- err
						}
					}

					w.readReady = false
				}
			}
		}
	}
}

func (w *ItemWatcher) Stop() {
	w.cancel()
}

func (w *ItemWatcher) closeConnection() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.WSConnection != nil {
		w.WSConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		w.WSConnection.Close()
		w.WSConnection = nil
	}
}

func (w *ItemWatcher) delayer() {
	ticker := time.NewTicker(w.delay)
	defer ticker.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.readReady = true
		}
	}
}

func increaseBackoff(current time.Duration) time.Duration {
	const maxBackoff = 30 * time.Second
	if current >= maxBackoff {
		return maxBackoff
	}
	return current * 2
}