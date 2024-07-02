package connections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"poebot/modules/connections/models"

	"github.com/andybalholm/brotli"
)

type Fetcher struct {
	Client *http.Client
	Header http.Header
}

func NewFetcher(client *http.Client, header http.Header) *Fetcher {
	return &Fetcher{
		Client: client,
		Header: header,
	}
}

func (f *Fetcher) FetchItems(items []string) (*[]models.FetchItem, error) {

	entities := []models.FetchItem{}

	for _, item := range items {
		reqBody := bytes.NewBuffer([]byte{})
		ItemReq, err := http.NewRequest("GET", fmt.Sprintf("https://www.pathofexile.com/api/trade/fetch/%v?query=LQEV88GUn", item), reqBody)
		if err != nil {
			return nil, fmt.Errorf("item request creation error: %v", err)

		}
		ItemReq.Header = f.Header

		ItemResp, err := f.Client.Do(ItemReq)
		if err != nil {
			return nil, fmt.Errorf("reqvest sending error: %v", err)

		}
		if ItemResp.StatusCode != http.StatusOK {
			log.Println("Get Item error: ", ItemResp.Status)
			continue
		}

		br := brotli.NewReader(ItemResp.Body)

		bodyBytes, err := io.ReadAll(br)
		if err != nil {
			return nil, fmt.Errorf("request reading error: %v", err)

		}

		itemInfo := models.FetchItem{}
		err = json.Unmarshal(bodyBytes, &itemInfo)
		if err != nil {
			return nil, fmt.Errorf("item unmarshal error: %v", err)
		}
		entities = append(entities, itemInfo)
		ItemResp.Body.Close()
	}

	return &entities, nil
}
