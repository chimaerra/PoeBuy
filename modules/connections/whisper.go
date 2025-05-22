package connections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"poebuy/modules/connections/models"
)

type Whisper struct {
	Client *http.Client
	Header http.Header
}

func NewWhisper(client *http.Client, header http.Header) *Whisper {
	return &Whisper{
		Client: client,
		Header: header,
	}
}

func (w *Whisper) Whisper(token string) error {
	jsonBody := []byte(fmt.Sprintf("{\"token\": \"%v\"}", token))
	whisperReq, err := http.NewRequest("POST", "https://www.pathofexile.com/api/trade/whisper", bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("whisper request creation error: %v", err)

	}
	whisperReq.Header = w.Header
	whisperResp, err := w.Client.Do(whisperReq)
	if err != nil {
		return fmt.Errorf("whisper request error: %v", err)

	}
	defer whisperResp.Body.Close()

	if whisperResp.StatusCode != 200 {
		errorMsg := &models.WhisperErrorResponse{}
		r, _ := io.ReadAll(whisperResp.Body)
		json.Unmarshal(r, errorMsg)
		return fmt.Errorf("Whisper error: %v %v", whisperResp.Status, errorMsg.Error.Message)
	}

	return nil
}
