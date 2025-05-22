package connections

import (
	"fmt"
	"net/http"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"

	"github.com/gorilla/websocket"
)

func NewWSConnection(poesessid string, league string, code string) (*websocket.Conn, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://www.pathofexile.com/api/trade/live/%v/%v", league, code), headers.GetLivesearchHeaders(poesessid))
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("invalid link code, check if it is correct")
		case http.StatusTooManyRequests:
			return nil, fmt.Errorf("live search limit reached, try to turn off one your other active links")
		default:
		}
	}
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
