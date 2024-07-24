package connections

import (
	"fmt"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"

	"github.com/gorilla/websocket"
)

func NewWSConnection(poesessid string, league string, code string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://www.pathofexile.com/api/trade/live/%v/%v", league, code), headers.GetLivesearchHeaders(poesessid))
	if err != nil {
		return nil, fmt.Errorf("websocket connection error: %v", err)
	}
	var r models.LivesearchAuthStatus
	conn.ReadJSON(&r)
	if !r.Auth {
		conn.Close()
		return nil, fmt.Errorf("authentification failed")
	}

	return conn, nil
}
