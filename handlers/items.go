package handlers

import (
	"log"
	"net/http"
	"poebot/connections"
	"poebot/connections/headers"
	"poebot/connections/models"
	"strings"

	"github.com/gorilla/websocket"
)

type ItemHandler struct {
	WSConnection *websocket.Conn
	Fetcher      *connections.Fetcher
	Whisper      *connections.Whisper
	Controller   chan int
	Errors       chan error
}

func NewItemHandler(poesseid string, controller chan int, errors chan error, conn *websocket.Conn) *ItemHandler {

	client := &http.Client{}

	return &ItemHandler{
		WSConnection: conn,
		Fetcher:      connections.NewFetcher(client, headers.GetFetchitemHeaders(poesseid)),
		Whisper:      connections.NewWhisper(client, headers.GetWhisperHeaders(poesseid)),
		Controller:   controller,
		Errors:       errors,
	}

}

func (h *ItemHandler) Serve() {

	for {
		var ls models.LivesearchNewItem
		err := h.WSConnection.ReadJSON(&ls)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("closed")
				break
			}
			log.Printf("read error: %v\n", err)
			h.Errors <- err
			continue
		}

		itemsInfo, err := h.Fetcher.FetchItems(ls.New)
		if err != nil {
			log.Printf("fetch error: %v\n", err)
			h.Errors <- err
			continue
		}
		for _, itemInfo := range *itemsInfo {
			err := h.Whisper.Whisper(itemInfo.Result[0].Listing.WhisperToken)
			if err != nil {
				log.Printf("Whisper error: %v\n", err)
				h.Errors <- err
				continue
			}
		}
	}
}

func (h *ItemHandler) Close() {

	h.WSConnection.Close()
}
