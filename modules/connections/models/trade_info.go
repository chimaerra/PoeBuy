package models

import "time"

type TradeInfo struct {
	Nickname string
	Leagues  []League
}

type League struct {
	ID    string `json:"id"`
	Realm string `json:"realm"`
	Text  string `json:"text"`
}

type PoeApiLeagueResponse struct {
	Leagues []struct {
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Realm       string      `json:"realm"`
		URL         string      `json:"url"`
		StartAt     time.Time   `json:"startAt"`
		EndAt       interface{} `json:"endAt"`
		Description string      `json:"description"`
		Category    struct {
			ID string `json:"id"`
		} `json:"category"`
		RegisterAt time.Time     `json:"registerAt,omitempty"`
		DelveEvent bool          `json:"delveEvent"`
		Rules      []interface{} `json:"rules"`
	} `json:"leagues"`
}

func (t *TradeInfo) GetLeagues() []string {

	leagues := make([]string, 0, 10)

	for _, leg := range t.Leagues {
		leagues = append(leagues, leg.ID)
	}

	return leagues
}
